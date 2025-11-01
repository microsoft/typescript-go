package lspservertests

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type projectInfo = *compiler.Program

type openFileInfo struct {
	defaultProjectName string
	allProjects        []string
}

type diffTableOptions struct {
	indent   string
	sortKeys bool
}

type diffTable struct {
	diff    collections.OrderedMap[string, string]
	options diffTableOptions
}

func (d *diffTable) add(key, value string) {
	d.diff.Set(key, value)
}

func (d *diffTable) print(w io.Writer, header string) {
	count := d.diff.Size()
	if count == 0 {
		return
	}
	if header != "" {
		fmt.Fprintf(w, "%s%s\n", d.options.indent, header)
	}
	diffKeys := make([]string, 0, count)
	keyWidth := 0
	indent := d.options.indent + "  "
	for key := range d.diff.Keys() {
		keyWidth = max(keyWidth, len(key))
		diffKeys = append(diffKeys, key)
	}
	if d.options.sortKeys {
		slices.Sort(diffKeys)
	}

	for _, key := range diffKeys {
		value := d.diff.GetOrZero(key)
		fmt.Fprintf(w, "%s%-*s %s\n", indent, keyWidth+1, key, value)
	}
}

type diffTableWriter struct {
	hasChange bool
	header    string
	diffs     map[string]func(io.Writer)
}

func newDiffTableWriter(header string) *diffTableWriter {
	return &diffTableWriter{header: header, diffs: make(map[string]func(io.Writer))}
}

func (d *diffTableWriter) setHasChange() {
	d.hasChange = true
}

func (d *diffTableWriter) add(key string, fn func(io.Writer)) {
	d.diffs[key] = fn
}

func (d *diffTableWriter) print(w io.Writer) {
	if d.hasChange {
		fmt.Fprintf(w, "%s::\n", d.header)
		keys := slices.Collect(maps.Keys(d.diffs))
		slices.Sort(keys)
		for _, key := range keys {
			d.diffs[key](w)
		}
	}
}

func areIterSeqEqual(a, b iter.Seq[tspath.Path]) bool {
	aSlice := slices.Collect(a)
	bSlice := slices.Collect(b)
	slices.Sort(aSlice)
	slices.Sort(bSlice)
	return slices.Equal(aSlice, bSlice)
}

func printSlicesWithDiffTable(w io.Writer, header string, newSlice []string, getOldSlice func() []string, options diffTableOptions, topChange string, isDefault func(entry string) bool) {
	var oldSlice []string
	if topChange == "*modified*" {
		oldSlice = getOldSlice()
	}
	table := diffTable{options: options}
	for _, entry := range newSlice {
		entryChange := ""
		if isDefault != nil && isDefault(entry) {
			entryChange = "(default) "
		}
		if topChange == "*modified*" && !slices.Contains(oldSlice, entry) {
			entryChange = "*new*"
		}
		table.add(entry, entryChange)
	}
	if topChange == "*modified*" {
		for _, entry := range oldSlice {
			if !slices.Contains(newSlice, entry) {
				table.add(entry, "*deleted*")
			}
		}
	}
	table.print(w, header)
}

func sliceFromIterSeqPath(seq iter.Seq[tspath.Path]) []string {
	var result []string
	for path := range seq {
		result = append(result, string(path))
	}
	slices.Sort(result)
	return result
}

func printPathIterSeqWithDiffTable(w io.Writer, header string, newIterSeq iter.Seq[tspath.Path], getOldIterSeq func() iter.Seq[tspath.Path], options diffTableOptions, topChange string) {
	printSlicesWithDiffTable(
		w,
		header,
		sliceFromIterSeqPath(newIterSeq),
		func() []string { return sliceFromIterSeqPath(getOldIterSeq()) },
		options,
		topChange,
		nil,
	)
}

type stateDiff struct {
	server   *testServer
	snapshot *project.Snapshot

	currentProjects map[string]projectInfo
	projectDiffs    *diffTableWriter

	currentOpenFiles map[string]*openFileInfo
	filesDiff        *diffTableWriter

	configFileRegistry   *project.ConfigFileRegistry
	configDiffs          *diffTableWriter
	configFileNamesDiffs *diffTableWriter
}

