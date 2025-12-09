package lsp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/scanner"
)

var (
	//nolint
	OutdatedProjectVersionError = errors.New("OutdatedTypeCheckerIdException")
	//nolint
	OutdatedProjectIdError = errors.New("OutdatedProjectIdException")
)

var (
	nextProjectId   = 1
	project2IdMap   = make(map[*project.Project]int)
	id2ProjectMap   = make(map[int]*project.Project)
	projectMapMutex = &sync.Mutex{}
	projectCacheMap = make(map[int]*ProjectCache)
)

type ProjectCache struct {
	projectVersion   uint64
	requestedTypeIds *collections.Set[checker.TypeId]
	seenTypeIds      map[checker.TypeId]*checker.Type
	seenSymbolIds    map[ast.SymbolId]*ast.Symbol
}

func getProjectId(project *project.Project) int {
	projectMapMutex.Lock()
	defer projectMapMutex.Unlock()
	result, ok := project2IdMap[project]
	if !ok {
		result = nextProjectId
		nextProjectId++
		project2IdMap[project] = result
		id2ProjectMap[result] = project
	}
	return result
}

func getProject(projectId int) (*project.Project, bool) {
	projectMapMutex.Lock()
	defer projectMapMutex.Unlock()
	id, ok := id2ProjectMap[projectId]
	return id, ok
}

func getProjectCache(projectId int, projectVersion uint64) *ProjectCache {
	projectMapMutex.Lock()
	defer projectMapMutex.Unlock()
	result, ok := projectCacheMap[projectId]
	if !ok || result.projectVersion != projectVersion {
		result = &ProjectCache{
			projectVersion:   projectVersion,
			requestedTypeIds: collections.NewSetWithSizeHint[checker.TypeId](10),
			seenTypeIds:      make(map[checker.TypeId]*checker.Type),
			seenSymbolIds:    make(map[ast.SymbolId]*ast.Symbol),
		}
		projectCacheMap[projectId] = result
	}
	return result
}

func CleanupProjectsCache(openedProjects []*project.Project, logger logging.Logger) {
	openedProjectsSet := collections.NewSetWithSizeHint[*project.Project](len(openedProjects))
	for _, p := range openedProjects {
		if p != nil {
			openedProjectsSet.Add(p)
		}
	}

	removedIds := cleanupProjectsCacheImpl(openedProjectsSet)

	if removedIds.Len() > 0 && len(openedProjects) > 0 && openedProjects[0] != nil {
		ids := make([]int, 0, removedIds.Len())
		for id := range removedIds.Keys() {
			ids = append(ids, id)
		}
		logger.Logf("CleanupProjectsCache: removed project ids: %v", ids)
	}
}

func cleanupProjectsCacheImpl(openedProjectsSet *collections.Set[*project.Project]) *collections.Set[int] {
	projectMapMutex.Lock()
	defer projectMapMutex.Unlock()

	removedIds := collections.NewSetWithSizeHint[int](len(id2ProjectMap) + len(projectCacheMap))

	for p := range project2IdMap {
		if !openedProjectsSet.Has(p) {
			delete(project2IdMap, p)
		}
	}

	validProjectIdsSet := collections.NewSetWithSizeHint[int](len(project2IdMap))
	for _, id := range project2IdMap {
		validProjectIdsSet.Add(id)
	}

	for id := range id2ProjectMap {
		if !validProjectIdsSet.Has(id) {
			delete(id2ProjectMap, id)
			removedIds.Add(id)
		}
	}

	for id := range projectCacheMap {
		if !validProjectIdsSet.Has(id) {
			delete(projectCacheMap, id)
			removedIds.Add(id)
		}
	}

	return removedIds
}

var (
	lspApiObjectId      = "lspApiObjectId"
	lspApiObjectIdRef   = "lspApiObjectIdRef"
	lspApiObjectType    = "lspApiObjectType"
	lspApiProjectId     = "lspApiProjectId"
	lspApiTypeCheckerId = "lspApiTypeCheckerId"
)

