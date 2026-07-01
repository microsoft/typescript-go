package api

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"gotest.tools/v3/assert"
)

// readResponsePayload decodes a MessagePack response tuple written to buf and
// returns the payload bytes.
func readResponsePayload(t *testing.T, buf *bytes.Buffer) []byte {
	t.Helper()
	p := NewMessagePackProtocol(buf)
	msgType, _, payload, err := p.readTuple()
	assert.NilError(t, err)
	assert.Equal(t, msgType, MessageTypeResponse)
	return payload
}

func TestMessagePackTimingFooter(t *testing.T) {
	t.Parallel()

	id := jsonrpc.NewIDString("getSourceFile")

	t.Run("appends little-endian footer for timed binary result", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewMessagePackProtocol(&buf)
		data := []byte{0x01, 0x02, 0x03}
		err := p.WriteResponse(id, serverTimedResult{
			result:               RawBinary(data),
			processingTimeMicros: 1234,
		})
		assert.NilError(t, err)

		payload := readResponsePayload(t, &buf)
		assert.Equal(t, len(payload), len(data)+4)
		assert.DeepEqual(t, payload[:len(data)], data)
		assert.Equal(t, binary.LittleEndian.Uint32(payload[len(data):]), uint32(1234))
	})

	t.Run("appends footer for timed JSON result", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewMessagePackProtocol(&buf)
		err := p.WriteResponse(id, serverTimedResult{
			result:               map[string]int{"x": 1},
			processingTimeMicros: 7,
		})
		assert.NilError(t, err)

		payload := readResponsePayload(t, &buf)
		assert.Equal(t, string(payload[:len(payload)-4]), `{"x":1}`)
		assert.Equal(t, binary.LittleEndian.Uint32(payload[len(payload)-4:]), uint32(7))
	})

	t.Run("no footer when result is not wrapped", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewMessagePackProtocol(&buf)
		data := []byte{0xAA, 0xBB}
		err := p.WriteResponse(id, RawBinary(data))
		assert.NilError(t, err)

		payload := readResponsePayload(t, &buf)
		assert.DeepEqual(t, payload, data)
	})
}

func TestDurationToMicros(t *testing.T) {
	t.Parallel()
	assert.Equal(t, durationToMicros(1500*time.Microsecond), uint32(1500))
	assert.Equal(t, durationToMicros(0), uint32(0))
	assert.Equal(t, durationToMicros(-5*time.Second), uint32(0))
	assert.Equal(t, durationToMicros(time.Duration(math.MaxInt64)), uint32(math.MaxUint32))
}

func TestJSONRPCTimingEnvelope(t *testing.T) {
	t.Parallel()

	id := jsonrpc.NewIDString("getSourceFile")

	t.Run("wraps timed result in an envelope", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewJSONRPCProtocol(&buf)
		err := p.WriteResponse(id, serverTimedResult{
			result:               map[string]int{"x": 1},
			processingTimeMicros: 42,
		})
		assert.NilError(t, err)

		msg, err := NewJSONRPCProtocol(&buf).ReadMessage()
		assert.NilError(t, err)
		assert.Assert(t, msg.IsResponse())

		var envelope struct {
			Result           map[string]int `json:"result"`
			ServerTimeMicros uint32         `json:"serverTimeMicros"`
		}
		assert.NilError(t, json.Unmarshal(msg.Result, &envelope))
		assert.Equal(t, envelope.Result["x"], 1)
		assert.Equal(t, envelope.ServerTimeMicros, uint32(42))
	})

	t.Run("wraps nil timed result as null", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewJSONRPCProtocol(&buf)
		err := p.WriteResponse(id, serverTimedResult{
			result:               nil,
			processingTimeMicros: 7,
		})
		assert.NilError(t, err)

		msg, err := NewJSONRPCProtocol(&buf).ReadMessage()
		assert.NilError(t, err)

		var envelope struct {
			Result           *int   `json:"result"`
			ServerTimeMicros uint32 `json:"serverTimeMicros"`
		}
		assert.NilError(t, json.Unmarshal(msg.Result, &envelope))
		assert.Assert(t, envelope.Result == nil)
		assert.Equal(t, envelope.ServerTimeMicros, uint32(7))
	})

	t.Run("no envelope when result is not wrapped", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		p := NewJSONRPCProtocol(&buf)
		err := p.WriteResponse(id, map[string]int{"x": 1})
		assert.NilError(t, err)

		msg, err := NewJSONRPCProtocol(&buf).ReadMessage()
		assert.NilError(t, err)

		var result map[string]int
		assert.NilError(t, json.Unmarshal(msg.Result, &result))
		assert.Equal(t, result["x"], 1)
	})
}
