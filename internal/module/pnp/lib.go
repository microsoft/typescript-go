// BSD 2-Clause License
//
// Copyright (c) 2016-present, Yarn Contributors.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package pnp

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ResolutionKind int

const (
	ResolutionSkipped ResolutionKind = iota
	ResolutionResolved
)

type PNPResolutionHost struct {
	FindPNPManifest func(start string) (*Manifest, error)
}

type ResolutionConfig struct {
	Host PNPResolutionHost
}

type Resolution struct {
	Kind       ResolutionKind
	Path       string
	ModulePath *string
}

func isNodeJSBuiltin(name string) bool {
	return core.NodeCoreModules()[name]
}

func parseScopedPackageName(spec string) (pkg string, sub *string, ok bool) {
	parts := strings.SplitN(spec, "/", 3)
	if len(parts) < 2 {
		return "", nil, false
	}
	scope, name := parts[0], parts[1]
	if scope == "" || name == "" {
		return "", nil, false
	}
	pkg = scope + "/" + name
	if len(parts) == 3 {
		s := parts[2]
		sub = &s
	}
	return pkg, sub, true
}

func parseGlobalPackageName(spec string) (pkg string, sub *string, ok bool) {
	parts := strings.SplitN(spec, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", nil, false
	}
	pkg = parts[0]
	if len(parts) == 2 {
		s := parts[1]
		sub = &s
	}
	return pkg, sub, true
}

func ParseBareIdentifier(spec string) (string, *string, error) {
	var (
		pkg string
		sub *string
		ok  bool
	)
	if strings.HasPrefix(spec, "@") {
		pkg, sub, ok = parseScopedPackageName(spec)
	} else {
		pkg, sub, ok = parseGlobalPackageName(spec)
	}
	if !ok {
		return "", nil, errors.New(diagnostics.Invalid_specifier.Format())
	}
	return pkg, sub, nil
}