func GetTypeOfElement(
	ctx context.Context,
	project *project.Project,
	fileName string,
	Range *lsproto.Range,
	forceReturnType bool,
	typeRequestKind lsproto.TypeRequestKind,
) (*collections.OrderedMap[string, interface{}], error) {
	projectIdNum := getProjectId(project)
	projectVersion := project.ProgramLastUpdate // TODO: possible race condition here
	program := project.GetProgram()
	sourceFile := program.GetSourceFile(fileName)
	if sourceFile == nil {
		return nil, nil
	}

	startOffset := scanner.GetECMAPositionOfLineAndCharacter(sourceFile, int(Range.Start.Line), int(Range.Start.Character))
	endOffset := scanner.GetECMAPositionOfLineAndCharacter(sourceFile, int(Range.End.Line), int(Range.End.Character))

	node := astnav.GetTokenAtPosition(sourceFile, startOffset).AsNode()
	for node != nil && node.End() < endOffset {
		node = node.Parent
	}
	if node == nil || node == sourceFile.AsNode() {
		return nil, nil
	}

	contextFlags := typeRequestKindToContextFlags(typeRequestKind)
	isContextual := contextFlags >= 0
	if isContextual && !ast.IsExpression(node) || (!isContextual && (ast.IsStringLiteral(node) || ast.IsNumericLiteral(node))) || ast.IsIdentifier(node) && ast.IsTypeReferenceNode(node.Parent) {
		if node.Pos() == node.Parent.Pos() && node.End() == node.Parent.End() {
			node = node.Parent
		}
	}

	typeChecker, done := program.GetTypeCheckerForFile(ctx, sourceFile)
	defer done()

	var t *checker.Type
	if contextFlags >= 0 {
		t = typeChecker.GetContextualType(node, checker.ContextFlags(contextFlags))
	} else {
		t = typeChecker.GetTypeAtLocation(node)
	}

	if t == nil && isContextual && node.Parent.Kind == ast.KindBinaryExpression {
		parentBinaryExpr := node.Parent.AsBinaryExpression()
		// from getContextualType in services/completions.ts
		right := parentBinaryExpr.Right
		if ls.IsEqualityOperatorKind(parentBinaryExpr.OperatorToken.Kind) {
			if node == right {
				t = typeChecker.GetTypeAtLocation(parentBinaryExpr.Left)
			} else {
				t = typeChecker.GetTypeAtLocation(right)
			}
		}
	}
	if t == nil {
		return nil, nil
	}

	convertContext := NewConvertContext(typeChecker, projectIdNum, projectVersion)

	var typeMap *collections.OrderedMap[string, interface{}]
	if forceReturnType || !convertContext.requestedTypeIds.Has(t.Id()) {
		typeMap = ConvertType(t, convertContext)
		convertContext.requestedTypeIds.Add(t.Id())
	} else {
		typeMap = collections.NewOrderedMapWithSizeHint[string, interface{}](4)
		typeMap.Set("id", t.Id())
		typeMap.Set(lspApiObjectType, "TypeObject")
	}
	typeMap.Set(lspApiProjectId, projectIdNum)
	typeMap.Set(lspApiTypeCheckerId, projectVersion)
	return typeMap, nil
}

func getConvertContext(
	ctx context.Context,
	projectId int,
	projectVersion uint64,
) (*ConvertContext, func(), error) {
	project, ok := getProject(projectId)
	if !ok {
		return nil, func() {}, OutdatedProjectIdError
	}

	if projectVersion != project.ProgramLastUpdate {
		return nil, func() {}, OutdatedProjectVersionError
	}
	checker, done := project.GetProgram().GetTypeChecker(ctx)
	return NewConvertContext(checker, projectId, projectVersion), func() {
		done()
	}, nil
}

func GetSymbolType(
	ctx context.Context,
	projectId int,
	projectVersion uint64,
	symbolId int,
) (*collections.OrderedMap[string, interface{}], error) {
	convertContext, done, err := getConvertContext(ctx, projectId, projectVersion)
	defer done()
	if err != nil {
		return nil, err
	}

	symbol, exists := convertContext.seenSymbolIds[ast.SymbolId(symbolId)]
	if !exists {
		return nil, errors.New("symbol not found")
	}

	t := convertContext.checker.GetTypeOfSymbol(symbol)
	result := ConvertType(t, convertContext)
	result.Set(lspApiProjectId, projectId)
	result.Set(lspApiTypeCheckerId, projectVersion)
	return result, nil
}

func GetTypeProperties(
	ctx context.Context,
	projectId int,
	projectVersion uint64,
	typeId int,
) (*collections.OrderedMap[string, interface{}], error) {
	convertContext, done, err := getConvertContext(ctx, projectId, projectVersion)
	defer done()
	if err != nil {
		return nil, err
	}

	t, exists := convertContext.seenTypeIds[checker.TypeId(typeId)]
	if !exists {
		return nil, errors.New("type not found")
	}

	result := ConvertTypeProperties(t, convertContext)
	result.Set(lspApiProjectId, projectId)
	result.Set(lspApiTypeCheckerId, projectVersion)
	return result, nil
}

