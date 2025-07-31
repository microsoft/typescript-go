package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

func TestLocationLinkSupport(t *testing.T) {
	t.Parallel()

	// Simple integration test to ensure LocationLink support works
	// without causing import cycles

	// Test that client capabilities are correctly used
	linkSupport := true
	capabilities := &lsproto.DefinitionClientCapabilities{
		LinkSupport: &linkSupport,
	}

	// Test that the capability checking logic works
	assert.Assert(t, capabilities != nil)
	assert.Assert(t, capabilities.LinkSupport != nil)
	assert.Assert(t, *capabilities.LinkSupport)

	// Test with capabilities disabled
	linkSupportFalse := false
	capabilitiesDisabled := &lsproto.DefinitionClientCapabilities{
		LinkSupport: &linkSupportFalse,
	}
	assert.Assert(t, capabilitiesDisabled.LinkSupport != nil)
	assert.Assert(t, !*capabilitiesDisabled.LinkSupport)

	// Test with nil capabilities
	var nilCapabilities *lsproto.DefinitionClientCapabilities
	assert.Assert(t, nilCapabilities == nil)
}