package lsutil

import (
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
)

var DefaultUserPreferences = &UserPreferences{
	IncludeCompletionsForModuleExports:    core.TSTrue,
	IncludeCompletionsForImportStatements: core.TSTrue,

	AllowRenameOfImportPath:            true,
	ProvideRefactorNotApplicableReason: true,
	IncludeCompletionsWithSnippetText:  core.TSTrue,
	DisplayPartsForJSDoc:               true,
	DisableLineTextInReferences:        true,
	ReportStyleChecksAsWarnings:        true,

	ExcludeLibrarySymbolsInNavTo: true,
}

// UserPreferences represents TypeScript language service preferences.
//
// Fields can be populated from two sources:
//  1. Direct property names (from VS Code's "unstable" config or direct tsserver protocol)
//     These are matched case-insensitively by Go field name.
//  2. VS Code's nested config structure (e.g., "suggest.autoImports", "inlayHints.parameterNames.enabled")
//     These use the `pref` tag to specify the dotted path.
//
// The `pref` tag format: "path.to.setting" or "path.to.setting,invert" for boolean inversion.
type UserPreferences struct {
	QuotePreference                           QuotePreference `pref:"preferences.quoteStyle"`
	LazyConfiguredProjectsFromExternalProject bool            // !!!

	// A positive integer indicating the maximum length of a hover text before it is truncated.
	//
	// Default: `500`
	MaximumHoverLength int // !!!

	// ------- Completions -------

	// If enabled, TypeScript will search through all external modules' exports and add them to the completions list.
	// This affects lone identifier completions but not completions on the right hand side of `obj.`.
	IncludeCompletionsForModuleExports core.Tristate `pref:"suggest.autoImports"`
	// Enables auto-import-style completions on partially-typed import statements. E.g., allows
	// `import write|` to be completed to `import { writeFile } from "fs"`.
	IncludeCompletionsForImportStatements core.Tristate `pref:"suggest.includeCompletionsForImportStatements"`
	// Unless this option is `false`,  member completion lists triggered with `.` will include entries
	// on potentially-null and potentially-undefined values, with insertion text to replace
	// preceding `.` tokens with `?.`.
	IncludeAutomaticOptionalChainCompletions core.Tristate `pref:"suggest.includeAutomaticOptionalChainCompletions"`
	// Allows completions to be formatted with snippet text, indicated by `CompletionItem["isSnippet"]`.
	IncludeCompletionsWithSnippetText core.Tristate // !!!
	// If enabled, completions for class members (e.g. methods and properties) will include
	// a whole declaration for the member.
	// E.g., `class A { f| }` could be completed to `class A { foo(): number {} }`, instead of
	// `class A { foo }`.
	IncludeCompletionsWithClassMemberSnippets core.Tristate `pref:"suggest.classMemberSnippets.enabled"` // !!!
	// If enabled, object literal methods will have a method declaration completion entry in addition
	// to the regular completion entry containing just the method name.
	// E.g., `const objectLiteral: T = { f| }` could be completed to `const objectLiteral: T = { foo(): void {} }`,
	// in addition to `const objectLiteral: T = { foo }`.
	IncludeCompletionsWithObjectLiteralMethodSnippets core.Tristate               `pref:"suggest.objectLiteralMethodSnippets.enabled"` // !!!
	JsxAttributeCompletionStyle                       JsxAttributeCompletionStyle `pref:"preferences.jsxAttributeCompletionStyle"`

	// ------- AutoImports --------

	ModuleSpecifier ModuleSpecifierUserPreferences

	IncludePackageJsonAutoImports IncludePackageJsonAutoImports `pref:"preferences.includePackageJsonAutoImports"` // !!!
	AutoImportFileExcludePatterns []string                      `pref:"preferences.autoImportFileExcludePatterns"` // !!!
	PreferTypeOnlyAutoImports     bool                          `pref:"preferences.preferTypeOnlyAutoImports"`     // !!!

	// ------- OrganizeImports -------

	// Indicates whether imports should be organized in a case-insensitive manner.
	//
	// Default: TSUnknown ("auto" in strada), will perform detection
	OrganizeImportsIgnoreCase core.Tristate `pref:"preferences.organizeImports.caseSensitivity"` // !!!
	// Indicates whether imports should be organized via an "ordinal" (binary) comparison using the numeric value of their
	// code points, or via "unicode" collation (via the Unicode Collation Algorithm (https://unicode.org/reports/tr10/#Scope))
	//
	// using rules associated with the locale specified in organizeImportsCollationLocale.
	//
	// Default: Ordinal
	OrganizeImportsCollation OrganizeImportsCollation `pref:"preferences.organizeImports.unicodeCollation"` // !!!
	// Indicates the locale to use for "unicode" collation. If not specified, the locale `"en"` is used as an invariant
	// for the sake of consistent sorting. Use `"auto"` to use the detected UI locale.
	//
	// This preference is ignored if organizeImportsCollation is not `unicode`.
	//
	// Default: `"en"`
	OrganizeImportsLocale string `pref:"preferences.organizeImports.locale"` // !!!
	// Indicates whether numeric collation should be used for digit sequences in strings. When `true`, will collate
	// strings such that `a1z < a2z < a100z`. When `false`, will collate strings such that `a1z < a100z < a2z`.
	//
	// This preference is ignored if organizeImportsCollation is not `unicode`.
	//
	// Default: `false`
	OrganizeImportsNumericCollation bool `pref:"preferences.organizeImports.numericCollation"` // !!!
	// Indicates whether accents and other diacritic marks are considered unequal for the purpose of collation. When
	// `true`, characters with accents and other diacritics will be collated in the order defined by the locale specified
	// in organizeImportsCollationLocale.
	//
	// This preference is ignored if organizeImportsCollation is not `unicode`.
	//
	// Default: `true`
	OrganizeImportsAccentCollation bool `pref:"preferences.organizeImports.accentCollation"` // !!!
	// Indicates whether upper case or lower case should sort first. When `false`, the default order for the locale
	// specified in organizeImportsCollationLocale is used.
	//
	// This preference is ignored if:
	//	- organizeImportsCollation is not `unicode`
	//	- organizeImportsIgnoreCase is `true`
	//	- organizeImportsIgnoreCase is `auto` and the auto-detected case sensitivity is case-insensitive.
	//
	// Default: `false`
	OrganizeImportsCaseFirst OrganizeImportsCaseFirst `pref:"preferences.organizeImports.caseFirst"` // !!!
	// Indicates where named type-only imports should sort. "inline" sorts named imports without regard to if the import is type-only.
	//
	// Default: `auto`, which defaults to `last`
	OrganizeImportsTypeOrder OrganizeImportsTypeOrder `pref:"preferences.organizeImports.typeOrder"` // !!!

	// ------- MoveToFile -------

	AllowTextChangesInNewFiles bool // !!!

	// ------- Rename -------

	// renamed from `providePrefixAndSuffixTextForRename`
	UseAliasesForRename     core.Tristate `pref:"preferences.useAliasesForRenames,alias:providePrefixAndSuffixTextForRename"`
	AllowRenameOfImportPath bool          // !!!

	// ------- CodeFixes/Refactors -------

	ProvideRefactorNotApplicableReason bool // !!!

	// ------- InlayHints -------

	InlayHints InlayHintsPreferences

	// ------- CodeLens -------

	CodeLens CodeLensUserPreferences

	// ------- Symbols -------

	ExcludeLibrarySymbolsInNavTo bool `pref:"workspaceSymbols.excludeLibrarySymbols"`

	// ------- Misc -------

	DisableSuggestions          bool // !!!
	DisableLineTextInReferences bool // !!!
	DisplayPartsForJSDoc        bool // !!!
	ReportStyleChecksAsWarnings bool // !!! If this changes, we need to ask the client to recompute diagnostics
}