func printStateDiff(server *testServer, w io.Writer) {
	if !server.isInitialized {
		return
	}
	snapshot, release := server.server.Session().Snapshot()
	defer release()

	stateDiff := &stateDiff{
		server:               server,
		snapshot:             snapshot,
		configFileRegistry:   snapshot.ProjectCollection.ConfigFileRegistry(),
		currentProjects:      make(map[string]projectInfo),
		currentOpenFiles:     make(map[string]*openFileInfo),
		projectDiffs:         newDiffTableWriter("Projects"),
		filesDiff:            newDiffTableWriter("Open Files"),
		configDiffs:          newDiffTableWriter("Config"),
		configFileNamesDiffs: newDiffTableWriter("Config File Names"),
	}
	stateDiff.checkProjects()
	stateDiff.checkOpenFiles()
	stateDiff.checkConfigFileRegistry()
	stateDiff.print(w)
}

func (d *stateDiff) checkProjects() {
	d.server.t.Helper()
	options := diffTableOptions{indent: "  "}
	for _, project := range d.snapshot.ProjectCollection.Projects() {
		program := project.GetProgram()
		var oldProgram *compiler.Program
		d.currentProjects[project.Name()] = program
		projectChange := ""
		if existing, ok := d.server.serializedProjects[project.Name()]; ok {
			oldProgram = existing
			if oldProgram != program {
				projectChange = "*modified*"
				d.projectDiffs.setHasChange()
			} else {
				projectChange = ""
			}
		} else {
			projectChange = "*new*"
			d.projectDiffs.setHasChange()
		}

		d.projectDiffs.add(project.Name(), func(w io.Writer) {
			fmt.Fprintf(w, "  [%s] %s\n", project.Name(), projectChange)
			subDiff := diffTable{options: options}
			if program != nil {
				for _, file := range program.GetSourceFiles() {
					fileDiff := ""
					// No need to write "*new*" for files as its obvious
					fileName := file.FileName()
					if projectChange == "*modified*" {
						if oldProgram == nil {
							if !d.server.isLibFile(fileName) {
								fileDiff = "*new*"
							}
						} else if oldFile := oldProgram.GetSourceFileByPath(file.Path()); oldFile == nil {
							fileDiff = "*new*"
						} else if oldFile != file {
							fileDiff = "*modified*"
						}
					}
					if fileDiff != "" || !d.server.isLibFile(fileName) {
						subDiff.add(fileName, fileDiff)
					}
				}
			}
			if oldProgram != program && oldProgram != nil {
				for _, file := range oldProgram.GetSourceFiles() {
					if program == nil || program.GetSourceFileByPath(file.Path()) == nil {
						subDiff.add(file.FileName(), "*deleted*")
					}
				}
			}
			subDiff.print(w, "")
		})
	}

	for projectName, info := range d.server.serializedProjects {
		if _, found := d.currentProjects[projectName]; !found {
			d.projectDiffs.setHasChange()
			d.projectDiffs.add(projectName, func(w io.Writer) {
				fmt.Fprintf(w, "  [%s] *deleted*\n", projectName)
				subDiff := diffTable{options: options}
				if info != nil {
					for _, file := range info.GetSourceFiles() {
						if fileName := file.FileName(); !d.server.isLibFile(fileName) {
							subDiff.add(fileName, "")
						}
					}
				}
				subDiff.print(w, "")
			})
		}
	}
	d.server.serializedProjects = d.currentProjects
}

