package options

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
)

type mapEntry[K comparable, V any] struct {
	key   K
	value V
}

func NewOrderedMapFromList[K comparable, V any](items []mapEntry[K, V]) *collections.OrderedMap[K, V] {
	mp := collections.NewOrderedMapWithSizeHint[K, V](len(items))
	for _, item := range items {
		mp.Set(item.key, item.value)
	}
	return mp
}

var libMap = NewOrderedMapFromList([]mapEntry[string, string]{
	// JavaScript only
	{"es5", "lib.es5.d.ts"},
	{"es6", "lib.es2015.d.ts"},
	{"es2015", "lib.es2015.d.ts"},
	{"es7", "lib.es2016.d.ts"},
	{"es2016", "lib.es2016.d.ts"},
	{"es2017", "lib.es2017.d.ts"},
	{"es2018", "lib.es2018.d.ts"},
	{"es2019", "lib.es2019.d.ts"},
	{"es2020", "lib.es2020.d.ts"},
	{"es2021", "lib.es2021.d.ts"},
	{"es2022", "lib.es2022.d.ts"},
	{"es2023", "lib.es2023.d.ts"},
	{"es2024", "lib.es2024.d.ts"},
	{"esnext", "lib.esnext.d.ts"},
	// Host only
	{"dom", "lib.dom.d.ts"},
	{"dom.iterable", "lib.dom.iterable.d.ts"},
	{"dom.asynciterable", "lib.dom.asynciterable.d.ts"},
	{"webworker", "lib.webworker.d.ts"},
	{"webworker.importscripts", "lib.webworker.importscripts.d.ts"},
	{"webworker.iterable", "lib.webworker.iterable.d.ts"},
	{"webworker.asynciterable", "lib.webworker.asynciterable.d.ts"},
	{"scripthost", "lib.scripthost.d.ts"},
	// ES2015 Or ESNext By-feature options
	{"es2015.string(core", "lib.es2015.string(core.d.ts"},
	{"es2015.collection", "lib.es2015.collection.d.ts"},
	{"es2015.generator", "lib.es2015.generator.d.ts"},
	{"es2015.iterable", "lib.es2015.iterable.d.ts"},
	{"es2015.promise", "lib.es2015.promise.d.ts"},
	{"es2015.proxy", "lib.es2015.proxy.d.ts"},
	{"es2015.reflect", "lib.es2015.reflect.d.ts"},
	{"es2015.symbol", "lib.es2015.symbol.d.ts"},
	{"es2015.symbol.wellknown", "lib.es2015.symbol.wellknown.d.ts"},
	{"es2016.array.include", "lib.es2016.array.include.d.ts"},
	{"es2016.intl", "lib.es2016.intl.d.ts"},
	{"es2017.arraybuffer", "lib.es2017.arraybuffer.d.ts"},
	{"es2017.date", "lib.es2017.date.d.ts"},
	{"es2017.object", "lib.es2017.object.d.ts"},
	{"es2017.sharedmemory", "lib.es2017.sharedmemory.d.ts"},
	{"es2017.string", "lib.es2017.string.d.ts"},
	{"es2017.intl", "lib.es2017.intl.d.ts"},
	{"es2017.typedarrays", "lib.es2017.typedarrays.d.ts"},
	{"es2018.asyncgenerator", "lib.es2018.asyncgenerator.d.ts"},
	{"es2018.asynciterable", "lib.es2018.asynciterable.d.ts"},
	{"es2018.intl", "lib.es2018.intl.d.ts"},
	{"es2018.promise", "lib.es2018.promise.d.ts"},
	{"es2018.regexp", "lib.es2018.regexp.d.ts"},
	{"es2019.array", "lib.es2019.array.d.ts"},
	{"es2019.object", "lib.es2019.object.d.ts"},
	{"es2019.string", "lib.es2019.string.d.ts"},
	{"es2019.symbol", "lib.es2019.symbol.d.ts"},
	{"es2019.intl", "lib.es2019.intl.d.ts"},
	{"es2020.bigint", "lib.es2020.bigint.d.ts"},
	{"es2020.date", "lib.es2020.date.d.ts"},
	{"es2020.promise", "lib.es2020.promise.d.ts"},
	{"es2020.sharedmemory", "lib.es2020.sharedmemory.d.ts"},
	{"es2020.string", "lib.es2020.string.d.ts"},
	{"es2020.symbol.wellknown", "lib.es2020.symbol.wellknown.d.ts"},
	{"es2020.intl", "lib.es2020.intl.d.ts"},
	{"es2020.number", "lib.es2020.number.d.ts"},
	{"es2021.promise", "lib.es2021.promise.d.ts"},
	{"es2021.string", "lib.es2021.string.d.ts"},
	{"es2021.weakref", "lib.es2021.weakref.d.ts"},
	{"es2021.intl", "lib.es2021.intl.d.ts"},
	{"es2022.array", "lib.es2022.array.d.ts"},
	{"es2022.error", "lib.es2022.error.d.ts"},
	{"es2022.intl", "lib.es2022.intl.d.ts"},
	{"es2022.object", "lib.es2022.object.d.ts"},
	{"es2022.string", "lib.es2022.string.d.ts"},
	{"es2022.regexp", "lib.es2022.regexp.d.ts"},
	{"es2023.array", "lib.es2023.array.d.ts"},
	{"es2023.collection", "lib.es2023.collection.d.ts"},
	{"es2023.intl", "lib.es2023.intl.d.ts"},
	{"es2024.arraybuffer", "lib.es2024.arraybuffer.d.ts"},
	{"es2024.collection", "lib.es2024.collection.d.ts"},
	{"es2024.object", "lib.es2024.object.d.ts"},
	{"es2024.promise", "lib.es2024.promise.d.ts"},
	{"es2024.regexp", "lib.es2024.regexp.d.ts"},
	{"es2024.sharedmemory", "lib.es2024.sharedmemory.d.ts"},
	{"es2024.string", "lib.es2024.string.d.ts"},
	{"esnext.array", "lib.esnext.array.d.ts"},
	{"esnext.collection", "lib.esnext.collection.d.ts"},
	{"esnext.symbol", "lib.es2019.symbol.d.ts"},
	{"esnext.asynciterable", "lib.es2018.asynciterable.d.ts"},
	{"esnext.intl", "lib.esnext.intl.d.ts"},
	{"esnext.disposable", "lib.esnext.disposable.d.ts"},
	{"esnext.bigint", "lib.es2020.bigint.d.ts"},
	{"esnext.string", "lib.es2024.string.d.ts"},
	{"esnext.promise", "lib.es2024.promise.d.ts"},
	{"esnext.weakref", "lib.es2021.weakref.d.ts"},
	{"esnext.decorators", "lib.esnext.decorators.d.ts"},
	{"esnext.object", "lib.es2024.object.d.ts"},
	{"esnext.regexp", "lib.es2024.regexp.d.ts"},
	{"esnext.iterator", "lib.esnext.iterator.d.ts"},
	{"esnext.promise", "lib.esnext.promise.d.ts"},
	{"decorators", "lib.decorators.d.ts"},
	{"decorators.legacy", "lib.decorators.legacy.d.ts"},
})

var moduleResolutionOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"node16", core.ModuleResolutionKindNode16.String()},
	{"nodenext", core.ModuleResolutionKindNodeNext.String()},
	{"bundler", core.ModuleResolutionKindBundler.String()},
})

var targetOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"es3", string(core.ScriptTargetES3)},
	{"es5", string(core.ScriptTargetES5)},
	{"es6", string(core.ScriptTargetES2015)},
	{"es2015", string(core.ScriptTargetES2015)},
	{"es2016", string(core.ScriptTargetES2016)},
	{"es2017", string(core.ScriptTargetES2017)},
	{"es2018", string(core.ScriptTargetES2018)},
	{"es2019", string(core.ScriptTargetES2019)},
	{"es2020", string(core.ScriptTargetES2020)},
	{"es2021", string(core.ScriptTargetES2021)},
	{"es2022", string(core.ScriptTargetES2022)},
	{"es2023", string(core.ScriptTargetES2023)},
	// {"es2024", string(core.ScriptTargetES2024)},
	{"esnext", string(core.ScriptTargetESNext)},
})

var moduleOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"none", string(core.ModuleKindNone)},
	{"commonjs", string(core.ModuleKindCommonJS)},
	{"amd", string(core.ModuleKindAMD)},
	{"system", string(core.ModuleKindSystem)},
	{"umd", string(core.ModuleKindUMD)},
	{"es6", string(core.ModuleKindES2015)},
	{"es2015", string(core.ModuleKindES2015)},
	{"es2020", string(core.ModuleKindES2020)},
	{"es2022", string(core.ModuleKindES2022)},
	{"esnext", string(core.ModuleKindESNext)},
	{"node16", string(core.ModuleKindNode16)},
	{"nodenext", string(core.ModuleKindNodeNext)},
	{"preserve", string(core.ModuleKindPreserve)},
})

// TODO: cleanup??
type ModuleDetectionKind int32

const (
	ModuleDetectionKindAuto   ModuleDetectionKind = 0
	ModuleDetectionKindLegacy ModuleDetectionKind = 1
	ModuleDetectionKindForce  ModuleDetectionKind = 2
)

var moduleDetectionOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"auto", string(ModuleDetectionKindAuto)},
	{"legacy", string(ModuleDetectionKindLegacy)},
	{"force", string(ModuleDetectionKindForce)},
})

var jsxOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"preserve", string(core.JsxEmitPreserve)},
	{"react-native", string(core.JsxEmitReactNative)},
	{"react", string(core.JsxEmitReact)},
	{"react-jsx", string(core.JsxEmitReactJSX)},
	{"react-jsxdev", string(core.JsxEmitReactJSXDev)},
})

var newLineOptionMap = NewOrderedMapFromList([]mapEntry[string, string]{
	{"crlf", string(rune(0))},
	{"lf", string(rune(1))},
})

// TODO: figure out where these definitions should be, or if they're still needed
// crlf: NewLineKind.CarriageReturnLineFeed,
// lf: NewLineKind.LineFeed,