type InlayHintsPreferences struct {
	IncludeInlayParameterNameHints                        IncludeInlayParameterNameHints `pref:"inlayHints.parameterNames.enabled"`
	IncludeInlayParameterNameHintsWhenArgumentMatchesName bool                           `pref:"inlayHints.parameterNames.suppressWhenArgumentMatchesName,invert"`
	IncludeInlayFunctionParameterTypeHints                bool                           `pref:"inlayHints.parameterTypes.enabled"`
	IncludeInlayVariableTypeHints                         bool                           `pref:"inlayHints.variableTypes.enabled"`
	IncludeInlayVariableTypeHintsWhenTypeMatchesName      bool                           `pref:"inlayHints.variableTypes.suppressWhenTypeMatchesName,invert"`
	IncludeInlayPropertyDeclarationTypeHints              bool                           `pref:"inlayHints.propertyDeclarationTypes.enabled"`
	IncludeInlayFunctionLikeReturnTypeHints               bool                           `pref:"inlayHints.functionLikeReturnTypes.enabled"`
	IncludeInlayEnumMemberValueHints                      bool                           `pref:"inlayHints.enumMemberValues.enabled"`
}

type CodeLensUserPreferences struct {
	ReferencesCodeLensEnabled                     bool `pref:"referencesCodeLens.enabled"`
	ImplementationsCodeLensEnabled                bool `pref:"implementationsCodeLens.enabled"`
	ReferencesCodeLensShowOnAllFunctions          bool `pref:"referencesCodeLens.showOnAllFunctions"`
	ImplementationsCodeLensShowOnInterfaceMethods bool `pref:"implementationsCodeLens.showOnInterfaceMethods"`
	ImplementationsCodeLensShowOnAllClassMethods  bool `pref:"implementationsCodeLens.showOnAllClassMethods"`
}