func GetTypeProperty(
	ctx context.Context,
	projectId int,
	projectVersion uint64,
	typeId int,
	propertyName string,
) (*collections.OrderedMap[string, interface{}], error) {
	convertContext, done, err := getConvertContext(ctx, projectId, projectVersion)
	defer done()
	if err != nil {
		return nil, err
	}

	t, exists := convertContext.seenTypeIds[checker.TypeId(typeId)]
	if !exists {
		return nil, errors.New("type not found")
	}

	symbol := convertContext.checker.GetPropertyOfType(t, propertyName)
	if symbol == nil {
		return nil, nil
	}
	result := ConvertSymbol(symbol, convertContext)
	return result, nil
}

func AreTypesMutuallyAssignable(
	ctx context.Context,
	projectId int,
	projectVersion uint64,
	type1Id int,
	type2Id int,
) (*collections.OrderedMap[string, interface{}], error) {
	convertCtx, done, err := getConvertContext(ctx, projectId, projectVersion)
	defer done()
	if err != nil {
		return nil, err
	}

	type1, exists := convertCtx.seenTypeIds[checker.TypeId(type1Id)]
	if !exists {
		return nil, errors.New("type1 not found")
	}

	type2, exists := convertCtx.seenTypeIds[checker.TypeId(type2Id)]
	if !exists {
		return nil, errors.New("type2 not found")
	}

	isType1To2 := convertCtx.checker.IsTypeAssignableTo(type1, type2)
	isType2To1 := convertCtx.checker.IsTypeAssignableTo(type2, type1)

	areMutuallyAssignable := isType1To2 && isType2To1

	result := collections.NewOrderedMapWithSizeHint[string, interface{}](2)
	result.Set("areMutuallyAssignable", areMutuallyAssignable)
	return result, nil
}

func GetResolvedSignature(
	ctx context.Context,
	project *project.Project,
	fileName string,
	Range lsproto.Range,
) (*collections.OrderedMap[string, interface{}], error) {
	projectId := getProjectId(project)
	projectVersion := project.ProgramLastUpdate
	program := project.Program
	sourceFile := program.GetSourceFile(fileName)
	if sourceFile == nil {
		return nil, nil
	}

	startOffset := scanner.GetECMAPositionOfLineAndCharacter(sourceFile, int(Range.Start.Line), int(Range.Start.Character))
	endOffset := scanner.GetECMAPositionOfLineAndCharacter(sourceFile, int(Range.End.Line), int(Range.End.Character))

	typeChecker, done := program.GetTypeCheckerForFile(ctx, sourceFile)
	defer done()

	// Find the node at the given position
	node := astnav.GetTokenAtPosition(sourceFile, startOffset).AsNode()
	for node != nil && node.End() < endOffset {
		node = node.Parent
	}

	if node == nil || node == sourceFile.AsNode() {
		return nil, nil
	}

	// Find the call expression
	for !ast.IsCallLikeExpression(node) {
		node = node.Parent
		if node == nil || node == sourceFile.AsNode() {
			return nil, nil
		}
	}

	// Get the resolved signature
	signature := typeChecker.GetResolvedSignature(node)
	if signature == nil {
		return nil, nil
	}

	// Return the signature information
	convertCtx := NewConvertContext(typeChecker, projectId, projectVersion)
	prepared := ConvertSignature(signature, convertCtx)
	prepared.Set(lspApiTypeCheckerId, projectVersion)
	prepared.Set(lspApiProjectId, projectId)
	prepared.Set(lspApiObjectType, "SignatureObject")

	return prepared, nil
}

type ConvertContext struct {
	nextId                  int
	createdObjectsLspApiIds map[interface{}]int
	checker                 *checker.Checker
	requestedTypeIds        *collections.Set[checker.TypeId]
	seenTypeIds             map[checker.TypeId]*checker.Type
	seenSymbolIds           map[ast.SymbolId]*ast.Symbol
}

