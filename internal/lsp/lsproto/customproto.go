package lsproto

type JSXClosingTagParams = TextDocumentPositionParams

type JSXClosingTagResponse struct {
	// TODO: TextRange?
	NewText *string `json:"newText"`
}

var TextDocumentJSXClosingTagInfo = RequestInfo[*JSXClosingTagParams, JSXClosingTagResponse]{Method: MethodTextDocumentTypeDefinition}
