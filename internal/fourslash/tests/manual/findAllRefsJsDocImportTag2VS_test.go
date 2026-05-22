package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsJsDocImportTag2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @checkJs: true
// @Filename: /component.js
export default class Component {
  constructor() {
    this.id_ = Math.random();
  }
  id() {
    return this.id_;
  }
}
// @Filename: /spatial-navigation.js
/** @import Component from './component.js' */

export class SpatialNavigation {
  /**
   * @param {Component} component
   */
  add(component) {}
}
// @Filename: /player.js
import Component from './component.js';

/**
 * @extends Component/*1*/
 */
export class Player extends Component {}`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1")
}