func (d *stateDiff) checkOpenFiles() {
	d.server.t.Helper()
	options := diffTableOptions{indent: "  ", sortKeys: true}
	for fileName := range d.server.openFiles {
		path := tspath.ToPath(fileName, "/", d.server.server.FS.UseCaseSensitiveFileNames())
		defaultProject := d.snapshot.ProjectCollection.GetDefaultProject(fileName, path)
		newFileInfo := &openFileInfo{}
		if defaultProject != nil {
			newFileInfo.defaultProjectName = defaultProject.Name()
		}
		for _, project := range d.snapshot.ProjectCollection.Projects() {
			if program := project.GetProgram(); program != nil && program.GetSourceFileByPath(path) != nil {
				newFileInfo.allProjects = append(newFileInfo.allProjects, project.Name())
			}
		}
		slices.Sort(newFileInfo.allProjects)
		d.currentOpenFiles[fileName] = newFileInfo
		openFileChange := ""
		var oldFileInfo *openFileInfo
		if existing, ok := d.server.serializedOpenFiles[fileName]; ok {
			oldFileInfo = existing
			if existing.defaultProjectName != newFileInfo.defaultProjectName || !slices.Equal(existing.allProjects, newFileInfo.allProjects) {
				openFileChange = "*modified*"
				d.filesDiff.setHasChange()
			} else {
				openFileChange = ""
			}
		} else {
			openFileChange = "*new*"
			d.filesDiff.setHasChange()
		}

		d.filesDiff.add(fileName, func(w io.Writer) {
			fmt.Fprintf(w, "  [%s] %s\n", fileName, openFileChange)
			printSlicesWithDiffTable(
				w,
				"",
				newFileInfo.allProjects,
				func() []string { return oldFileInfo.allProjects },
				options,
				openFileChange,
				func(projectName string) bool { return projectName == newFileInfo.defaultProjectName },
			)
		})
	}
	for fileName := range d.server.serializedOpenFiles {
		if _, found := d.currentOpenFiles[fileName]; !found {
			d.filesDiff.setHasChange()
			d.filesDiff.add(fileName, func(w io.Writer) {
				fmt.Fprintf(w, "  [%s] *closed*\n", fileName)
			})
		}
	}
	d.server.serializedOpenFiles = d.currentOpenFiles
}