func NewConvertContext(checker *checker.Checker, projectId int, projectVersion uint64) *ConvertContext {
	cache := getProjectCache(projectId, projectVersion)
	return &ConvertContext{
		nextId:                  0,
		createdObjectsLspApiIds: make(map[interface{}]int),
		checker:                 checker,
		requestedTypeIds:        cache.requestedTypeIds,
		seenTypeIds:             cache.seenTypeIds,
		seenSymbolIds:           cache.seenSymbolIds,
	}
}

func (ctx *ConvertContext) GetLspApiObjectId(obj interface{}) (int, bool) {
	id, exists := ctx.createdObjectsLspApiIds[obj]
	return id, exists
}

func (ctx *ConvertContext) RegisterLspApiObject(obj interface{}) int {
	id := ctx.nextId
	ctx.nextId++
	ctx.createdObjectsLspApiIds[obj] = id
	return id
}

func FindReferenceOrConvert(sourceObj interface{},
	convertTarget func(lspApiObjectId int) *collections.OrderedMap[string, interface{}], ctx *ConvertContext,
) *collections.OrderedMap[string, interface{}] {
	lspApiObjectId, exists := ctx.GetLspApiObjectId(sourceObj)
	if exists {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](1)
		result.Set(lspApiObjectIdRef, lspApiObjectId)
		return result
	}

	lspApiObjectId = ctx.RegisterLspApiObject(sourceObj)
	newObject := convertTarget(lspApiObjectId)
	return newObject
}