type ModuleSpecifierUserPreferences struct {
	ImportModuleSpecifierPreference modulespecifiers.ImportModuleSpecifierPreference `pref:"preferences.importModuleSpecifier"` // !!!
	// Determines whether we import `foo/index.ts` as "foo", "foo/index", or "foo/index.js"
	ImportModuleSpecifierEnding       modulespecifiers.ImportModuleSpecifierEndingPreference `pref:"preferences.importModuleSpecifierEnding"`       // !!!
	AutoImportSpecifierExcludeRegexes []string                                               `pref:"preferences.autoImportSpecifierExcludeRegexes"` // !!!
}

// --- Enum Types ---

type QuotePreference string

const (
	QuotePreferenceUnknown QuotePreference = ""
	QuotePreferenceAuto    QuotePreference = "auto"
	QuotePreferenceDouble  QuotePreference = "double"
	QuotePreferenceSingle  QuotePreference = "single"
)

type JsxAttributeCompletionStyle string

const (
	JsxAttributeCompletionStyleUnknown JsxAttributeCompletionStyle = ""
	JsxAttributeCompletionStyleAuto    JsxAttributeCompletionStyle = "auto"
	JsxAttributeCompletionStyleBraces  JsxAttributeCompletionStyle = "braces"
	JsxAttributeCompletionStyleNone    JsxAttributeCompletionStyle = "none"
)

type IncludeInlayParameterNameHints string

const (
	IncludeInlayParameterNameHintsNone     IncludeInlayParameterNameHints = ""
	IncludeInlayParameterNameHintsAll      IncludeInlayParameterNameHints = "all"
	IncludeInlayParameterNameHintsLiterals IncludeInlayParameterNameHints = "literals"
)

type IncludePackageJsonAutoImports string

const (
	IncludePackageJsonAutoImportsUnknown IncludePackageJsonAutoImports = ""
	IncludePackageJsonAutoImportsAuto    IncludePackageJsonAutoImports = "auto"
	IncludePackageJsonAutoImportsOn      IncludePackageJsonAutoImports = "on"
	IncludePackageJsonAutoImportsOff     IncludePackageJsonAutoImports = "off"
)

type OrganizeImportsCollation bool

const (
	OrganizeImportsCollationOrdinal OrganizeImportsCollation = false
	OrganizeImportsCollationUnicode OrganizeImportsCollation = true
)

type OrganizeImportsCaseFirst int

const (
	OrganizeImportsCaseFirstFalse OrganizeImportsCaseFirst = 0
	OrganizeImportsCaseFirstLower OrganizeImportsCaseFirst = 1
	OrganizeImportsCaseFirstUpper OrganizeImportsCaseFirst = 2
)

type OrganizeImportsTypeOrder int

const (
	OrganizeImportsTypeOrderAuto   OrganizeImportsTypeOrder = 0
	OrganizeImportsTypeOrderLast   OrganizeImportsTypeOrder = 1
	OrganizeImportsTypeOrderInline OrganizeImportsTypeOrder = 2
	OrganizeImportsTypeOrderFirst  OrganizeImportsTypeOrder = 3
)

// --- Reflection-based parsing infrastructure ---

