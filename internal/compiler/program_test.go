package compiler

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type testFile struct {
	FileName string `json:"name"`
	Contents string `json:"contents"`
}

type programTest struct {
	TestName      string            `json:"name"`
	Files         []testFile        `json:"files"`
	ExpectedFiles []string          `json:"expectedFiles"`
	Target        core.ScriptTarget `json:"target"`
}

var esnextLibs = []string{
	"lib.es5.d.ts",
	"lib.es2015.d.ts",
	"lib.es2016.d.ts",
	"lib.es2017.d.ts",
	"lib.es2018.d.ts",
	"lib.es2019.d.ts",
	"lib.es2020.d.ts",
	"lib.es2021.d.ts",
	"lib.es2022.d.ts",
	"lib.es2023.d.ts",
	"lib.esnext.d.ts",
	"lib.dom.d.ts",
	"lib.dom.iterable.d.ts",
	"lib.dom.asynciterable.d.ts",
	"lib.webworker.importscripts.d.ts",
	"lib.scripthost.d.ts",
	"lib.es2015.core.d.ts",
	"lib.es2015.collection.d.ts",
	"lib.es2015.generator.d.ts",
	"lib.es2015.iterable.d.ts",
	"lib.es2015.promise.d.ts",
	"lib.es2015.proxy.d.ts",
	"lib.es2015.reflect.d.ts",
	"lib.es2015.symbol.d.ts",
	"lib.es2015.symbol.wellknown.d.ts",
	"lib.es2016.array.include.d.ts",
	"lib.es2016.intl.d.ts",
	"lib.es2017.date.d.ts",
	"lib.es2017.object.d.ts",
	"lib.es2017.sharedmemory.d.ts",
	"lib.es2017.string.d.ts",
	"lib.es2017.intl.d.ts",
	"lib.es2017.typedarrays.d.ts",
	"lib.es2018.asyncgenerator.d.ts",
	"lib.es2018.asynciterable.d.ts",
	"lib.es2018.intl.d.ts",
	"lib.es2018.promise.d.ts",
	"lib.es2018.regexp.d.ts",
	"lib.es2019.array.d.ts",
	"lib.es2019.object.d.ts",
	"lib.es2019.string.d.ts",
	"lib.es2019.symbol.d.ts",
	"lib.es2019.intl.d.ts",
	"lib.es2020.bigint.d.ts",
	"lib.es2020.date.d.ts",
	"lib.es2020.promise.d.ts",
	"lib.es2020.sharedmemory.d.ts",
	"lib.es2020.string.d.ts",
	"lib.es2020.symbol.wellknown.d.ts",
	"lib.es2020.intl.d.ts",
	"lib.es2020.number.d.ts",
	"lib.es2021.promise.d.ts",
	"lib.es2021.string.d.ts",
	"lib.es2021.weakref.d.ts",
	"lib.es2021.intl.d.ts",
	"lib.es2022.array.d.ts",
	"lib.es2022.error.d.ts",
	"lib.es2022.intl.d.ts",
	"lib.es2022.object.d.ts",
	"lib.es2022.sharedmemory.d.ts",
	"lib.es2022.string.d.ts",
	"lib.es2022.regexp.d.ts",
	"lib.es2023.array.d.ts",
	"lib.es2023.collection.d.ts",
	"lib.es2023.intl.d.ts",
	"lib.esnext.array.d.ts",
	"lib.esnext.collection.d.ts",
	"lib.esnext.intl.d.ts",
	"lib.esnext.disposable.d.ts",
	"lib.esnext.string.d.ts",
	"lib.esnext.promise.d.ts",
	"lib.esnext.decorators.d.ts",
	"lib.esnext.object.d.ts",
	"lib.esnext.regexp.d.ts",
	"lib.esnext.iterator.d.ts",
	"lib.decorators.d.ts",
	"lib.decorators.legacy.d.ts",
	"lib.esnext.full.d.ts",
}

