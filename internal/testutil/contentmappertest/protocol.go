// Package contentmappertest provides realistic content mapper implementations used by tests.
package contentmappertest

import (
	"context"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

type noNotifications struct{}

func (noNotifications) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

func initializeResult(source string) contentmapper.InitializeResult {
	return contentmapper.InitializeResult{
		ProtocolVersion:  contentmapper.ProtocolVersion,
		PositionEncoding: contentmapper.PositionEncodingUTF8,
		DiagnosticSource: source,
	}
}

func identityMappedOutput(content string) (contentmapper.MappedOutput, error) {
	mappings, err := spanmap.New([]spanmap.Segment{{
		VirtualEnd:  core.TextPos(len(content)),
		OriginalEnd: core.TextPos(len(content)),
		Kind:        spanmap.KindVerbatim,
		Features:    spanmap.FeatureAll,
	}}).Marshal()
	if err != nil {
		return contentmapper.MappedOutput{}, err
	}
	return contentmapper.MappedOutput{Text: content, Extension: ".ts", Mappings: json.Value(mappings)}, nil
}

type staticProjectHandler struct{ ipc.Handler }

func (h staticProjectHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodOpenProject:
		return contentmapper.OpenProjectResult{}, nil
	case contentmapper.MethodCloseProject:
		return nil, nil
	default:
		return h.Handler.HandleRequest(ctx, method, params)
	}
}