// typeParsers maps reflect.Type to a function that parses a value into that type.
var typeParsers = map[reflect.Type]func(any) any{
	reflect.TypeFor[core.Tristate](): func(val any) any {
		if b, ok := val.(bool); ok {
			if b {
				return core.TSTrue
			}
			return core.TSFalse
		}
		return core.TSUnknown
	},
	reflect.TypeFor[QuotePreference](): func(val any) any {
		if s, ok := val.(string); ok {
			switch strings.ToLower(s) {
			case "auto":
				return QuotePreferenceAuto
			case "double":
				return QuotePreferenceDouble
			case "single":
				return QuotePreferenceSingle
			}
		}
		return QuotePreferenceUnknown
	},
	reflect.TypeFor[JsxAttributeCompletionStyle](): func(val any) any {
		if s, ok := val.(string); ok {
			switch strings.ToLower(s) {
			case "braces":
				return JsxAttributeCompletionStyleBraces
			case "none":
				return JsxAttributeCompletionStyleNone
			}
		}
		return JsxAttributeCompletionStyleAuto
	},
	reflect.TypeFor[IncludeInlayParameterNameHints](): func(val any) any {
		if s, ok := val.(string); ok {
			switch s {
			case "all":
				return IncludeInlayParameterNameHintsAll
			case "literals":
				return IncludeInlayParameterNameHintsLiterals
			}
		}
		return IncludeInlayParameterNameHintsNone
	},
	reflect.TypeFor[IncludePackageJsonAutoImports](): func(val any) any {
		if s, ok := val.(string); ok {
			switch strings.ToLower(s) {
			case "on":
				return IncludePackageJsonAutoImportsOn
			case "off":
				return IncludePackageJsonAutoImportsOff
			default:
				return IncludePackageJsonAutoImportsAuto
			}
		}
		return IncludePackageJsonAutoImportsUnknown
	},
	reflect.TypeFor[OrganizeImportsCollation](): func(val any) any {
		if s, ok := val.(string); ok && strings.ToLower(s) == "unicode" {
			return OrganizeImportsCollationUnicode
		}
		return OrganizeImportsCollationOrdinal
	},
	reflect.TypeFor[OrganizeImportsCaseFirst](): func(val any) any {
		if s, ok := val.(string); ok {
			switch s {
			case "lower":
				return OrganizeImportsCaseFirstLower
			case "upper":
				return OrganizeImportsCaseFirstUpper
			}
		}
		return OrganizeImportsCaseFirstFalse
	},
	reflect.TypeFor[OrganizeImportsTypeOrder](): func(val any) any {
		if s, ok := val.(string); ok {
			switch s {
			case "last":
				return OrganizeImportsTypeOrderLast
			case "inline":
				return OrganizeImportsTypeOrderInline
			case "first":
				return OrganizeImportsTypeOrderFirst
			}
		}
		return OrganizeImportsTypeOrderAuto
	},
	reflect.TypeFor[modulespecifiers.ImportModuleSpecifierPreference](): func(val any) any {
		if s, ok := val.(string); ok {
			switch strings.ToLower(s) {
			case "project-relative":
				return modulespecifiers.ImportModuleSpecifierPreferenceProjectRelative
			case "relative":
				return modulespecifiers.ImportModuleSpecifierPreferenceRelative
			case "non-relative":
				return modulespecifiers.ImportModuleSpecifierPreferenceNonRelative
			}
		}
		return modulespecifiers.ImportModuleSpecifierPreferenceShortest
	},
	reflect.TypeFor[modulespecifiers.ImportModuleSpecifierEndingPreference](): func(val any) any {
		if s, ok := val.(string); ok {
			switch strings.ToLower(s) {
			case "minimal":
				return modulespecifiers.ImportModuleSpecifierEndingPreferenceMinimal
			case "index":
				return modulespecifiers.ImportModuleSpecifierEndingPreferenceIndex
			case "js":
				return modulespecifiers.ImportModuleSpecifierEndingPreferenceJs
			}
		}
		return modulespecifiers.ImportModuleSpecifierEndingPreferenceAuto
	},
}

// typeSerializers maps reflect.Type to a function that serializes a value of that type.
var typeSerializers = map[reflect.Type]func(any) any{
	reflect.TypeFor[core.Tristate](): func(val any) any {
		switch val.(core.Tristate) {
		case core.TSTrue:
			return true
		case core.TSFalse:
			return false
		default:
			return nil
		}
	},
	reflect.TypeFor[IncludeInlayParameterNameHints](): func(val any) any {
		s := val.(IncludeInlayParameterNameHints)
		if s == "" {
			return "none"
		}
		return string(s)
	},
	reflect.TypeFor[OrganizeImportsCollation](): func(val any) any {
		if val.(OrganizeImportsCollation) == OrganizeImportsCollationUnicode {
			return "unicode"
		}
		return "ordinal"
	},
	reflect.TypeFor[OrganizeImportsCaseFirst](): func(val any) any {
		switch val.(OrganizeImportsCaseFirst) {
		case OrganizeImportsCaseFirstLower:
			return "lower"
		case OrganizeImportsCaseFirstUpper:
			return "upper"
		default:
			return "default"
		}
	},
	reflect.TypeFor[OrganizeImportsTypeOrder](): func(val any) any {
		switch val.(OrganizeImportsTypeOrder) {
		case OrganizeImportsTypeOrderLast:
			return "last"
		case OrganizeImportsTypeOrderInline:
			return "inline"
		case OrganizeImportsTypeOrderFirst:
			return "first"
		default:
			return "auto"
		}
	},
}