func ConvertType(t *checker.Type, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(t, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		tscType := collections.NewOrderedMapWithSizeHint[string, interface{}](15)
		tscType.Set(lspApiObjectId, _lspApiObjectId)
		tscType.Set(lspApiObjectType, "TypeObject")
		tscType.Set("flags", strconv.Itoa(int(t.Flags()))) // Flags are u32, LSP number is s32

		// Handle aliasTypeArguments
		aliasType := t.Alias()
		if aliasType != nil && aliasType.TypeArguments() != nil {
			aliasArgs := make([]interface{}, 0)
			for _, t := range aliasType.TypeArguments() {
				aliasArgs = append(aliasArgs, ConvertType(t, ctx))
			}
			tscType.Set("aliasTypeArguments", aliasArgs)
		}

		if t.Flags()&checker.TypeFlagsObject != 0 {
			if target := t.Target(); target != nil {
				tscType.Set("target", ConvertType(target, ctx))
			}

			if (t.ObjectFlags() & checker.ObjectFlagsReference) != 0 {
				resolvedArgs := make([]interface{}, 0)
				typeArgs := ctx.checker.GetTypeArguments(t)
				for _, t := range typeArgs {
					// Filter out 'this' type
					if t.Flags()&checker.TypeFlagsTypeParameter == 0 || !t.AsTypeParameter().IsThisType() {
						resolvedArgs = append(resolvedArgs, ConvertType(t, ctx))
					}
				}
				tscType.Set("resolvedTypeArguments", resolvedArgs)
			}
		}

		// For UnionOrIntersection and TemplateLiteral types
		if t.Flags()&(checker.TypeFlagsUnionOrIntersection|checker.TypeFlagsTemplateLiteral) != 0 {
			typesArr := t.Types()
			types := make([]interface{}, 0)
			for _, t := range typesArr {
				types = append(types, ConvertType(t, ctx))
			}
			tscType.Set("types", types)
		}

		// For Literal types with freshType
		if t.Flags()&checker.TypeFlagsLiteral != 0 {
			literalType := t.AsLiteralType()
			if freshType := literalType.FreshType(); freshType != nil {
				tscType.Set("freshType", ConvertType(freshType, ctx))
			}
		}

		// For TypeParameter types
		if t.Flags()&checker.TypeFlagsTypeParameter != 0 {
			if constraint := ctx.checker.GetBaseConstraintOfType(t); constraint != nil {
				tscType.Set("constraint", ConvertType(constraint, ctx))
			}
		}

		// For Index types
		if t.Flags()&checker.TypeFlagsIndex != 0 {
			indexType := t.AsIndexType()
			tscType.Set("type", ConvertType(indexType.Target(), ctx))
		}

		// For IndexedAccess types
		if t.Flags()&checker.TypeFlagsIndexedAccess != 0 {
			indexedAccessType := t.AsIndexedAccessType()
			tscType.Set("objectType", ConvertType(indexedAccessType.ObjectType(), ctx))
			tscType.Set("indexType", ConvertType(indexedAccessType.IndexType(), ctx))
		}

		// For Conditional types
		if t.Flags()&checker.TypeFlagsConditional != 0 {
			conditionalType := t.AsConditionalType()
			tscType.Set("checkType", ConvertType(conditionalType.CheckType(), ctx))
			tscType.Set("extendsType", ConvertType(conditionalType.ExtendsType(), ctx))
		}

		// For Substitution types
		if t.Flags()&checker.TypeFlagsSubstitution != 0 {
			substitutionType := t.AsSubstitutionType()
			tscType.Set("baseType", ConvertType(substitutionType.BaseType(), ctx))
		}

		// Add symbol and aliasSymbol
		if t.Symbol() != nil {
			tscType.Set("symbol", ConvertSymbol(t.Symbol(), ctx))
		}
		if t.Alias() != nil && t.Alias().Symbol() != nil {
			tscType.Set("aliasSymbol", ConvertSymbol(t.Alias().Symbol(), ctx))
		}

		// Handle object flags
		if t.Flags()&checker.TypeFlagsObject != 0 {
			tscType.Set("objectFlags", strconv.Itoa(int(t.ObjectFlags())))
		}

		// Handle literal type values
		if t.Flags()&checker.TypeFlagsLiteral != 0 {
			literalType := t.AsLiteralType()
			if t.Flags()&checker.TypeFlagsBigIntLiteral != 0 {
				// Convert BigInt literal value
				tscType.Set("value", ConvertPseudoBigInt(literalType.Value().(jsnum.PseudoBigInt), ctx))
			} else {
				// For other literal types
				tscType.Set("value", literalType.Value())
			}
		}

		// Handle enum literal
		if t.Flags()&checker.TypeFlagsEnumLiteral != 0 {
			tscType.Set("nameType", getEnumQualifiedName(t))
		}

		// Handle template literal
		if t.Flags()&checker.TypeFlagsTemplateLiteral != 0 {
			templateLiteralType := t.AsTemplateLiteralType()
			tscType.Set("texts", templateLiteralType.Texts())
		}

		// Handle type parameter isThisType
		if t.Flags()&checker.TypeFlagsTypeParameter != 0 {
			typeParam := t.AsTypeParameter()
			if typeParam.IsThisType() {
				tscType.Set("isThisType", true)
			}
		}

		// Handle intrinsic name
		if t.Flags()&checker.TypeFlagsIntrinsic != 0 {
			intrinsicType := t.AsIntrinsicType()
			if intrinsicType != nil && intrinsicType.IntrinsicName() != "" {
				tscType.Set("intrinsicName", intrinsicType.IntrinsicName())
			}
		}

		// Handle tuple element flags
		if t.Flags()&checker.TypeFlagsObject != 0 && t.ObjectFlags()&checker.ObjectFlagsTuple != 0 {
			tupleType := t.AsTupleType()
			elementInfos := tupleType.ElementInfos()
			elementFlags := make([]interface{}, len(elementInfos))
			for i, elementInfo := range elementInfos {
				elementFlags[i] = strconv.Itoa(int(elementInfo.ElementFlags()))
			}
			tscType.Set("elementFlags", elementFlags)
		}

		// Add type ID
		typeId := t.Id()
		tscType.Set("id", typeId)
		ctx.seenTypeIds[typeId] = t
		return tscType
	}, ctx)
}

func getSourceFileParent(node *ast.Node) *ast.Node {
	if ast.IsSourceFile(node) {
		return nil
	}
	current := node.Parent
	for current != nil {
		if ast.IsSourceFile(current) {
			return current
		}
		current = current.Parent
	}
	return nil
}

func ConvertNode(node *ast.Node, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(node, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](1)
		result.Set(lspApiObjectId, _lspApiObjectId)

		if ast.IsSourceFile(node) {
			result.Set(lspApiObjectType, "SourceFileObject")
			result.Set("fileName", node.AsSourceFile().FileName())
			return result
		} else {
			result.Set(lspApiObjectType, "NodeObject")
			sourceFileParent := getSourceFileParent(node)
			if sourceFileParent == nil || node.Pos() == -1 || node.End() == -1 {
				if sourceFileParent != nil {
					result.Set("parent", ConvertNode(sourceFileParent, ctx))
				}
				return result
			}

			// Add range information
			if sourceFileParent != nil {
				sourceFile := sourceFileParent.AsSourceFile()
				startLine, startChar := scanner.GetECMALineAndCharacterOfPosition(sourceFile, node.Pos())
				endLine, endChar := scanner.GetECMALineAndCharacterOfPosition(sourceFile, node.End())
				result.Set("range", &lsproto.Range{
					Start: lsproto.Position{Line: uint32(startLine), Character: uint32(startChar)},
					End:   lsproto.Position{Line: uint32(endLine), Character: uint32(endChar)},
				})
			}

			// Add parent information
			if sourceFileParent != nil {
				result.Set("parent", ConvertNode(sourceFileParent, ctx))
			}

			// Check for computed property
			if node.Name() != nil && ast.IsComputedPropertyName(node.Name()) {
				result.Set("computedProperty", true)
			}
			return result
		}
	}, ctx)
}