func FindClosestPNPManifestPath(start string) (string, bool) {
	dir := tspath.NormalizePath(start)

	for {
		candidate := tspath.CombinePaths(dir, ".pnp.cjs")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, true
		}
		parent := tspath.GetDirectoryPath(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

func extractJSONFromPnPFile(path string, content string) (string, error) {
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: tspath.NormalizePath(path),
		Path:     tspath.Path(path),
	}, content, core.ScriptKindJS)
	if sourceFile == nil {
		return "", errors.New(diagnostics.We_failed_to_locate_the_PnP_data_payload_inside_its_manifest_file_Did_you_manually_edit_the_file.Format())
	}

	var jsonData string
	var found bool

	var visitNode func(*ast.Node) bool
	visitNode = func(node *ast.Node) bool {
		if node == nil {
			return false
		}

		// Look for variable declaration lists with RAW_RUNTIME_STATE (Yarn v4)
		if node.Kind == ast.KindVariableDeclarationList {
			declList := node.AsVariableDeclarationList()
			if declList != nil && declList.Declarations != nil && len(declList.Declarations.Nodes) > 0 {
				declarator := declList.Declarations.Nodes[0]
				if declarator != nil && declarator.Name() != nil {
					if name := declarator.Name().AsIdentifier(); name != nil && name.Text == "RAW_RUNTIME_STATE" {
						if declarator.Initializer() != nil {
							if init := declarator.Initializer().AsStringLiteral(); init != nil {
								jsonData = init.Text
								found = true
								return true
							}
						}
					}
				}
			}
		}

		// Look for call expressions with hydrateRuntimeState(JSON.parse('...')) (Yarn v3)
		if node.Kind == ast.KindCallExpression {
			call := node.AsCallExpression()
			if call != nil && call.Expression != nil {
				if expr := call.Expression.AsIdentifier(); expr != nil && expr.Text == "hydrateRuntimeState" {
					if call.Arguments != nil && len(call.Arguments.Nodes) > 0 && call.Arguments.Nodes[0] != nil {
						if arg := call.Arguments.Nodes[0].AsCallExpression(); arg != nil && arg.Expression != nil {
							if argExpr := arg.Expression.AsPropertyAccessExpression(); argExpr != nil {
								if obj := argExpr.Expression.AsIdentifier(); obj != nil && obj.Text == "JSON" {
									if prop := argExpr.Name().AsIdentifier(); prop != nil && prop.Text == "parse" {
										if arg.Arguments != nil && len(arg.Arguments.Nodes) > 0 && arg.Arguments.Nodes[0] != nil {
											if jsonArg := arg.Arguments.Nodes[0].AsStringLiteral(); jsonArg != nil {
												jsonData = jsonArg.Text
												found = true
												return true
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		return node.ForEachChild(visitNode)
	}

	sourceFile.ForEachChild(visitNode)

	if !found {
		return "", errors.New(diagnostics.We_failed_to_locate_the_PnP_data_payload_inside_its_manifest_file_Did_you_manually_edit_the_file.Format())
	}

	return jsonData, nil
}

func InitPNPManifest(m *Manifest, manifestPath string) error {
	m.ManifestPath = manifestPath
	m.ManifestDir = tspath.GetDirectoryPath(manifestPath)

	m.LocationTrie = *NewLocationTrie[PackageLocator]()

	for name, ranges := range m.PackageRegistryData {
		for reference, info := range ranges {
			joined := tspath.CombinePaths(m.ManifestDir, info.PackageLocation)

			norm := tspath.NormalizePath(joined)

			info.PackageLocation = norm
			ranges[reference] = info

			if !info.DiscardFromLookup {
				m.LocationTrie.Insert(
					info.PackageLocation,
					PackageLocator{Name: name, Reference: reference},
				)
			}
		}
	}

	ranges, ok := m.PackageRegistryData[""]
	if !ok {
		return errors.New(diagnostics.X_assertion_failed_Colon_should_have_a_top_level_name_key.Format())
	}
	top, ok := ranges[""]
	if !ok {
		return errors.New(diagnostics.X_assertion_failed_Colon_should_have_a_top_level_range_key.Format())
	}

	if m.FallbackPool == nil {
		m.FallbackPool = make(FallbackPool)
	}
	for depName, dep := range top.PackageDependencies {
		if _, exists := m.FallbackPool[depName]; !exists {
			m.FallbackPool[depName] = dep
		}
	}

	return nil
}

func LoadPNPManifest(p string) (Manifest, error) {
	content, err := os.ReadFile(p)
	if err != nil {
		return Manifest{}, errors.New(diagnostics.We_failed_to_read_the_content_of_the_manifest_Colon_0.Format(err.Error()))
	}

	jsonData, err := extractJSONFromPnPFile(p, string(content))
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err = json.Unmarshal([]byte(jsonData), &manifest); err != nil {
		return Manifest{}, errors.New(diagnostics.We_failed_to_parse_the_PnP_data_payload_as_proper_JSON_Did_you_manually_edit_the_file_Colon_0.Format(err.Error()))
	}

	err = InitPNPManifest(&manifest, p)
	if err != nil {
		return Manifest{}, errors.New(diagnostics.We_failed_to_init_the_PnP_manifest_Colon_0.Format(err.Error()))
	}
	return manifest, nil
}

func FindPNPManifest(parent string) (*Manifest, error) {
	path, ok := FindClosestPNPManifestPath(parent)
	if !ok {
		return nil, nil
	}
	manifest, err := LoadPNPManifest(path)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func IsDependencyTreeRoot(m *Manifest, loc *PackageLocator) bool {
	return slices.Contains(m.DependencyTreeRoots, PackageLocator{
		Name:      loc.Name,
		Reference: loc.Reference,
	})
}

func FindLocator(manifest *Manifest, path string) *PackageLocator {
	rel := tspath.GetRelativePathFromDirectory(manifest.ManifestDir, path, tspath.ComparePathsOptions{UseCaseSensitiveFileNames: true})

	if manifest.IgnorePatternData != nil {
		_, err := manifest.IgnorePatternData.reg.MatchString(tspath.NormalizePath(rel))
		if err != nil {
			return nil
		}
	}

	normFull := tspath.NormalizePath(path)
	if v, ok := manifest.LocationTrie.GetAncestorValue(normFull); ok && v != nil {
		return v
	}
	return nil
}

func GetPackage(manifest *Manifest, locator *PackageLocator) (*PackageInformation, error) {
	refs, ok := manifest.PackageRegistryData[locator.Name]
	if !ok {
		return nil, errors.New(diagnostics.X_should_have_an_entry_in_the_package_registry_for_0.Format(locator.Name))
	}
	info, ok := refs[locator.Reference]
	if !ok {
		return nil, errors.New(diagnostics.X_should_have_an_entry_in_the_package_registry_for_0.Format(locator.Reference))
	}
	return &info, nil
}

func IsExcludedFromFallback(manifest *Manifest, locator *PackageLocator) bool {
	if refs, ok := manifest.FallbackExclusionList[locator.Name]; ok {
		for _, r := range refs {
			if r == locator.Reference {
				return true
			}
		}
	}
	return false
}

func FindBrokenPeerDependencies(specifier string, parent *PackageLocator) []PackageLocator {
	return []PackageLocator{}
}

func viaSuffix(specifier string, ident string) string {
	if ident != specifier {
		return diagnostics.X_via_0.Format(ident)
	}
	return ""
}

func ResolveToUnqualifiedViaManifest(
	manifest *Manifest,
	specifier string,
	parentPath string,
) (Resolution, error) {
	ident, modulePath, err := ParseBareIdentifier(specifier)
	if err != nil {
		return Resolution{}, err
	}

	parentLocator := FindLocator(manifest, parentPath)
	if parentLocator == nil {
		return Resolution{Kind: ResolutionSkipped}, nil
	}

	parentPkg, err := GetPackage(manifest, parentLocator)
	if err != nil {
		return Resolution{}, err
	}

	var refOrAlias *PackageDependency

	if dep, ok := parentPkg.PackageDependencies[ident]; ok && dep != nil {
		refOrAlias = dep
	}

	if refOrAlias == nil && manifest.EnableTopLevelFallback && !IsExcludedFromFallback(manifest, parentLocator) {
		if dep, ok := manifest.FallbackPool[ident]; ok {
			refOrAlias = dep
		}
	}

	if refOrAlias == nil {
		if isNodeJSBuiltin(specifier) {
			if IsDependencyTreeRoot(manifest, parentLocator) {
				return Resolution{}, errors.New(diagnostics.Your_application_tried_to_access_0_While_this_module_is_usually_interpreted_as_a_Node_builtin_your_resolver_is_running_inside_a_non_Node_resolution_context_where_such_builtins_are_ignored_Since_0_isn_t_otherwise_declared_in_your_dependencies_this_makes_the_require_call_ambiguous_and_unsound_Required_package_Colon_0_1_Required_by_Colon_2.Format(ident, ident, viaSuffix(specifier, ident), parentPath))
			}
			return Resolution{}, errors.New(diagnostics.X_0_tried_to_access_1_While_this_module_is_usually_interpreted_as_a_Node_builtin_your_resolver_is_running_inside_a_non_Node_resolution_context_where_such_builtins_are_ignored_Since_1_isn_t_otherwise_declared_in_0_s_dependencies_this_makes_the_require_call_ambiguous_and_unsound_Required_package_Colon_1_2_Required_by_Colon_3.Format(parentLocator.Name, ident, ident, parentLocator.Name, viaSuffix(specifier, ident), parentPath))
		}

		if IsDependencyTreeRoot(manifest, parentLocator) {
			return Resolution{}, errors.New(diagnostics.Your_application_tried_to_access_0_but_it_isn_t_declared_in_your_dependencies_this_makes_the_require_call_ambiguous_and_unsound_Required_package_Colon_0_1_Required_by_Colon_2.Format(ident, viaSuffix(specifier, ident), parentPath))
		}

		brokenAncestors := FindBrokenPeerDependencies(specifier, parentLocator)
		allBrokenAreRoots := len(brokenAncestors) > 0
		if allBrokenAreRoots {
			for _, brokenAncestor := range brokenAncestors {
				if !IsDependencyTreeRoot(manifest, &brokenAncestor) {
					allBrokenAreRoots = false
					break
				}
			}
		}

		if len(brokenAncestors) > 0 && allBrokenAreRoots {
			return Resolution{}, errors.New(diagnostics.X_0_tried_to_access_1_a_peer_dependency_but_it_isn_t_provided_by_your_application_this_makes_the_require_call_ambiguous_and_unsound_Required_package_Colon_1_2_Required_by_Colon_0_3_via_4.Format(parentLocator.Name, ident, viaSuffix(specifier, ident), parentLocator.Reference, parentPath))
		} else {
			return Resolution{}, errors.New(diagnostics.X_0_tried_to_access_1_a_peer_dependency_but_it_isn_t_provided_by_its_ancestors_this_makes_the_require_call_ambiguous_and_unsound_Required_package_Colon_1_2_Required_by_Colon_0_3_via_4.Format(parentLocator.Name, ident, viaSuffix(specifier, ident), parentLocator.Reference, parentPath))
		}
	}

	var target PackageLocator
	if refOrAlias.IsAlias {
		target = PackageLocator{
			Name:      refOrAlias.Alias[0],
			Reference: refOrAlias.Alias[1],
		}
	} else {
		target = PackageLocator{
			Name:      ident,
			Reference: refOrAlias.Reference,
		}
	}
	depPkg, err := GetPackage(manifest, &target)
	if err != nil {
		return Resolution{}, err
	}

	return Resolution{
		Kind:       ResolutionResolved,
		Path:       depPkg.PackageLocation,
		ModulePath: modulePath,
	}, nil
}

func ResolveToUnqualified(specifier, parentPath string, cfg *ResolutionConfig) (Resolution, error) {
	if cfg == nil || cfg.Host.FindPNPManifest == nil {
		return Resolution{}, errors.New(diagnostics.X_no_host_configured.Format())
	}
	m, err := cfg.Host.FindPNPManifest(parentPath)
	if err != nil {
		return Resolution{}, err
	}
	if m == nil {
		return Resolution{Kind: ResolutionSkipped}, nil
	}
	return ResolveToUnqualifiedViaManifest(m, specifier, parentPath)
}