type fieldInfo struct {
	lowerName string   // lowercase Go field name for direct matching
	aliases   []string // additional lowercase names (e.g., "provideprefixandsuffixtextforrename")
	path      string   // dotted path for VS Code config (e.g., "preferences.quoteStyle")
	fieldPath []int    // index path to field in struct
	invert    bool     // whether to invert boolean values
}

var fieldInfoCache = sync.OnceValue(func() []fieldInfo {
	var infos []fieldInfo
	collectFieldInfos(reflect.TypeFor[UserPreferences](), nil, &infos)
	return infos
})

// lowerNameIndex maps lowercase field names to fieldInfo index for O(1) lookup
var lowerNameIndex = sync.OnceValue(func() map[string]int {
	infos := fieldInfoCache()
	index := make(map[string]int, len(infos)*2)
	for i, info := range infos {
		index[info.lowerName] = i
		for _, alias := range info.aliases {
			index[alias] = i
		}
	}
	return index
})

func collectFieldInfos(t reflect.Type, indexPath []int, infos *[]fieldInfo) {
	for i := range t.NumField() {
		field := t.Field(i)
		currentPath := append(slices.Clone(indexPath), i)

		if field.Type.Kind() == reflect.Struct && field.Tag.Get("pref") == "" {
			collectFieldInfos(field.Type, currentPath, infos)
			continue
		}

		info := fieldInfo{
			lowerName: strings.ToLower(field.Name),
			fieldPath: currentPath,
		}

		tag := field.Tag.Get("pref")
		if tag != "" {
			// Parse tag: "path" or "path,invert" or "path,alias:name"
			parts := strings.Split(tag, ",")
			info.path = parts[0]
			for _, part := range parts[1:] {
				if part == "invert" {
					info.invert = true
				} else if alias, ok := strings.CutPrefix(part, "alias:"); ok {
					info.aliases = append(info.aliases, strings.ToLower(alias))
				}
			}
		}

		*infos = append(*infos, info)
	}
}

func getNestedValue(config map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := any(config)
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func setNestedValue(config map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	current := config
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			next = make(map[string]any)
			current[part] = next
		}
		current = next
	}
	current[parts[len(parts)-1]] = value
}

func (p *UserPreferences) parseWorker(config map[string]any) {
	v := reflect.ValueOf(p).Elem()
	infos := fieldInfoCache()
	index := lowerNameIndex()

	// Process "unstable" first - these are spread directly by field name
	if unstable, ok := config["unstable"].(map[string]any); ok {
		for name, value := range unstable {
			if idx, found := index[strings.ToLower(name)]; found {
				info := infos[idx]
				field := getFieldByPath(v, info.fieldPath)
				if info.invert {
					if b, ok := value.(bool); ok {
						value = !b
					}
				}
				setFieldFromValue(field, value)
			}
		}
	}

	// Process path-based config (VS Code style)
	for _, info := range infos {
		if info.path == "" {
			continue
		}
		val, ok := getNestedValue(config, info.path)
		if !ok {
			continue
		}

		field := getFieldByPath(v, info.fieldPath)
		if info.invert {
			if b, ok := val.(bool); ok {
				val = !b
			}
		}
		setFieldFromValue(field, val)
	}

	// Process direct field names at root level (non-VS Code clients, fourslash tests)
	for name, value := range config {
		if name == "unstable" {
			continue
		}
		// Skip known VS Code config sections that are handled by path-based parsing
		switch name {
		case "preferences", "suggest", "inlayHints", "referencesCodeLens",
			"implementationsCodeLens", "workspaceSymbols", "format", "tsserver", "tsc", "experimental":
			continue
		}
		if idx, found := index[strings.ToLower(name)]; found {
			info := infos[idx]
			field := getFieldByPath(v, info.fieldPath)
			if info.invert {
				if b, ok := value.(bool); ok {
					value = !b
				}
			}
			setFieldFromValue(field, value)
		}
	}
}