func getEnumQualifiedName(t *checker.Type) string {
	if t == nil || t.Symbol() == nil {
		return ""
	}

	qName := ""
	current := t.Symbol().Parent
	for current != nil && !(current.ValueDeclaration != nil && ast.IsSourceFile(current.ValueDeclaration)) {
		if qName == "" {
			qName = current.Name
		} else {
			qName = current.Name + "." + qName
		}
		current = current.Parent
	}
	return qName
}

func ConvertPseudoBigInt(pseudoBigInt jsnum.PseudoBigInt, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(pseudoBigInt, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](4)
		result.Set(lspApiObjectId, _lspApiObjectId)
		result.Set("negative", pseudoBigInt.Negative)
		result.Set("base10Value", pseudoBigInt.Base10Value)
		return result
	}, ctx)
}

func ConvertIndexInfo(indexInfo *checker.IndexInfo, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(indexInfo, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](7)
		result.Set(lspApiObjectId, _lspApiObjectId)
		result.Set(lspApiObjectType, "IndexInfo")
		result.Set("keyType", ConvertType(indexInfo.KeyType(), ctx))
		result.Set("type", ConvertType(indexInfo.ValueType(), ctx))
		result.Set("isReadonly", indexInfo.IsReadonly())

		if indexInfo.Declaration() != nil {
			result.Set("declaration", ConvertNode(indexInfo.Declaration(), ctx))
		}
		return result
	}, ctx)
}

func ConvertSymbol(symbol *ast.Symbol, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(symbol, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](7)
		result.Set(lspApiObjectId, _lspApiObjectId)
		result.Set(lspApiObjectType, "SymbolObject")
		result.Set("flags", strconv.Itoa(int(symbol.Flags)))
		escapedName := symbol.Name
		if strings.Contains(escapedName, ast.InternalSymbolNamePrefix) {
			escapedName = strings.ReplaceAll(escapedName, ast.InternalSymbolNamePrefix, "__")
		}
		result.Set("escapedName", escapedName)

		if symbol.Declarations != nil && len(symbol.Declarations) > 0 {
			declarations := make([]interface{}, 0)
			for _, d := range symbol.Declarations {
				declarations = append(declarations, ConvertNode(d, ctx))
			}
			result.Set("declarations", declarations)
		}

		if symbol.ValueDeclaration != nil {
			result.Set("valueDeclaration", ConvertNode(symbol.ValueDeclaration, ctx))
		}

		// Get symbol ID
		symbolId := ast.GetSymbolId(symbol)
		ctx.seenSymbolIds[symbolId] = symbol
		result.Set("id", symbolId)

		return result
	}, ctx)
}

func ConvertTypeProperties(t *checker.Type, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(t, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		prepared := collections.NewOrderedMapWithSizeHint[string, interface{}](10)
		prepared.Set(lspApiObjectId, _lspApiObjectId)
		prepared.Set(lspApiObjectType, "TypeObject")
		prepared.Set("flags", strconv.Itoa(int(t.Flags())))
		prepared.Set("objectFlags", strconv.Itoa(int(t.ObjectFlags())))

		if t.Flags()&checker.TypeFlagsObject != 0 {
			assignObjectTypeProperties(t.AsObjectType(), ctx, prepared)
		}
		if t.Flags()&checker.TypeFlagsUnionOrIntersection != 0 {
			assignUnionOrIntersectionTypeProperties(t.AsUnionOrIntersectionType(), ctx, prepared)
		}
		if t.Flags()&checker.TypeFlagsConditional != 0 {
			assignConditionalTypeProperties(t.AsConditionalType(), ctx, prepared)
		}
		return prepared
	}, ctx)
}