func (d *stateDiff) checkConfigFileRegistry() {
	if d.server.serializedConfigFileRegistry == d.configFileRegistry {
		return
	}
	options := diffTableOptions{indent: "    ", sortKeys: true}
	d.configFileRegistry.ForEachTestConfigEntry(func(path tspath.Path, entry *project.TestConfigEntry) {
		configChange := ""
		oldEntry := d.server.serializedConfigFileRegistry.GetTestConfigEntry(path)
		if oldEntry == nil {
			configChange = "*new*"
			d.configDiffs.setHasChange()
		} else if oldEntry != entry {
			if !areIterSeqEqual(oldEntry.RetainingProjects, entry.RetainingProjects) ||
				!areIterSeqEqual(oldEntry.RetainingOpenFiles, entry.RetainingOpenFiles) ||
				!areIterSeqEqual(oldEntry.RetainingConfigs, entry.RetainingConfigs) {
				configChange = "*modified*"
				d.configDiffs.setHasChange()
			}
		}
		d.configDiffs.add(string(path), func(w io.Writer) {
			fmt.Fprintf(w, "  [%s] %s\n", entry.FileName, configChange)
			// Print the details of the config entry
			var retainingProjectsModified string
			var retainingOpenFilesModified string
			var retainingConfigsModified string
			if configChange == "*modified*" {
				if !areIterSeqEqual(entry.RetainingProjects, oldEntry.RetainingProjects) {
					retainingProjectsModified = " *modified*"
				}
				if !areIterSeqEqual(entry.RetainingOpenFiles, oldEntry.RetainingOpenFiles) {
					retainingOpenFilesModified = " *modified*"
				}
				if !areIterSeqEqual(entry.RetainingConfigs, oldEntry.RetainingConfigs) {
					retainingConfigsModified = " *modified*"
				}
			}
			printPathIterSeqWithDiffTable(w, "RetainingProjects:"+retainingProjectsModified, entry.RetainingProjects, func() iter.Seq[tspath.Path] { return oldEntry.RetainingProjects }, options, configChange)
			printPathIterSeqWithDiffTable(w, "RetainingOpenFiles:"+retainingOpenFilesModified, entry.RetainingOpenFiles, func() iter.Seq[tspath.Path] { return oldEntry.RetainingOpenFiles }, options, configChange)
			printPathIterSeqWithDiffTable(w, "RetainingConfigs:"+retainingConfigsModified, entry.RetainingConfigs, func() iter.Seq[tspath.Path] { return oldEntry.RetainingConfigs }, options, configChange)
		})
	})
	d.configFileRegistry.ForEachTestConfigFileNamesEntry(func(path tspath.Path, entry *project.TestConfigFileNamesEntry) {
		configFileNamesChange := ""
		oldEntry := d.server.serializedConfigFileRegistry.GetTestConfigFileNamesEntry(path)
		if oldEntry == nil {
			configFileNamesChange = "*new*"
			d.configFileNamesDiffs.setHasChange()
		} else if oldEntry.NearestConfigFileName != entry.NearestConfigFileName ||
			!maps.Equal(oldEntry.Ancestors, entry.Ancestors) {
			configFileNamesChange = "*modified*"
			d.configFileNamesDiffs.setHasChange()
		}
		d.configFileNamesDiffs.add(string(path), func(w io.Writer) {
			fmt.Fprintf(w, "  [%s] %s\n", path, configFileNamesChange)
			var nearestConfigFileNameModified string
			var ancestorDiffModified string
			if configFileNamesChange == "*modified*" {
				if oldEntry.NearestConfigFileName != entry.NearestConfigFileName {
					nearestConfigFileNameModified = " *modified*"
				}
				if !maps.Equal(oldEntry.Ancestors, entry.Ancestors) {
					ancestorDiffModified = " *modified*"
				}
			}
			fmt.Fprintf(w, "    NearestConfigFileName: %s%s\n", entry.NearestConfigFileName, nearestConfigFileNameModified)
			ancestorDiff := diffTable{options: options}
			for config, ancestorOfConfig := range entry.Ancestors {
				ancestorChange := ""
				if configFileNamesChange == "*modified*" {
					if oldConfigFileName, ok := oldEntry.Ancestors[config]; ok {
						if oldConfigFileName != ancestorOfConfig {
							ancestorChange = "*modified*"
						}
					} else {
						ancestorChange = "*new*"
					}
				}
				ancestorDiff.add(config, fmt.Sprintf("%s %s", ancestorOfConfig, ancestorChange))
			}
			if configFileNamesChange == "*modified*" {
				for ancestorPath, oldConfigFileName := range oldEntry.Ancestors {
					if _, ok := entry.Ancestors[ancestorPath]; !ok {
						ancestorDiff.add(ancestorPath, oldConfigFileName+" *deleted*")
					}
				}
			}
			ancestorDiff.print(w, "Ancestors:"+ancestorDiffModified)
		})
	})

	d.server.serializedConfigFileRegistry.ForEachTestConfigEntry(func(path tspath.Path, entry *project.TestConfigEntry) {
		if d.configFileRegistry.GetTestConfigEntry(path) == nil {
			d.configDiffs.setHasChange()
			d.configDiffs.add(string(path), func(w io.Writer) {
				fmt.Fprintf(w, "  [%s] *deleted*\n", entry.FileName)
			})
		}
	})
	d.server.serializedConfigFileRegistry.ForEachTestConfigFileNamesEntry(func(path tspath.Path, entry *project.TestConfigFileNamesEntry) {
		if d.configFileRegistry.GetTestConfigFileNamesEntry(path) == nil {
			d.configFileNamesDiffs.setHasChange()
			d.configFileNamesDiffs.add(string(path), func(w io.Writer) {
				fmt.Fprintf(w, "  [%s] *deleted*\n", path)
			})
		}
	})
	d.server.serializedConfigFileRegistry = d.configFileRegistry
}

func (d *stateDiff) print(w io.Writer) {
	d.projectDiffs.print(w)
	d.filesDiff.print(w)
	d.configDiffs.print(w)
	d.configFileNamesDiffs.print(w)
}