var programTestCases = []programTest{
	{
		TestName: "BasicFileOrdering",
		Files: []testFile{
			{FileName: "c:/dev/src/index.ts", Contents: "/// <reference path='c:/dev/src2/a/5.ts' />\n/// <reference path='c:/dev/src2/a/10.ts' />"},
			{FileName: "c:/dev/src2/a/5.ts", Contents: "/// <reference path='4.ts' />"},
			{FileName: "c:/dev/src2/a/4.ts", Contents: "/// <reference path='b/3.ts' />"},
			{FileName: "c:/dev/src2/a/b/3.ts", Contents: "/// <reference path='2.ts' />"},
			{FileName: "c:/dev/src2/a/b/2.ts", Contents: "/// <reference path='c/1.ts' />"},
			{FileName: "c:/dev/src2/a/b/c/1.ts", Contents: "console.log('hello');"},
			{FileName: "c:/dev/src2/a/10.ts", Contents: "/// <reference path='b/c/d/9.ts' />"},
			{FileName: "c:/dev/src2/a/b/c/d/9.ts", Contents: "/// <reference path='e/8.ts' />"},
			{FileName: "c:/dev/src2/a/b/c/d/e/8.ts", Contents: "/// <reference path='7.ts' />"},
			{FileName: "c:/dev/src2/a/b/c/d/e/7.ts", Contents: "/// <reference path='f/6.ts' />"},
			{FileName: "c:/dev/src2/a/b/c/d/e/f/6.ts", Contents: "console.log('world!');"},
		},
		ExpectedFiles: append(esnextLibs,
			"c:/dev/src2/a/b/c/1.ts",
			"c:/dev/src2/a/b/2.ts",
			"c:/dev/src2/a/b/3.ts",
			"c:/dev/src2/a/4.ts",
			"c:/dev/src2/a/5.ts",
			"c:/dev/src2/a/b/c/d/e/f/6.ts",
			"c:/dev/src2/a/b/c/d/e/7.ts",
			"c:/dev/src2/a/b/c/d/e/8.ts",
			"c:/dev/src2/a/b/c/d/9.ts",
			"c:/dev/src2/a/10.ts",
			"c:/dev/src/index.ts",
		),
		Target: core.ScriptTargetESNext,
	},
	{
		TestName: "FileOrderingImports",
		Files: []testFile{
			{FileName: "c:/dev/src/index.ts", Contents: "import * as five from '../src2/a/5.ts';\nimport * as ten from '../src2/a/10.ts';"},
			{FileName: "c:/dev/src2/a/5.ts", Contents: "import * as four from './4.ts';"},
			{FileName: "c:/dev/src2/a/4.ts", Contents: "import * as three from './b/3.ts';"},
			{FileName: "c:/dev/src2/a/b/3.ts", Contents: "import * as two from './2.ts';"},
			{FileName: "c:/dev/src2/a/b/2.ts", Contents: "import * as one from './c/1.ts';"},
			{FileName: "c:/dev/src2/a/b/c/1.ts", Contents: "console.log('hello');"},
			{FileName: "c:/dev/src2/a/10.ts", Contents: "import * as nine from './b/c/d/9.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/9.ts", Contents: "import * as eight from './e/8.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/8.ts", Contents: "import * as seven from './7.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/7.ts", Contents: "import * as six from './f/6.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/f/6.ts", Contents: "console.log('world!');"},
		},
		ExpectedFiles: append(esnextLibs,
			"c:/dev/src2/a/b/c/1.ts",
			"c:/dev/src2/a/b/2.ts",
			"c:/dev/src2/a/b/3.ts",
			"c:/dev/src2/a/4.ts",
			"c:/dev/src2/a/5.ts",
			"c:/dev/src2/a/b/c/d/e/f/6.ts",
			"c:/dev/src2/a/b/c/d/e/7.ts",
			"c:/dev/src2/a/b/c/d/e/8.ts",
			"c:/dev/src2/a/b/c/d/9.ts",
			"c:/dev/src2/a/10.ts",
			"c:/dev/src/index.ts",
		),
		Target: core.ScriptTargetESNext,
	},
	{
		TestName: "FileOrderingCycles",
		Files: []testFile{
			{FileName: "c:/dev/src/index.ts", Contents: "import * as five from '../src2/a/5.ts';\nimport * as ten from '../src2/a/10.ts';"},
			{FileName: "c:/dev/src2/a/5.ts", Contents: "import * as four from './4.ts';"},
			{FileName: "c:/dev/src2/a/4.ts", Contents: "import * as three from './b/3.ts';"},
			{FileName: "c:/dev/src2/a/b/3.ts", Contents: "import * as two from './2.ts';\nimport * as cycle from 'c:/dev/src/index.ts'; "},
			{FileName: "c:/dev/src2/a/b/2.ts", Contents: "import * as one from './c/1.ts';"},
			{FileName: "c:/dev/src2/a/b/c/1.ts", Contents: "console.log('hello');"},
			{FileName: "c:/dev/src2/a/10.ts", Contents: "import * as nine from './b/c/d/9.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/9.ts", Contents: "import * as eight from './e/8.ts';\nimport * as cycle from 'c:/dev/src/index.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/8.ts", Contents: "import * as seven from './7.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/7.ts", Contents: "import * as six from './f/6.ts';"},
			{FileName: "c:/dev/src2/a/b/c/d/e/f/6.ts", Contents: "console.log('world!');"},
		},
		ExpectedFiles: append(esnextLibs,
			"c:/dev/src2/a/b/c/1.ts",
			"c:/dev/src2/a/b/2.ts",
			"c:/dev/src2/a/b/3.ts",
			"c:/dev/src2/a/4.ts",
			"c:/dev/src2/a/5.ts",
			"c:/dev/src2/a/b/c/d/e/f/6.ts",
			"c:/dev/src2/a/b/c/d/e/7.ts",
			"c:/dev/src2/a/b/c/d/e/8.ts",
			"c:/dev/src2/a/b/c/d/9.ts",
			"c:/dev/src2/a/10.ts",
			"c:/dev/src/index.ts",
		),
		Target: core.ScriptTargetESNext,
	},
}

func TestProgram(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	for _, testCase := range programTestCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			libPrefix := bundled.LibPath() + "/"
			fs := vfstest.FromMapFS(fstest.MapFS{}, false /*useCaseSensitiveFileNames*/)
			fs = bundled.WrapFS(fs)

			for _, testFile := range testCase.Files {
				_ = fs.WriteFile(testFile.FileName, testFile.Contents, false)
			}

			opts := core.CompilerOptions{Target: testCase.Target}

			program := NewProgram(ProgramOptions{
				RootPath:           "c:/dev/src",
				Host:               NewCompilerHost(&opts, "c:/dev/src", fs),
				Options:            &opts,
				DefaultLibraryPath: bundled.LibPath(),
				SingleThreaded:     false,
			})

			actualFiles := []string{}
			for _, file := range program.files {
				actualFiles = append(actualFiles, strings.TrimPrefix(file.FileName(), libPrefix))
			}

			assert.DeepEqual(t, testCase.ExpectedFiles, actualFiles)
		})
	}
}