func assignObjectTypeProperties(t *checker.ObjectType, ctx *ConvertContext, tscType *collections.OrderedMap[string, interface{}]) {
	constructSignatures := make([]interface{}, 0)
	for _, s := range ctx.checker.GetConstructSignatures(t.AsType()) {
		constructSignatures = append(constructSignatures, ConvertSignature(s, ctx))
	}
	tscType.Set("constructSignatures", constructSignatures)

	callSignatures := make([]interface{}, 0)
	for _, s := range ctx.checker.GetCallSignatures(t.AsType()) {
		callSignatures = append(callSignatures, ConvertSignature(s, ctx))
	}
	tscType.Set("callSignatures", callSignatures)

	properties := make([]interface{}, 0)
	for _, p := range ctx.checker.GetPropertiesOfType(t.AsType()) {
		properties = append(properties, ConvertSymbol(p, ctx))
	}
	tscType.Set("properties", properties)

	indexInfos := make([]interface{}, 0)
	for _, info := range ctx.checker.GetIndexInfosOfType(t.AsType()) {
		indexInfos = append(indexInfos, ConvertIndexInfo(info, ctx))
	}
	tscType.Set("indexInfos", indexInfos)
}

func assignUnionOrIntersectionTypeProperties(t *checker.UnionOrIntersectionType, ctx *ConvertContext, tscType *collections.OrderedMap[string, interface{}]) {
	resolvedProperties := make([]interface{}, 0)
	for _, p := range ctx.checker.GetPropertiesOfType(t.AsType()) {
		resolvedProperties = append(resolvedProperties, ConvertSymbol(p, ctx))
	}
	tscType.Set("resolvedProperties", resolvedProperties)

	callSignatures := make([]interface{}, 0)
	for _, s := range ctx.checker.GetCallSignatures(t.AsType()) {
		callSignatures = append(callSignatures, ConvertSignature(s, ctx))
	}
	tscType.Set("callSignatures", callSignatures)

	constructSignatures := make([]interface{}, 0)
	for _, s := range ctx.checker.GetConstructSignatures(t.AsType()) {
		constructSignatures = append(constructSignatures, ConvertSignature(s, ctx))
	}
	tscType.Set("constructSignatures", constructSignatures)
}

func assignConditionalTypeProperties(t *checker.ConditionalType, ctx *ConvertContext, tscType *collections.OrderedMap[string, interface{}]) {
	trueType := ctx.checker.GetTrueTypeFromConditionalType(t.AsType())
	if trueType != nil {
		tscType.Set("resolvedTrueType", ConvertType(trueType, ctx))
	}

	falseType := ctx.checker.GetFalseTypeFromConditionalType(t.AsType())
	if falseType != nil {
		tscType.Set("resolvedFalseType", ConvertType(falseType, ctx))
	}
}

func ConvertSignature(signature *checker.Signature, ctx *ConvertContext) *collections.OrderedMap[string, interface{}] {
	return FindReferenceOrConvert(signature, func(_lspApiObjectId int) *collections.OrderedMap[string, interface{}] {
		result := collections.NewOrderedMapWithSizeHint[string, interface{}](5)
		result.Set(lspApiObjectId, _lspApiObjectId)
		result.Set(lspApiObjectType, "SignatureObject")

		if declaration := signature.Declaration(); declaration != nil {
			result.Set("declaration", ConvertNode(declaration, ctx))
		}

		if returnType := ctx.checker.GetReturnTypeOfSignature(signature); returnType != nil {
			result.Set("resolvedReturnType", ConvertType(returnType, ctx))
		}

		parameters := make([]interface{}, 0)
		for _, param := range signature.Parameters() {
			parameters = append(parameters, ConvertSymbol(param, ctx))
		}
		result.Set("parameters", parameters)

		result.Set("flags", strconv.Itoa(convertSignatureFlags(signature)))

		return result
	}, ctx)
}

func convertSignatureFlags(signature *checker.Signature) int {
	var result = 0
	if signature.IsSignatureCandidateForOverloadFailure() {
		result = result | (1 << 7)
	}
	return result
}

func typeRequestKindToContextFlags(typeRequestKind lsproto.TypeRequestKind) int {
	switch typeRequestKind {
	case lsproto.TypeRequestKindDefault:
		return -1
	case lsproto.TypeRequestKindContextual:
		return 0
	case lsproto.TypeRequestKindContextualCompletions:
		return 4
	default:
		panic(fmt.Sprintf("Unexpected typeRequestKind %s", typeRequestKind))
	}
}
