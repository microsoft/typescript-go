package fourslash_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func ptrTo[T any](v T) *T {
	return &v
}

var ignored = struct{}{}

var defaultCommitCharacters = []string{".", ",", ";"}

var completionGlobalThisItem = &lsproto.CompletionItem{
	Label:    "globalThis",
	Kind:     ptrTo(lsproto.CompletionItemKindModule),
	SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
}

var completionUndefinedVarItem = &lsproto.CompletionItem{
	Label:    "undefined",
	Kind:     ptrTo(lsproto.CompletionItemKindVariable),
	SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
}

var completionGlobalVars = []fourslash.ExpectedCompletionItem{
	&lsproto.CompletionItem{
		Label:    "AbortController",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AbortSignal",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AbstractRange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ActiveXObject",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AnalyserNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Animation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AnimationEffect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AnimationEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AnimationPlaybackEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AnimationTimeline",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ArrayBuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Attr",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Audio",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioBuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioBufferSourceNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioData",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioDecoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioDestinationNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioEncoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioListener",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioParam",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioParamMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioScheduledSourceNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioWorklet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioWorkletNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AuthenticatorAssertionResponse",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AuthenticatorAttestationResponse",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AuthenticatorResponse",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BarProp",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BaseAudioContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BeforeUnloadEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BiquadFilterNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Blob",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BlobEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Boolean",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "BroadcastChannel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ByteLengthQueuingStrategy",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CDATASection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSS",
		Kind:     ptrTo(lsproto.CompletionItemKindModule),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSAnimation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSConditionRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSContainerRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSCounterStyleRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSFontFaceRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSFontFeatureValuesRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSFontPaletteValuesRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSGroupingRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSImageValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSImportRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSKeyframeRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSKeyframesRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSKeywordValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSLayerBlockRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSLayerStatementRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathClamp",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathInvert",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathMax",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathMin",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathNegate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathProduct",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathSum",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMathValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMatrixComponent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSMediaRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSNamespaceRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSNestedDeclarations",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSNumericArray",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSNumericValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSPageRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSPerspective",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSPropertyRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSRotate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSRuleList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSScale",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSScopeRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSSkew",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSSkewX",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSSkewY",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSStartingStyleRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSStyleDeclaration",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSStyleRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSStyleSheet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSStyleValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSSupportsRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSTransformComponent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSTransformValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSTransition",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSTranslate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSUnitValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSUnparsedValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSVariableReferenceValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CSSViewTransitionRule",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Cache",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CacheStorage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CanvasCaptureMediaStreamTrack",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CanvasGradient",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CanvasPattern",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CanvasRenderingContext2D",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CaretPosition",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ChannelMergerNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ChannelSplitterNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CharacterData",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Clipboard",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ClipboardEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ClipboardItem",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CloseEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Comment",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CompositionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CompressionStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ConstantSourceNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ContentVisibilityAutoStateChangeEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ConvolverNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CountQueuingStrategy",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Credential",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CredentialsContainer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Crypto",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CryptoKey",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CustomElementRegistry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CustomEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "CustomStateSet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMException",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMImplementation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMMatrix",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMMatrixReadOnly",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMParser",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMPoint",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMPointReadOnly",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMQuad",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMRect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMRectList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMRectReadOnly",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMStringList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMStringMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DOMTokenList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DataTransfer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DataTransferItem",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DataTransferItemList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DataView",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Date",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DecompressionStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DelayNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DeviceMotionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DeviceOrientationEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Document",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DocumentFragment",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DocumentTimeline",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DocumentType",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DragEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "DynamicsCompressorNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Element",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ElementInternals",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EncodedAudioChunk",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EncodedVideoChunk",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Enumerator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Error",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ErrorEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EvalError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Event",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EventCounts",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EventSource",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "EventTarget",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "File",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileReader",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystem",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemDirectoryEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemDirectoryHandle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemDirectoryReader",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemFileEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemFileHandle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemHandle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FileSystemWritableFileStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Float32Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Float64Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FocusEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FontFace",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FontFaceSet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FontFaceSetLoadEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FormData",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FormDataEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "FragmentDirective",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Function",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GainNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Gamepad",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GamepadButton",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GamepadEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GamepadHapticActuator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Geolocation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GeolocationCoordinates",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GeolocationPosition",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "GeolocationPositionError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLAllCollection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLAnchorElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLAreaElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLAudioElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLBRElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLBaseElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLBodyElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLButtonElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLCanvasElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLCollection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDListElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDataElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDataListElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDetailsElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDialogElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDivElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLEmbedElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFieldSetElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFormControlsCollection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFormElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLHRElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLHeadElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLHeadingElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLHtmlElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLIFrameElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLImageElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLInputElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLLIElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLLabelElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLLegendElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLLinkElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMapElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMediaElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMenuElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMetaElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMeterElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLModElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLOListElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLObjectElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLOptGroupElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLOptionElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLOptionsCollection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLOutputElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLParagraphElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLPictureElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLPreElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLProgressElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLQuoteElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLScriptElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLSelectElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLSlotElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLSourceElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLSpanElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLStyleElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableCaptionElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableCellElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableColElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableRowElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTableSectionElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTemplateElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTextAreaElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTimeElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTitleElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLTrackElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLUListElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLUnknownElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLVideoElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HashChangeEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Headers",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Highlight",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "HighlightRegistry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "History",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBCursor",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBCursorWithValue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBDatabase",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBFactory",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBIndex",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBKeyRange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBObjectStore",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBOpenDBRequest",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBRequest",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBTransaction",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IDBVersionChangeEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IIRFilterNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IdleDeadline",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Image",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageBitmap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageBitmapRenderingContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageData",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageDecoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageTrack",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ImageTrackList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Infinity",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "InputDeviceInfo",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "InputEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Int16Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Int32Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Int8Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IntersectionObserver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "IntersectionObserverEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Intl",
		Kind:     ptrTo(lsproto.CompletionItemKindModule),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "JSON",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "KeyboardEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "KeyframeEffect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "LargestContentfulPaint",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Location",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Lock",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "LockManager",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIAccess",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIConnectionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIInput",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIInputMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIMessageEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIOutput",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIOutputMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MIDIPort",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Math",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MathMLElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaCapabilities",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaDeviceInfo",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaDevices",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaElementAudioSourceNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaEncryptedEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaKeyMessageEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaKeySession",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaKeyStatusMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaKeySystemAccess",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaKeys",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaMetadata",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaQueryList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaQueryListEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaRecorder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaSession",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaSource",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaSourceHandle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaStreamAudioDestinationNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaStreamAudioSourceNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaStreamTrack",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MediaStreamTrackEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MessageChannel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MessageEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MessagePort",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MouseEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MutationObserver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "MutationRecord",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NaN",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NamedNodeMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NavigationActivation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NavigationHistoryEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NavigationPreloadManager",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Navigator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Node",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NodeFilter",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NodeIterator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "NodeList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Notification",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Number",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Object",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OfflineAudioCompletionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OfflineAudioContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OffscreenCanvas",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OffscreenCanvasRenderingContext2D",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Option",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OscillatorNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "OverconstrainedError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PageRevealEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PageSwapEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PageTransitionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PannerNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Path2D",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PaymentAddress",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PaymentMethodChangeEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PaymentRequest",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PaymentRequestUpdateEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PaymentResponse",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Performance",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceEventTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceMark",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceMeasure",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceNavigationTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceObserver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceObserverEntryList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformancePaintTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceResourceTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceServerTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PeriodicWave",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PermissionStatus",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Permissions",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PictureInPictureEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PictureInPictureWindow",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PointerEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PopStateEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ProcessingInstruction",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ProgressEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PromiseRejectionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PublicKeyCredential",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PushManager",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PushSubscription",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "PushSubscriptionOptions",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCCertificate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCDTMFSender",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCDTMFToneChangeEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCDataChannel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCDataChannelEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCDtlsTransport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCEncodedAudioFrame",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCEncodedVideoFrame",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCErrorEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCIceCandidate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCIceTransport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCPeerConnection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCPeerConnectionIceErrorEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCPeerConnectionIceEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCRtpReceiver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCRtpScriptTransform",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCRtpSender",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCRtpTransceiver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCSctpTransport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCSessionDescription",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCStatsReport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RTCTrackEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RadioNodeList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Range",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RangeError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableByteStreamController",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableStreamBYOBReader",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableStreamBYOBRequest",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableStreamDefaultController",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReadableStreamDefaultReader",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReferenceError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RegExp",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "RemotePlayback",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Report",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReportBody",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ReportingObserver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Request",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ResizeObserver",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ResizeObserverEntry",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ResizeObserverSize",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Response",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAngle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimateElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimateMotionElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimateTransformElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedAngle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedBoolean",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedEnumeration",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedInteger",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedLength",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedLengthList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedNumber",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedNumberList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedPreserveAspectRatio",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedRect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedString",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimatedTransformList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGAnimationElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGCircleElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGClipPathElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGComponentTransferFunctionElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGDefsElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGDescElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGEllipseElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEBlendElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEColorMatrixElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEComponentTransferElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFECompositeElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEConvolveMatrixElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEDiffuseLightingElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEDisplacementMapElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEDistantLightElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEDropShadowElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEFloodElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEFuncAElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEFuncBElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEFuncGElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEFuncRElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEGaussianBlurElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEImageElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEMergeElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEMergeNodeElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEMorphologyElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEOffsetElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFEPointLightElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFESpecularLightingElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFESpotLightElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFETileElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFETurbulenceElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGFilterElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGForeignObjectElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGGElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGGeometryElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGGradientElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGGraphicsElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGImageElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGLength",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGLengthList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGLineElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGLinearGradientElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGMPathElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGMarkerElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGMaskElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGMatrix",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGMetadataElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGNumber",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGNumberList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPathElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPatternElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPoint",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPointList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPolygonElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPolylineElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGPreserveAspectRatio",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGRadialGradientElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGRect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGRectElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGSVGElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGScriptElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGSetElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGStopElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGStringList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGStyleElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGSwitchElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGSymbolElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTSpanElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTextContentElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTextElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTextPathElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTextPositioningElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTitleElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTransform",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGTransformList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGUnitTypes",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGUseElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SVGViewElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SafeArray",
		Kind:     ptrTo(lsproto.CompletionItemKindClass),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Screen",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ScreenOrientation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SecurityPolicyViolationEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Selection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ServiceWorker",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ServiceWorkerContainer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ServiceWorkerRegistration",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ShadowRoot",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SharedWorker",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SourceBuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SourceBufferList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechRecognitionAlternative",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechRecognitionResult",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechRecognitionResultList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechSynthesis",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechSynthesisErrorEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechSynthesisEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechSynthesisUtterance",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SpeechSynthesisVoice",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StaticRange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StereoPannerNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Storage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StorageEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StorageManager",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "String",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StylePropertyMap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StylePropertyMapReadOnly",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StyleSheet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "StyleSheetList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SubmitEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SubtleCrypto",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "SyntaxError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Text",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextDecoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextDecoderStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextEncoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextEncoderStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextMetrics",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextTrack",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextTrackCue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextTrackCueList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TextTrackList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TimeRanges",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ToggleEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Touch",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TouchEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TouchList",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TrackEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TransformStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TransformStreamDefaultController",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TransitionEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TreeWalker",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "TypeError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "UIEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "URIError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "URL",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "URLSearchParams",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Uint16Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Uint32Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Uint8Array",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Uint8ClampedArray",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "UserActivation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VBArray",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VTTCue",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VTTRegion",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ValidityState",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VarDate",
		Kind:     ptrTo(lsproto.CompletionItemKindClass),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VideoColorSpace",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VideoDecoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VideoEncoder",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VideoFrame",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VideoPlaybackQuality",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ViewTransition",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ViewTransitionTypeSet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "VisualViewport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WSH",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WScript",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WakeLock",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WakeLockSentinel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WaveShaperNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebAssembly",
		Kind:     ptrTo(lsproto.CompletionItemKindModule),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGL2RenderingContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLActiveInfo",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLBuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLContextEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLFramebuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLProgram",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLQuery",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLRenderbuffer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLRenderingContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLSampler",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLShader",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLShaderPrecisionFormat",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLSync",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLTexture",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLTransformFeedback",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLUniformLocation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebGLVertexArrayObject",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebKitCSSMatrix",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebSocket",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebTransport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebTransportBidirectionalStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebTransportDatagramDuplexStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WebTransportError",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WheelEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Window",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Worker",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "Worklet",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WritableStream",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WritableStreamDefaultController",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "WritableStreamDefaultWriter",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XMLDocument",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XMLHttpRequest",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XMLHttpRequestEventTarget",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XMLHttpRequestUpload",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XMLSerializer",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XPathEvaluator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XPathExpression",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XPathResult",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "XSLTProcessor",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "addEventListener",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "alert",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "atob",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "btoa",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "caches",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "cancelAnimationFrame",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "cancelIdleCallback",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "clearInterval",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "clearTimeout",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "close",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "closed",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "confirm",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "console",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "createImageBitmap",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "crossOriginIsolated",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "crypto",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "customElements",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "decodeURI",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "decodeURIComponent",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "devicePixelRatio",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "dispatchEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "document",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "encodeURI",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "encodeURIComponent",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "eval",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "fetch",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "focus",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "frameElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "frames",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "getComputedStyle",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "getSelection",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "history",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "importScripts",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "indexedDB",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "innerHeight",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "innerWidth",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "isFinite",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "isNaN",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "isSecureContext",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "length",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "localStorage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "location",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "locationbar",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "matchMedia",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "menubar",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "moveBy",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "moveTo",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "navigator",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onabort",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onafterprint",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onanimationcancel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onanimationend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onanimationiteration",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onanimationstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onauxclick",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onbeforeinput",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onbeforeprint",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onbeforetoggle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onbeforeunload",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onblur",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncancel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncanplay",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncanplaythrough",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onclick",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onclose",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncontextlost",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncontextmenu",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncontextrestored",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncopy",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncuechange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oncut",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondblclick",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondevicemotion",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondeviceorientation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondeviceorientationabsolute",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondrag",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondragend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondragenter",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondragleave",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondragover",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondragstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondrop",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ondurationchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onemptied",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onended",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onerror",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onfocus",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onformdata",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ongamepadconnected",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ongamepaddisconnected",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ongotpointercapture",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onhashchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oninput",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "oninvalid",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onkeydown",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onkeyup",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onlanguagechange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onload",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onloadeddata",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onloadedmetadata",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onloadstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onlostpointercapture",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmessage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmessageerror",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmousedown",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmouseenter",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmouseleave",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmousemove",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmouseout",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmouseover",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onmouseup",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onoffline",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ononline",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpagehide",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpagereveal",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpageshow",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpageswap",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpaste",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpause",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onplay",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onplaying",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointercancel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerdown",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerenter",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerleave",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointermove",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerout",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerover",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpointerup",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onpopstate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onprogress",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onratechange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onrejectionhandled",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onreset",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onresize",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onscroll",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onscrollend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onsecuritypolicyviolation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onseeked",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onseeking",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onselect",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onselectionchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onselectstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onslotchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onstalled",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onstorage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onsubmit",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onsuspend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontimeupdate",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontoggle",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontouchcancel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontouchend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontouchmove",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontouchstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontransitioncancel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontransitionend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontransitionrun",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "ontransitionstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onunhandledrejection",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onvolumechange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onwaiting",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "onwheel",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "open",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "opener",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "origin",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "outerHeight",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "outerWidth",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "pageXOffset",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "pageYOffset",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "parent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "parseFloat",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "parseInt",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "performance",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "personalbar",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "postMessage",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "print",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "prompt",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "queueMicrotask",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "removeEventListener",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "reportError",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "requestAnimationFrame",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "requestIdleCallback",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "resizeBy",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "resizeTo",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "screen",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "screenLeft",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "screenTop",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "screenX",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "screenY",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scroll",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scrollBy",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scrollTo",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scrollX",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scrollY",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "scrollbars",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "self",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "sessionStorage",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "setInterval",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "setTimeout",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "speechSynthesis",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "statusbar",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "stop",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "structuredClone",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "toString",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "toolbar",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "top",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "visualViewport",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "webkitURL",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "window",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "AudioProcessingEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "External",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDirectoryElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLDocument",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFontElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFrameElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLFrameSetElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLMarqueeElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "HTMLParamElement",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "MimeType",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "MimeTypeArray",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceNavigation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "PerformanceTiming",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "Plugin",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "PluginArray",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "ScriptProcessorNode",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "TextEvent",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "blur",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "captureEvents",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "clientInformation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "escape",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "event",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "external",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "name",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onkeypress",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onorientationchange",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onunload",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onwebkitanimationend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onwebkitanimationiteration",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onwebkitanimationstart",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "onwebkittransitionend",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "orientation",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "releaseEvents",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "status",
		Kind:     ptrTo(lsproto.CompletionItemKindVariable),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
	&lsproto.CompletionItem{
		Label:    "unescape",
		Kind:     ptrTo(lsproto.CompletionItemKindFunction),
		SortText: ptrTo(string(ls.DeprecateSortText(ls.SortTextGlobalsOrKeywords))),
	},
}

var completionGlobalKeywords = []fourslash.ExpectedCompletionItem{
	&lsproto.CompletionItem{
		Label:    "abstract",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "any",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "as",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "asserts",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "async",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "await",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "bigint",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "boolean",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "break",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "case",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "catch",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "class",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "const",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "continue",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "debugger",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "declare",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "default",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "delete",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "do",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "else",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "enum",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "export",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "extends",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "false",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "finally",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "for",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "function",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "if",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "implements",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "import",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "in",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "infer",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "instanceof",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "interface",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "keyof",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "let",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "module",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "namespace",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "never",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "new",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "null",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "number",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "object",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "package",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "readonly",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "return",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "satisfies",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "string",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "super",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "switch",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "symbol",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "this",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "throw",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "true",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "try",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "type",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "typeof",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "unique",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "unknown",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "using",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "var",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "void",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "while",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "with",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
	&lsproto.CompletionItem{
		Label:    "yield",
		Kind:     ptrTo(lsproto.CompletionItemKindKeyword),
		SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
	},
}

var completionGlobals = sortCompletionItems(append(
	append(completionGlobalVars, completionGlobalKeywords...),
	completionGlobalThisItem,
	completionUndefinedVarItem,
))

func sortCompletionItems(items []fourslash.ExpectedCompletionItem) []fourslash.ExpectedCompletionItem {
	items = slices.Clone(items)
	slices.SortStableFunc(items, func(a fourslash.ExpectedCompletionItem, b fourslash.ExpectedCompletionItem) int {
		defaultSortText := string(ls.SortTextLocationPriority)
		var aSortText, bSortText string
		switch a := a.(type) {
		case *lsproto.CompletionItem:
			if a.SortText != nil {
				aSortText = *a.SortText
			}
		}
		switch b := b.(type) {
		case *lsproto.CompletionItem:
			if b.SortText != nil {
				bSortText = *b.SortText
			}
		}
		aSortText = core.OrElse(aSortText, defaultSortText)
		bSortText = core.OrElse(bSortText, defaultSortText)
		bySortText := cmp.Compare(aSortText, bSortText)
		if bySortText != 0 {
			return bySortText
		}
		var aLabel, bLabel string
		switch a := a.(type) {
		case *lsproto.CompletionItem:
			aLabel = a.Label
		case string:
			aLabel = a
		default:
			panic(fmt.Sprintf("unexpected completion item type: %T", a))
		}
		switch b := b.(type) {
		case *lsproto.CompletionItem:
			bLabel = b.Label
		case string:
			bLabel = b
		default:
			panic(fmt.Sprintf("unexpected completion item type: %T", b))
		}
		return cmp.Compare(aLabel, bLabel)
	})
	return items
}

func completionGlobalsPlus(items []fourslash.ExpectedCompletionItem, noLib bool) []fourslash.ExpectedCompletionItem {
	var all []fourslash.ExpectedCompletionItem
	if noLib {
		all = append(
			append(items, completionGlobalThisItem, completionUndefinedVarItem),
			completionGlobalKeywords...,
		)
	} else {
		all = append(items, completionGlobals...)
	}
	return sortCompletionItems(all)
}