func getFieldByPath(v reflect.Value, path []int) reflect.Value {
	for _, idx := range path {
		v = v.Field(idx)
	}
	return v
}

func setFieldFromValue(field reflect.Value, val any) {
	if val == nil {
		return
	}

	// Check custom parsers first (for types like Tristate, OrganizeImportsCollation, etc.)
	if parser, ok := typeParsers[field.Type()]; ok {
		field.Set(reflect.ValueOf(parser(val)))
		return
	}

	switch field.Kind() {
	case reflect.Bool:
		if b, ok := val.(bool); ok {
			field.SetBool(b)
		}
	case reflect.Int:
		switch v := val.(type) {
		case int:
			field.SetInt(int64(v))
		case float64:
			field.SetInt(int64(v))
		}
	case reflect.String:
		if s, ok := val.(string); ok {
			field.SetString(s)
		}
	case reflect.Slice:
		if arr, ok := val.([]any); ok {
			result := reflect.MakeSlice(field.Type(), 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = reflect.Append(result, reflect.ValueOf(s))
				}
			}
			field.Set(result)
		}
	}
}

func (p *UserPreferences) MarshalJSONTo(enc *jsontext.Encoder) error {
	config := make(map[string]any)
	v := reflect.ValueOf(p).Elem()

	for _, info := range fieldInfoCache() {
		field := getFieldByPath(v, info.fieldPath)

		val := serializeField(field)
		if val == nil {
			continue
		}
		if info.invert {
			if b, ok := val.(bool); ok {
				val = !b
			}
		}

		// Use the path if available, otherwise use the lowercase field name at root level
		if info.path != "" {
			setNestedValue(config, info.path, val)
		} else {
			config[info.lowerName] = val
		}
	}

	return json.MarshalEncode(enc, config)
}

func serializeField(field reflect.Value) any {
	// Check custom serializers first (for types like Tristate, OrganizeImportsCollation, etc.)
	if serializer, ok := typeSerializers[field.Type()]; ok {
		return serializer(field.Interface())
	}

	switch field.Kind() {
	case reflect.Bool:
		return field.Bool()
	case reflect.Int:
		return int(field.Int())
	case reflect.String:
		return field.String()
	case reflect.Slice:
		if field.IsNil() {
			return nil
		}
		result := make([]string, field.Len())
		for i := range field.Len() {
			result[i] = field.Index(i).String()
		}
		return result
	default:
		return field.Interface()
	}
}

func (p *UserPreferences) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var config map[string]any
	if err := json.UnmarshalDecode(dec, &config); err != nil {
		return err
	}
	// Start with defaults, then overlay parsed values
	*p = *DefaultUserPreferences.Copy()
	p.parseWorker(config)
	return nil
}

// --- Helper methods ---

func deepCopy[T any](src T) T {
	var dst T
	deepCopyValue(reflect.ValueOf(&dst).Elem(), reflect.ValueOf(src))
	return dst
}

func deepCopyValue(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Pointer:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		dst.Set(reflect.New(src.Type().Elem()))
		deepCopyValue(dst.Elem(), src.Elem())
	case reflect.Struct:
		for i := range src.NumField() {
			deepCopyValue(dst.Field(i), src.Field(i))
		}
	case reflect.Slice:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Len()))
		for i := range src.Len() {
			deepCopyValue(dst.Index(i), src.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		for _, key := range src.MapKeys() {
			val := src.MapIndex(key)
			copiedVal := reflect.New(val.Type()).Elem()
			deepCopyValue(copiedVal, val)
			dst.SetMapIndex(key, copiedVal)
		}
	default:
		dst.Set(src)
	}
}

func (p *UserPreferences) Copy() *UserPreferences {
	return deepCopy(p)
}

func (p *UserPreferences) ModuleSpecifierPreferences() modulespecifiers.UserPreferences {
	return modulespecifiers.UserPreferences(p.ModuleSpecifier)
}

// ParseUserPreferences parses user preferences from a config map or returns existing preferences.
// For config maps: returns a fresh *UserPreferences with defaults applied, then overlaid with parsed values.
// For *UserPreferences: returns the same pointer (caller should not mutate).
// Returns nil if item is nil or unrecognized type.
func ParseUserPreferences(item any) *UserPreferences {
	if item == nil {
		return nil
	}
	if config, ok := item.(map[string]any); ok {
		p := DefaultUserPreferences.Copy()
		p.parseWorker(config)
		return p
	}
	if prefs, ok := item.(*UserPreferences); ok {
		return prefs
	}
	return nil
}
