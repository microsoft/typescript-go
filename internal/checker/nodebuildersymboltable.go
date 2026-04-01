package checker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// symbolTableSerializationState holds the mutable state for symbolTableToDeclarationStatements.
// In Strada this is captured as closure state.
type symbolTableSerializationState struct {
	b                     *NodeBuilderImpl
	results               []*ast.Node
	visitedSymbols        collections.Set[ast.SymbolId]
	deferredPrivatesStack []map[ast.SymbolId]*ast.Symbol
	addingDeclare         bool
	enclosingDeclaration  *ast.Node
	usedSymbolNames       collections.Set[string]
	remappedSymbolNames   map[ast.SymbolId]string
}

func (b *NodeBuilderImpl) symbolTableToDeclarationStatements(symbolTable *ast.SymbolTable) []*ast.Node {
	s := &symbolTableSerializationState{
		b:                    b,
		enclosingDeclaration: b.ctx.enclosingDeclaration,
		addingDeclare:        true,
		usedSymbolNames:      collections.Set[string]{},
		remappedSymbolNames:  make(map[ast.SymbolId]string),
	}

	// Cache symbol names
	for name, sym := range *symbolTable {
		s.getInternalSymbolName(sym, name)
	}

	s.visitSymbolTable(symbolTable, false, false)
	return s.results
}

func (s *symbolTableSerializationState) visitSymbolTable(symbolTable *ast.SymbolTable, suppressNewPrivateContext bool, propertyAsAlias bool) {
	if !suppressNewPrivateContext {
		s.deferredPrivatesStack = append(s.deferredPrivatesStack, make(map[ast.SymbolId]*ast.Symbol))
	}
	i := 0
	size := len(*symbolTable)
	symbols := make([]*ast.Symbol, 0, size)
	for _, symbol := range *symbolTable {
		symbols = append(symbols, symbol)
	}
	s.b.ch.sortSymbols(symbols)
	for _, symbol := range symbols {
		i++
		if s.b.checkTruncationLengthIfExpanding() && (i+2 < size-1) {
			s.b.ctx.out.Truncated = true
			s.results = append(s.results, s.createTruncationStatement(fmt.Sprintf("... (%d more ...)", size-i)))
			s.serializeSymbol(symbols[len(symbols)-1], false, propertyAsAlias)
			break
		}
		s.serializeSymbol(symbol, false, propertyAsAlias)
	}
	if !suppressNewPrivateContext {
		// deferredPrivates will be filled up by visiting the symbol table
		// And will continue to iterate as elements are added while visited `deferredPrivates`
		last := s.deferredPrivatesStack[len(s.deferredPrivatesStack)-1]
		deferredSymbols := make([]*ast.Symbol, 0, len(last))
		for _, symbol := range last {
			deferredSymbols = append(deferredSymbols, symbol)
		}
		s.b.ch.sortSymbols(deferredSymbols)
		for _, symbol := range deferredSymbols {
			s.serializeSymbol(symbol, true, propertyAsAlias)
		}
		s.deferredPrivatesStack = s.deferredPrivatesStack[:len(s.deferredPrivatesStack)-1]
	}
}

func (s *symbolTableSerializationState) serializeSymbol(symbol *ast.Symbol, isPrivate bool, propertyAsAlias bool) {
	s.b.ch.getPropertiesOfType(s.b.ch.getTypeOfSymbol(symbol)) // resolve properties to trigger merges
	visitedSym := s.b.ch.getMergedSymbol(symbol)
	if s.visitedSymbols.Has(ast.GetSymbolId(visitedSym)) {
		return
	}
	s.visitedSymbols.Add(ast.GetSymbolId(visitedSym))

	skipMembershipCheck := !isPrivate
	if skipMembershipCheck || (len(symbol.Declarations) > 0 && core.Some(symbol.Declarations, func(d *ast.Node) bool {
		return ast.FindAncestor(d, func(n *ast.Node) bool { return n == s.enclosingDeclaration }) != nil
	})) {
		scopeCleanup := cloneNodeBuilderContext(s.b.ctx)
		s.serializeSymbolWorker(symbol, isPrivate, propertyAsAlias, symbol.Name)
		scopeCleanup()
	}
}

func (s *symbolTableSerializationState) serializeSymbolWorker(symbol *ast.Symbol, isPrivate bool, propertyAsAlias bool, symbolName string) {
	escapedSymbolName := symbol.Name
	isDefault := escapedSymbolName == ast.InternalSymbolNameDefault
	if isPrivate && s.b.ctx.flags&nodebuilder.FlagsAllowAnonymousIdentifier == 0 && ast.IsNonContextualKeyword(scanner.StringToToken(symbolName)) && !isDefault {
		s.b.ctx.encounteredError = true
		return
	}

	needsPostExportDefault := isDefault && (symbol.Flags&ast.SymbolFlagsExportDoesNotSupportDefaultModifier != 0 ||
		(symbol.Flags&ast.SymbolFlagsFunction != 0 && len(s.b.ch.getPropertiesOfType(s.b.ch.getTypeOfSymbol(symbol))) > 0)) &&
		symbol.Flags&ast.SymbolFlagsAlias == 0
	needsExportDeclaration := !needsPostExportDefault && !isPrivate && ast.IsNonContextualKeyword(scanner.StringToToken(symbolName)) && !isDefault

	if needsPostExportDefault || needsExportDeclaration {
		isPrivate = true
	}

	modifierFlags := ast.ModifierFlagsNone
	if !isPrivate {
		modifierFlags |= ast.ModifierFlagsExport
	}
	if isDefault && !needsPostExportDefault {
		modifierFlags |= ast.ModifierFlagsDefault
	}

	isConstMergedWithNS := symbol.Flags&ast.SymbolFlagsModule != 0 &&
		symbol.Flags&(ast.SymbolFlagsBlockScopedVariable|ast.SymbolFlagsFunctionScopedVariable|ast.SymbolFlagsProperty) != 0 &&
		escapedSymbolName != ast.InternalSymbolNameExportEquals
	isConstMergedWithNSPrintableAsSignatureMerge := isConstMergedWithNS && s.isTypeRepresentableAsFunctionNamespaceMerge(s.b.ch.getTypeOfSymbol(symbol), symbol)

	if symbol.Flags&(ast.SymbolFlagsFunction|ast.SymbolFlagsMethod) != 0 || isConstMergedWithNSPrintableAsSignatureMerge {
		s.serializeAsFunctionNamespaceMerge(s.b.ch.getTypeOfSymbol(symbol), symbol, s.getInternalSymbolName(symbol, symbolName), modifierFlags)
	}
	if symbol.Flags&ast.SymbolFlagsTypeAlias != 0 {
		s.serializeTypeAlias(symbol, symbolName, modifierFlags)
	}
	if symbol.Flags&(ast.SymbolFlagsBlockScopedVariable|ast.SymbolFlagsFunctionScopedVariable|ast.SymbolFlagsProperty|ast.SymbolFlagsAccessor) != 0 &&
		escapedSymbolName != ast.InternalSymbolNameExportEquals &&
		symbol.Flags&ast.SymbolFlagsPrototype == 0 &&
		symbol.Flags&ast.SymbolFlagsClass == 0 &&
		symbol.Flags&ast.SymbolFlagsMethod == 0 &&
		!isConstMergedWithNSPrintableAsSignatureMerge {
		s.serializeVariableOrProperty(symbol, symbolName, isPrivate, needsPostExportDefault, modifierFlags, propertyAsAlias)
	}
	if symbol.Flags&ast.SymbolFlagsEnum != 0 {
		s.serializeEnum(symbol, symbolName, modifierFlags)
	}
	if symbol.Flags&ast.SymbolFlagsClass != 0 {
		if symbol.Flags&ast.SymbolFlagsProperty != 0 &&
			symbol.ValueDeclaration != nil &&
			ast.IsBinaryExpression(symbol.ValueDeclaration.Parent) &&
			ast.IsClassExpression(symbol.ValueDeclaration.Parent.AsBinaryExpression().Right) {
			s.serializeAsAlias(symbol, s.getInternalSymbolName(symbol, symbolName), modifierFlags)
		} else {
			s.serializeAsClass(symbol, s.getInternalSymbolName(symbol, symbolName), modifierFlags)
		}
	}
	if (symbol.Flags&(ast.SymbolFlagsValueModule|ast.SymbolFlagsNamespaceModule) != 0 && (!isConstMergedWithNS || s.isTypeOnlyNamespace(symbol))) || isConstMergedWithNSPrintableAsSignatureMerge {
		s.serializeModule(symbol, symbolName, modifierFlags)
	}
	if symbol.Flags&ast.SymbolFlagsInterface != 0 && symbol.Flags&ast.SymbolFlagsClass == 0 {
		s.serializeInterface(symbol, symbolName, modifierFlags)
	}
	if symbol.Flags&ast.SymbolFlagsAlias != 0 {
		s.serializeAsAlias(symbol, s.getInternalSymbolName(symbol, symbolName), modifierFlags)
	}
	if symbol.Flags&ast.SymbolFlagsProperty != 0 && symbol.Name == ast.InternalSymbolNameExportEquals {
		s.serializeMaybeAliasAssignment(symbol)
	}

	if needsPostExportDefault {
		internalSymbolName := s.getInternalSymbolName(symbol, symbolName)
		s.b.ctx.approximateLength += 16 + len(internalSymbolName)
		s.results = append(s.results, s.b.f.NewExportAssignment(nil, false, nil, s.b.f.NewIdentifier(internalSymbolName)))
	} else if needsExportDeclaration {
		internalSymbolName := s.getInternalSymbolName(symbol, symbolName)
		s.b.ctx.approximateLength += 22 + len(symbolName) + len(internalSymbolName)
		s.results = append(s.results, s.b.f.NewExportDeclaration(
			nil,
			false,
			s.b.f.NewNamedExports(s.b.f.NewNodeList([]*ast.Node{
				s.b.f.NewExportSpecifier(false, s.b.f.NewIdentifier(internalSymbolName), s.b.f.NewIdentifier(symbolName)),
			})),
			nil,
			nil,
		))
	}
}

func (s *symbolTableSerializationState) includePrivateSymbol(symbol *ast.Symbol) {
	if core.Some(symbol.Declarations, isPartOfParameterDeclaration) {
		return
	}
	s.getUnusedName(symbol.Name, symbol)
	isExternalImportAlias := symbol.Flags&ast.SymbolFlagsAlias != 0 && !core.Some(symbol.Declarations, func(d *ast.Node) bool {
		return ast.FindAncestor(d, ast.IsExportDeclaration) != nil ||
			ast.IsNamespaceExport(d) ||
			(ast.IsImportEqualsDeclaration(d) && !ast.IsExternalModuleReference(d.AsImportEqualsDeclaration().ModuleReference))
	})
	idx := len(s.deferredPrivatesStack) - 1
	if isExternalImportAlias {
		idx = 0
	}
	s.deferredPrivatesStack[idx][ast.GetSymbolId(symbol)] = symbol
}

func isPartOfParameterDeclaration(d *ast.Node) bool {
	return ast.FindAncestor(d, ast.IsParameter) != nil
}

func (s *symbolTableSerializationState) isExportingScope(enclosingDeclaration *ast.Node) bool {
	return (ast.IsSourceFile(enclosingDeclaration) && (ast.IsExternalOrCommonJSModule(enclosingDeclaration.AsSourceFile()) || ast.IsJsonSourceFile(enclosingDeclaration.AsSourceFile()))) ||
		(ast.IsAmbientModule(enclosingDeclaration) && !ast.IsGlobalScopeAugmentation(enclosingDeclaration))
}

func (s *symbolTableSerializationState) addResult(node *ast.Node, additionalModifierFlags ast.ModifierFlags) {
	if ast.CanHaveModifiers(node) {
		newModifierFlags := ast.ModifierFlagsNone
		enclosingDeclaration := s.b.ctx.enclosingDeclaration
		if enclosingDeclaration != nil && (enclosingDeclaration.Kind == ast.KindJSDocTypedefTag || enclosingDeclaration.Kind == ast.KindJSDocCallbackTag) {
			enclosingDeclaration = ast.GetSourceFileOfNode(enclosingDeclaration).AsNode()
		}

		canExport := ast.IsEnumDeclaration(node) || ast.IsVariableStatement(node) || ast.IsFunctionDeclaration(node) || ast.IsClassDeclaration(node) ||
			ast.IsInterfaceDeclaration(node) || ast.IsTypeAliasDeclaration(node) || (ast.IsModuleDeclaration(node) && node.Parent != nil && !ast.IsExternalModuleAugmentation(node) && !ast.IsGlobalScopeAugmentation(node))
		if additionalModifierFlags&ast.ModifierFlagsExport != 0 &&
			enclosingDeclaration != nil &&
			(s.isExportingScope(enclosingDeclaration) || ast.IsModuleDeclaration(enclosingDeclaration)) &&
			canExport {
			newModifierFlags |= ast.ModifierFlagsExport
		}
		if s.addingDeclare &&
			newModifierFlags&ast.ModifierFlagsExport == 0 &&
			(enclosingDeclaration == nil || enclosingDeclaration.Flags&ast.NodeFlagsAmbient == 0) &&
			(ast.IsEnumDeclaration(node) || ast.IsVariableStatement(node) || ast.IsFunctionDeclaration(node) || ast.IsClassDeclaration(node) || ast.IsModuleDeclaration(node)) {
			newModifierFlags |= ast.ModifierFlagsAmbient
		}
		if additionalModifierFlags&ast.ModifierFlagsDefault != 0 && (ast.IsClassDeclaration(node) || ast.IsInterfaceDeclaration(node) || ast.IsFunctionDeclaration(node)) {
			newModifierFlags |= ast.ModifierFlagsDefault
		}
		if newModifierFlags != ast.ModifierFlagsNone {
			oldModifierFlags := node.ModifierFlags()
			modifiers := ast.CreateModifiersFromModifierFlags(newModifierFlags|oldModifierFlags, s.b.f.NewModifier)
			node = ast.ReplaceModifiers(s.b.f, node, s.b.f.NewModifierList(modifiers))
		}
	}
	s.results = append(s.results, node)
}

func (s *symbolTableSerializationState) serializeTypeAlias(symbol *ast.Symbol, symbolName string, modifierFlags ast.ModifierFlags) {
	aliasType := s.b.ch.getDeclaredTypeOfTypeAlias(symbol)
	typeParams := s.b.ch.getLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
	typeParamDecls := core.Map(typeParams, func(p *Type) *ast.Node { return s.b.typeParameterToDeclaration(p) })
	restoreFlags := s.b.saveRestoreFlags()
	s.b.ctx.flags |= nodebuilder.FlagsInTypeAlias
	typeNode := s.b.typeToTypeNode(aliasType)
	internalSymbolName := s.getInternalSymbolName(symbol, symbolName)
	s.b.ctx.approximateLength += 8 + len(internalSymbolName)
	s.addResult(
		s.b.f.NewTypeAliasDeclaration(nil, s.b.f.NewIdentifier(internalSymbolName), s.b.f.NewNodeList(typeParamDecls), typeNode),
		modifierFlags,
	)
	restoreFlags()
}

func (s *symbolTableSerializationState) serializeInterface(symbol *ast.Symbol, symbolName string, modifierFlags ast.ModifierFlags) {
	internalSymbolName := s.getInternalSymbolName(symbol, symbolName)
	s.b.ctx.approximateLength += 14 + len(internalSymbolName)
	interfaceType := s.b.ch.getDeclaredTypeOfClassOrInterface(symbol)
	localParams := s.b.ch.getLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
	typeParamDecls := core.Map(localParams, func(p *Type) *ast.Node { return s.b.typeParameterToDeclaration(p) })
	baseTypes := s.b.ch.getBaseTypes(interfaceType)
	var baseType *Type
	if len(baseTypes) > 0 {
		baseType = s.b.ch.getIntersectionType(baseTypes)
	}
	members := s.serializePropertySymbolsForInterface(s.b.ch.getPropertiesOfType(interfaceType), baseType)
	callSignatures := s.serializeSignatures(SignatureKindCall, interfaceType, baseType, ast.KindCallSignature)
	constructSignatures := s.serializeSignatures(SignatureKindConstruct, interfaceType, baseType, ast.KindConstructSignature)
	indexSignatures := s.serializeIndexSignatures(interfaceType, baseType)

	var heritageClauses []*ast.Node
	if len(baseTypes) > 0 {
		var hcTypes []*ast.Node
		for _, bt := range baseTypes {
			if ref := s.trySerializeAsTypeReference(bt, ast.SymbolFlagsValue); ref != nil {
				hcTypes = append(hcTypes, ref)
			}
		}
		if len(hcTypes) > 0 {
			heritageClauses = []*ast.Node{s.b.f.NewHeritageClause(ast.KindExtendsKeyword, s.b.f.NewNodeList(hcTypes))}
		}
	}

	allMembers := make([]*ast.Node, 0, len(indexSignatures)+len(constructSignatures)+len(callSignatures)+len(members))
	allMembers = append(allMembers, indexSignatures...)
	allMembers = append(allMembers, constructSignatures...)
	allMembers = append(allMembers, callSignatures...)
	allMembers = append(allMembers, members...)

	s.addResult(
		s.b.f.NewInterfaceDeclaration(
			nil,
			s.b.f.NewIdentifier(internalSymbolName),
			s.b.f.NewNodeList(typeParamDecls),
			s.b.f.NewNodeList(heritageClauses),
			s.b.f.NewNodeList(allMembers),
		),
		modifierFlags,
	)
}

func (s *symbolTableSerializationState) serializeAsClass(symbol *ast.Symbol, localName string, modifierFlags ast.ModifierFlags) {
	s.b.ctx.approximateLength += 9 + len(localName)
	originalDecl := core.Find(symbol.Declarations, ast.IsClassLike)
	oldEnclosing := s.b.ctx.enclosingDeclaration
	if originalDecl != nil {
		s.b.ctx.enclosingDeclaration = originalDecl
	}
	localParams := s.b.ch.getLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
	typeParamDecls := core.Map(localParams, func(p *Type) *ast.Node { return s.b.typeParameterToDeclaration(p) })
	classType := s.b.ch.getTypeWithThisArgument(s.b.ch.getDeclaredTypeOfClassOrInterface(symbol), nil, false)
	baseTypes := s.b.ch.getBaseTypes(s.b.ch.getTargetType(classType))
	staticType := s.b.ch.getTypeOfSymbol(symbol)
	isClass := staticType.symbol != nil && staticType.symbol.ValueDeclaration != nil && ast.IsClassLike(staticType.symbol.ValueDeclaration)
	var staticBaseType *Type
	if isClass {
		staticBaseType = s.b.ch.getBaseConstructorTypeOfClass(s.b.ch.getDeclaredTypeOfClassOrInterface(symbol))
	} else {
		staticBaseType = s.b.ch.anyType
	}

	// Heritage clauses
	var heritageClauses []*ast.Node
	if len(baseTypes) > 0 {
		extendsTypes := core.Map(baseTypes, func(bt *Type) *ast.Node { return s.serializeBaseType(bt, staticBaseType, localName) })
		heritageClauses = append(heritageClauses, s.b.f.NewHeritageClause(ast.KindExtendsKeyword, s.b.f.NewNodeList(extendsTypes)))
	}
	implementsTypes := s.getImplementsTypes(classType)
	if len(implementsTypes) > 0 {
		var implExprs []*ast.Node
		for _, t := range implementsTypes {
			if impl := s.serializeImplementedType(t); impl != nil {
				implExprs = append(implExprs, impl)
			}
		}
		if len(implExprs) > 0 {
			heritageClauses = append(heritageClauses, s.b.f.NewHeritageClause(ast.KindImplementsKeyword, s.b.f.NewNodeList(implExprs)))
		}
	}

	symbolProps := s.getNonInheritedProperties(classType, baseTypes, s.b.ch.getPropertiesOfType(classType))
	publicSymbolProps := core.Filter(symbolProps, func(sym *ast.Symbol) bool { return !isHashPrivate(sym) })
	hasPrivateIdentifier := core.Some(symbolProps, isHashPrivate)

	var privateProperties []*ast.Node
	if hasPrivateIdentifier {
		if isExpanding(s.b.ctx) {
			privateProperties = s.serializePropertySymbolsForClass(core.Filter(symbolProps, isHashPrivate), false, core.FirstOrNil(baseTypes))
		} else {
			privateProperties = []*ast.Node{s.b.f.NewPropertyDeclaration(
				nil,
				s.b.f.NewPrivateIdentifier("#private"),
				nil, nil, nil,
			)}
		}
	}

	publicProperties := s.serializePropertySymbolsForClass(publicSymbolProps, false, core.FirstOrNil(baseTypes))
	staticMembers := s.serializePropertySymbolsForClass(
		core.Filter(s.b.ch.getPropertiesOfType(staticType), func(p *ast.Symbol) bool {
			return p.Flags&ast.SymbolFlagsPrototype == 0 && p.Name != "prototype" && !s.isNamespaceMember(p)
		}),
		true,
		staticBaseType,
	)

	isNonConstructableClassLikeInJsFile := !isClass &&
		symbol.ValueDeclaration != nil &&
		ast.IsInJSFile(symbol.ValueDeclaration) &&
		!core.Some(s.b.ch.getSignaturesOfType(staticType, SignatureKindConstruct), func(_ *Signature) bool { return true })

	var constructors []*ast.Node
	if isNonConstructableClassLikeInJsFile {
		s.b.ctx.approximateLength += 21
		modifiers := ast.CreateModifiersFromModifierFlags(ast.ModifierFlagsPrivate, s.b.f.NewModifier)
		constructors = []*ast.Node{s.b.f.NewConstructorDeclaration(s.b.f.NewModifierList(modifiers), nil, s.b.f.NewNodeList(nil), nil, nil, nil)}
	} else {
		constructors = s.serializeSignatures(SignatureKindConstruct, staticType, staticBaseType, ast.KindConstructor)
	}

	indexSignatures := s.serializeIndexSignatures(classType, core.FirstOrNil(baseTypes))

	s.b.ctx.enclosingDeclaration = oldEnclosing

	allMembers := make([]*ast.Node, 0, len(indexSignatures)+len(staticMembers)+len(constructors)+len(publicProperties)+len(privateProperties))
	allMembers = append(allMembers, indexSignatures...)
	allMembers = append(allMembers, staticMembers...)
	allMembers = append(allMembers, constructors...)
	allMembers = append(allMembers, publicProperties...)
	allMembers = append(allMembers, privateProperties...)

	s.addResult(
		s.b.f.NewClassDeclaration(
			nil,
			s.b.f.NewIdentifier(localName),
			s.b.f.NewNodeList(typeParamDecls),
			s.b.f.NewNodeList(heritageClauses),
			s.b.f.NewNodeList(allMembers),
		),
		modifierFlags,
	)
}

func (s *symbolTableSerializationState) serializeEnum(symbol *ast.Symbol, symbolName string, modifierFlags ast.ModifierFlags) {
	internalSymbolName := s.getInternalSymbolName(symbol, symbolName)
	s.b.ctx.approximateLength += 9 + len(internalSymbolName)
	var members []*ast.Node
	memberProps := core.Filter(s.b.ch.getPropertiesOfType(s.b.ch.getTypeOfSymbol(symbol)), func(p *ast.Symbol) bool {
		return p.Flags&ast.SymbolFlagsEnumMember != 0
	})
	for i, p := range memberProps {
		if s.b.checkTruncationLengthIfExpanding() && (i+2 < len(memberProps)-1) {
			s.b.ctx.out.Truncated = true
			members = append(members, s.b.f.NewEnumMember(s.b.f.NewStringLiteral(fmt.Sprintf(" ... %d more ... ", len(memberProps)-i), 0), nil))
			last := memberProps[len(memberProps)-1]
			initializedValue := s.getEnumMemberInitializer(last)
			memberName := last.Name
			members = append(members, s.b.f.NewEnumMember(s.b.f.NewIdentifier(memberName), initializedValue))
			break
		}
		memberDecl := core.Find(p.Declarations, ast.IsEnumMember)
		var initializer *ast.Node
		if isExpanding(s.b.ctx) && memberDecl != nil && memberDecl.AsEnumMember().Initializer != nil {
			initializer = s.b.f.DeepCloneNode(memberDecl.AsEnumMember().Initializer)
		} else {
			initializer = s.getEnumMemberInitializer(p)
		}
		memberName := p.Name
		s.b.ctx.approximateLength += 4 + len(memberName)
		members = append(members, s.b.f.NewEnumMember(s.b.f.NewIdentifier(memberName), initializer))
	}

	constModifier := ast.ModifierFlagsNone
	if isConstEnumSymbol(symbol) {
		constModifier = ast.ModifierFlagsConst
	}
	var mods *ast.ModifierList
	if constModifier != 0 {
		mods = s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(constModifier, s.b.f.NewModifier))
	}
	s.addResult(
		s.b.f.NewEnumDeclaration(
			mods,
			s.b.f.NewIdentifier(internalSymbolName),
			s.b.f.NewNodeList(members),
		),
		modifierFlags,
	)
}

func (s *symbolTableSerializationState) getEnumMemberInitializer(p *ast.Symbol) *ast.Node {
	memberDecl := core.Find(p.Declarations, ast.IsEnumMember)
	if memberDecl == nil {
		return nil
	}
	initializedValue := s.b.ch.GetConstantValue(memberDecl)
	if initializedValue == nil {
		return nil
	}
	switch v := initializedValue.(type) {
	case string:
		return s.b.f.NewStringLiteral(v, 0)
	case jsnum.Number:
		return s.b.f.NewNumericLiteral(v.String(), 0)
	}
	return nil
}

func (s *symbolTableSerializationState) serializeModule(symbol *ast.Symbol, symbolName string, modifierFlags ast.ModifierFlags) {
	members := s.getNamespaceMembersForSerialization(symbol)
	expanding := isExpanding(s.b.ctx)

	// Split NS members up by declaration - members whose parent symbol is the ns symbol vs those whose is not (but were added in later via merging)
	var realMembers, mergedMembers []*ast.Symbol
	for _, m := range members {
		if (m.Parent != nil && m.Parent == symbol) || expanding {
			realMembers = append(realMembers, m)
		} else {
			mergedMembers = append(mergedMembers, m)
		}
	}

	// TODO: `suppressNewPrivateContext` is questionable - we need to simply be emitting privates in whatever scope they were declared in, rather
	// than whatever scope we traverse to them in. That's a bit of a complex rewrite, since we're not _actually_ tracking privates at all in advance,
	// so we don't even have placeholders to fill in.
	if len(realMembers) > 0 || expanding {
		var localName *ast.Node
		if expanding {
			// Use the same name as symbol display.
			oldFlags := s.b.ctx.flags
			s.b.ctx.flags |= nodebuilder.FlagsWriteTypeParametersInQualifiedName | nodebuilder.Flags(SymbolFormatFlagsUseOnlyExternalAliasing)
			localName = s.b.symbolToNode(symbol, ast.SymbolFlagsAll)
			s.b.ctx.flags = oldFlags
		} else {
			localName = s.b.f.NewIdentifier(s.getInternalSymbolName(symbol, symbolName))
		}
		s.serializeAsNamespaceDeclaration(realMembers, localName, modifierFlags, false)
	}
	// Handle merged members as variable/namespace if needed
	if len(mergedMembers) > 0 && !expanding {
		if symbol.Flags&(ast.SymbolFlagsValueModule|ast.SymbolFlagsNamespaceModule) == 0 || (symbol.Exports != nil && len(symbol.Exports) != 0) {
			return
		}
		props := core.Filter(s.b.ch.getPropertiesOfType(s.b.ch.getTypeOfSymbol(symbol)), func(p *ast.Symbol) bool { return s.isNamespaceMember(p) })
		localName := s.getInternalSymbolName(symbol, symbolName)
		s.b.ctx.approximateLength += len(localName)
		s.serializeAsNamespaceDeclaration(props, s.b.f.NewIdentifier(localName), modifierFlags, true)
	}
}

func (s *symbolTableSerializationState) serializeAsNamespaceDeclaration(props []*ast.Symbol, localName *ast.Node, modifierFlags ast.ModifierFlags, suppressNewPrivateContext bool) {
	expanding := isExpanding(s.b.ctx)
	// Use "namespace" for identifier names, "module" for string literal names (ambient modules)
	keyword := ast.KindNamespaceKeyword
	if !ast.IsIdentifier(localName) {
		keyword = ast.KindModuleKeyword
	}
	if len(props) > 0 {
		s.b.ctx.approximateLength += 14
		// Separate local vs remote props
		// handle remote props first - we need to make an `import` declaration that points at the module containing each remote
		// prop in the outermost scope
		// TODO: implement handling for remote props - should be difficult to trigger, as only interesting cross-file js merges should make this possible
		var localProps []*ast.Symbol
		for _, p := range props {
			if len(p.Declarations) == 0 || core.Some(p.Declarations, func(d *ast.Node) bool {
				return ast.GetSourceFileOfNode(d) == ast.GetSourceFileOfNode(s.b.ctx.enclosingDeclaration)
			}) || expanding {
				localProps = append(localProps, p)
			}
		}

		// Add a namespace
		// Create namespace as non-synthetic so it is usable as an enclosing declaration
		fakespace := s.b.f.NewModuleDeclaration(nil, keyword, localName, s.b.f.NewModuleBlock(s.b.f.NewNodeList(nil)))
		fakespace.Flags &^= ast.NodeFlagsSynthesized
		fakespace.Parent = s.b.ctx.enclosingDeclaration
		// Set locals and symbol
		localTable := make(ast.SymbolTable)
		for _, p := range props {
			localTable[p.Name] = p
		}
		fakespace.LocalsContainerData().Locals = localTable
		if len(props) > 0 && props[0].Parent != nil {
			fakespace.DeclarationData().Symbol = props[0].Parent
		}

		oldResults := s.results
		s.results = nil
		oldAddingDeclare := s.addingDeclare
		s.addingDeclare = false
		oldEnclosingDeclaration := s.b.ctx.enclosingDeclaration
		s.b.ctx.enclosingDeclaration = fakespace

		localSymbolTable := make(ast.SymbolTable)
		for _, p := range localProps {
			localSymbolTable[p.Name] = p
		}
		s.visitSymbolTable(&localSymbolTable, suppressNewPrivateContext, true)

		s.b.ctx.enclosingDeclaration = oldEnclosingDeclaration
		s.addingDeclare = oldAddingDeclare
		declarations := s.results
		s.results = oldResults

		// Strip export modifiers if all declarations are exported
		allExported := len(declarations) > 0 && core.Every(declarations, func(d *ast.Node) bool {
			return ast.HasSyntacticModifier(d, ast.ModifierFlagsExport)
		})
		if allExported {
			declarations = core.Map(declarations, func(d *ast.Node) *ast.Node {
				return s.removeExportModifier(d)
			})
		}

		// replace namespace with synthetic version
		fakespace = s.b.f.UpdateModuleDeclaration(fakespace.AsModuleDeclaration(), fakespace.Modifiers(), keyword, fakespace.Name(), s.b.f.NewModuleBlock(s.b.f.NewNodeList(declarations)))
		fakespace.Parent = s.b.ctx.enclosingDeclaration
		s.addResult(fakespace, modifierFlags) // namespaces can never be default exported
	} else if expanding {
		s.b.ctx.approximateLength += 14
		s.addResult(
			s.b.f.NewModuleDeclaration(nil, keyword, localName, s.b.f.NewModuleBlock(s.b.f.NewNodeList(nil))),
			modifierFlags,
		)
	}
}

func (s *symbolTableSerializationState) removeExportModifier(node *ast.Node) *ast.Node {
	if !ast.CanHaveModifiers(node) {
		return node
	}
	flags := node.ModifierFlags() &^ ast.ModifierFlagsExport
	modifiers := ast.CreateModifiersFromModifierFlags(flags, s.b.f.NewModifier)
	return ast.ReplaceModifiers(s.b.f, node, s.b.f.NewModifierList(modifiers))
}

func (s *symbolTableSerializationState) serializePropertySymbolsForInterface(props []*ast.Symbol, baseType *Type) []*ast.Node {
	var elements []*ast.Node
	for i, prop := range props {
		if s.b.checkTruncationLengthIfExpanding() && (i+2 < len(props)-1) {
			s.b.ctx.out.Truncated = true
			elements = append(elements, s.createTruncationProperty(fmt.Sprintf("... %d more ... ", len(props)-i), false))
			result := s.serializePropertySymbolForInterface(props[len(props)-1], baseType)
			elements = append(elements, result...)
			break
		}
		s.b.ctx.approximateLength += 1
		result := s.serializePropertySymbolForInterface(prop, baseType)
		elements = append(elements, result...)
	}
	return elements
}

func (s *symbolTableSerializationState) serializePropertySymbolsForClass(props []*ast.Symbol, isStatic bool, baseType *Type) []*ast.Node {
	var elements []*ast.Node
	for i, prop := range props {
		if s.b.checkTruncationLengthIfExpanding() && (i+2 < len(props)-1) {
			s.b.ctx.out.Truncated = true
			elements = append(elements, s.createTruncationProperty(fmt.Sprintf("... %d more ... ", len(props)-i), true))
			result := s.serializePropertySymbolForClass(props[len(props)-1], isStatic, baseType)
			elements = append(elements, result...)
			break
		}
		s.b.ctx.approximateLength += 1
		result := s.serializePropertySymbolForClass(prop, isStatic, baseType)
		elements = append(elements, result...)
	}
	return elements
}

func (s *symbolTableSerializationState) createTruncationProperty(text string, isClass bool) *ast.Node {
	if isClass {
		return s.b.f.NewPropertyDeclaration(nil, s.b.f.NewStringLiteral(text, 0), nil, nil, nil)
	}
	return s.b.f.NewPropertySignatureDeclaration(nil, s.b.f.NewStringLiteral(text, 0), nil, nil, nil)
}

func (s *symbolTableSerializationState) createTruncationStatement(text string) *ast.Node {
	return s.b.f.NewExpressionStatement(s.b.f.NewIdentifier(text))
}

func (s *symbolTableSerializationState) serializePropertySymbolForInterface(p *ast.Symbol, baseType *Type) []*ast.Node {
	return s.makeSerializePropertySymbol(p, false, baseType, false)
}

func (s *symbolTableSerializationState) serializePropertySymbolForClass(p *ast.Symbol, isStatic bool, baseType *Type) []*ast.Node {
	return s.makeSerializePropertySymbol(p, isStatic, baseType, true)
}

func (s *symbolTableSerializationState) makeSerializePropertySymbol(p *ast.Symbol, isStatic bool, baseType *Type, isClass bool) []*ast.Node {
	modFlags := getDeclarationModifierFlagsFromSymbol(p)
	omitType := modFlags&ast.ModifierFlagsPrivate != 0 && !isExpanding(s.b.ctx)

	if isStatic && p.Flags&(ast.SymbolFlagsType|ast.SymbolFlagsNamespace|ast.SymbolFlagsAlias) != 0 {
		return nil
	}

	if p.Flags&ast.SymbolFlagsPrototype != 0 || p.Name == "constructor" {
		return nil
	}
	if baseType != nil {
		baseProp := s.b.ch.getPropertyOfType(baseType, p.Name)
		if baseProp != nil &&
			s.b.ch.isReadonlySymbol(baseProp) == s.b.ch.isReadonlySymbol(p) &&
			(p.Flags&ast.SymbolFlagsOptional) == (baseProp.Flags&ast.SymbolFlagsOptional) &&
			s.b.ch.isTypeIdenticalTo(s.b.ch.getTypeOfSymbol(p), s.b.ch.getTypeOfPropertyOfType(baseType, p.Name)) {
			return nil
		}
	}

	flag := modFlags &^ ast.ModifierFlagsAsync
	if isStatic {
		flag |= ast.ModifierFlagsStatic
	}
	name := s.b.getPropertyNameNodeForSymbol(p)

	if p.Flags&ast.SymbolFlagsAccessor != 0 && isClass {
		var result []*ast.Node
		if p.Flags&ast.SymbolFlagsSetAccessor != 0 {
			setterDecl := core.Find(p.Declarations, ast.IsSetAccessorDeclaration)
			s.b.ctx.approximateLength += 7
			result = append(result, s.b.f.NewSetAccessorDeclaration(
				s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(flag, s.b.f.NewModifier)),
				name,
				nil,
				s.b.f.NewNodeList([]*ast.Node{
					s.b.f.NewParameterDeclaration(nil, nil, s.b.f.NewIdentifier("value"), nil,
						core.IfElse(omitType, nil, s.b.serializeTypeForDeclaration(setterDecl, nil, p, true)),
						nil),
				}),
				nil, nil, nil,
			))
		}
		if p.Flags&ast.SymbolFlagsGetAccessor != 0 {
			getterDecl := core.Find(p.Declarations, ast.IsGetAccessorDeclaration)
			_ = getterDecl
			s.b.ctx.approximateLength += 8
			result = append(result, s.b.f.NewGetAccessorDeclaration(
				s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(flag, s.b.f.NewModifier)),
				name,
				nil,
				s.b.f.NewNodeList(nil),
				core.IfElse(omitType, nil, s.b.serializeTypeForDeclaration(nil, nil, p, true)),
				nil, nil,
			))
		}
		return result
	}

	if p.Flags&(ast.SymbolFlagsProperty|ast.SymbolFlagsVariable|ast.SymbolFlagsAccessor) != 0 {
		mf := flag
		if s.b.ch.isReadonlySymbol(p) {
			mf |= ast.ModifierFlagsReadonly
		}
		s.b.ctx.approximateLength += 2
		var typeNode *ast.Node
		if !omitType {
			typeNode = s.b.serializeTypeForDeclaration(core.Find(p.Declarations, ast.IsSetAccessorDeclaration), nil, p, true)
		}
		var questionToken *ast.Node
		if p.Flags&ast.SymbolFlagsOptional != 0 {
			questionToken = s.b.f.NewToken(ast.KindQuestionToken)
		}
		if isClass {
			return []*ast.Node{s.b.f.NewPropertyDeclaration(
				s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(mf, s.b.f.NewModifier)),
				name, questionToken, typeNode, nil,
			)}
		}
		return []*ast.Node{s.b.f.NewPropertySignatureDeclaration(
			s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(mf, s.b.f.NewModifier)),
			name, questionToken, typeNode, nil,
		)}
	}

	if p.Flags&(ast.SymbolFlagsMethod|ast.SymbolFlagsFunction) != 0 {
		t := s.b.ch.getTypeOfSymbol(p)
		signatures := s.b.ch.getSignaturesOfType(t, SignatureKindCall)
		if omitType {
			mf := flag
			if s.b.ch.isReadonlySymbol(p) {
				mf |= ast.ModifierFlagsReadonly
			}
			s.b.ctx.approximateLength += 1
			var questionToken *ast.Node
			if p.Flags&ast.SymbolFlagsOptional != 0 {
				questionToken = s.b.f.NewToken(ast.KindQuestionToken)
			}
			if isClass {
				return []*ast.Node{s.b.f.NewPropertyDeclaration(
					s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(mf, s.b.f.NewModifier)),
					name, questionToken, nil, nil,
				)}
			}
			return []*ast.Node{s.b.f.NewPropertySignatureDeclaration(
				s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(mf, s.b.f.NewModifier)),
				name, questionToken, nil, nil,
			)}
		}

		var result []*ast.Node
		methodKind := core.IfElse(isClass, ast.KindMethodDeclaration, ast.KindMethodSignature)
		for _, sig := range signatures {
			s.b.ctx.approximateLength += 1
			var questionToken *ast.Node
			if p.Flags&ast.SymbolFlagsOptional != 0 {
				questionToken = s.b.f.NewToken(ast.KindQuestionToken)
			}
			var modList *ast.ModifierList
			if flag != 0 {
				modList = s.b.f.NewModifierList(ast.CreateModifiersFromModifierFlags(flag, s.b.f.NewModifier))
			}
			var modifiers []*ast.Node
			if modList != nil {
				modifiers = modList.Nodes
			}
			decl := s.b.signatureToSignatureDeclarationHelper(sig, methodKind, &SignatureToSignatureDeclarationOptions{
				name:          name,
				questionToken: questionToken,
				modifiers:     modifiers,
			})
			result = append(result, decl)
		}
		return result
	}

	return nil
}

func (s *symbolTableSerializationState) serializeSignatures(kind SignatureKind, input *Type, baseType *Type, outputKind ast.Kind) []*ast.Node {
	signatures := s.b.ch.getSignaturesOfType(input, kind)
	if kind == SignatureKindConstruct {
		if baseType == nil && core.Every(signatures, func(sig *Signature) bool { return len(sig.parameters) == 0 }) {
			return nil
		}
		if baseType != nil {
			baseSigs := s.b.ch.getSignaturesOfType(baseType, SignatureKindConstruct)
			if len(baseSigs) == 0 && core.Every(signatures, func(sig *Signature) bool { return len(sig.parameters) == 0 }) {
				return nil
			}
			if len(baseSigs) == len(signatures) {
				failed := false
				for i := range baseSigs {
					if s.b.ch.compareSignaturesIdentical(signatures[i], baseSigs[i], false, false, true, s.b.ch.compareTypesIdentical) != TernaryTrue {
						failed = true
						break
					}
				}
				if !failed {
					return nil
				}
			}
			var privateProtected ast.ModifierFlags
			for _, sig := range signatures {
				if sig.declaration != nil {
					privateProtected |= sig.declaration.ModifierFlags() & (ast.ModifierFlagsPrivate | ast.ModifierFlagsProtected)
				}
			}
			if privateProtected != 0 {
				modifiers := ast.CreateModifiersFromModifierFlags(privateProtected, s.b.f.NewModifier)
				return []*ast.Node{s.b.f.NewConstructorDeclaration(
					s.b.f.NewModifierList(modifiers),
					nil, s.b.f.NewNodeList(nil), nil, nil, nil,
				)}
			}
		}
	}

	var result []*ast.Node
	for _, sig := range signatures {
		s.b.ctx.approximateLength += 1
		decl := s.b.signatureToSignatureDeclarationHelper(sig, outputKind, nil)
		result = append(result, decl)
	}
	return result
}

func (s *symbolTableSerializationState) serializeIndexSignatures(input *Type, baseType *Type) []*ast.Node {
	var result []*ast.Node
	for _, info := range s.b.ch.getIndexInfosOfType(input) {
		if baseType != nil {
			baseInfo := s.b.ch.getIndexInfoOfType(baseType, info.keyType)
			if baseInfo != nil && s.b.ch.isTypeIdenticalTo(info.valueType, baseInfo.valueType) {
				continue
			}
		}
		result = append(result, s.b.indexInfoToIndexSignatureDeclarationHelper(info, nil))
	}
	return result
}

func (s *symbolTableSerializationState) serializeBaseType(t *Type, staticType *Type, rootName string) *ast.Node {
	ref := s.trySerializeAsTypeReference(t, ast.SymbolFlagsValue)
	if ref != nil {
		return ref
	}
	tempName := s.getUnusedName(rootName+"_base", nil)
	stmt := s.b.f.NewVariableStatement(
		nil,
		s.b.f.NewVariableDeclarationList(
			ast.NodeFlagsConst,
			s.b.f.NewNodeList([]*ast.Node{
				s.b.f.NewVariableDeclaration(s.b.f.NewIdentifier(tempName), nil, s.b.typeToTypeNode(staticType), nil),
			}),
		),
	)
	s.addResult(stmt, ast.ModifierFlagsNone)
	return s.b.f.NewExpressionWithTypeArguments(s.b.f.NewIdentifier(tempName), nil)
}

func (s *symbolTableSerializationState) trySerializeAsTypeReference(t *Type, flags ast.SymbolFlags) *ast.Node {
	var typeArgs []*ast.Node
	var reference *ast.Node

	if t.Target() != nil && s.b.ch.IsSymbolAccessibleByFlags(t.Target().symbol, s.enclosingDeclaration, flags) {
		typeArgs = core.Map(s.b.ch.getTypeArguments(t), func(arg *Type) *ast.Node { return s.b.typeToTypeNode(arg) })
		reference = s.b.symbolToExpression(t.Target().symbol, ast.SymbolFlagsType)
	} else if t.symbol != nil && s.b.ch.IsSymbolAccessibleByFlags(t.symbol, s.enclosingDeclaration, flags) {
		reference = s.b.symbolToExpression(t.symbol, ast.SymbolFlagsType)
	}
	if reference != nil {
		return s.b.f.NewExpressionWithTypeArguments(reference, s.b.f.NewNodeList(typeArgs))
	}
	return nil
}

func (s *symbolTableSerializationState) serializeImplementedType(t *Type) *ast.Node {
	ref := s.trySerializeAsTypeReference(t, ast.SymbolFlagsType)
	if ref != nil {
		return ref
	}
	if t.symbol != nil {
		return s.b.f.NewExpressionWithTypeArguments(s.b.symbolToExpression(t.symbol, ast.SymbolFlagsType), nil)
	}
	return nil
}

func (s *symbolTableSerializationState) serializeVariableOrProperty(symbol *ast.Symbol, symbolName string, isPrivate bool, needsPostExportDefault bool, modifierFlags ast.ModifierFlags, propertyAsAlias bool) {
	if propertyAsAlias {
		s.serializeMaybeAliasAssignment(symbol)
		return
	}
	_ = s.b.ch.getTypeOfSymbol(symbol)
	localName := s.getInternalSymbolName(symbol, symbolName)

	var flags ast.NodeFlags
	if symbol.Flags&ast.SymbolFlagsBlockScopedVariable == 0 {
		if symbol.Parent != nil && symbol.Parent.ValueDeclaration != nil && ast.IsSourceFile(symbol.Parent.ValueDeclaration) {
			flags = ast.NodeFlagsConst
		}
	} else if isConstantVariable(symbol) {
		flags = ast.NodeFlagsConst
	} else {
		flags = ast.NodeFlagsLet
	}

	var name string
	if !needsPostExportDefault && symbol.Flags&ast.SymbolFlagsProperty == 0 {
		name = localName
	} else {
		name = s.getUnusedName(localName, symbol)
	}

	s.b.ctx.approximateLength += 7 + len(name)
	stmt := s.b.f.NewVariableStatement(
		nil,
		s.b.f.NewVariableDeclarationList(
			flags,
			s.b.f.NewNodeList([]*ast.Node{
				s.b.f.NewVariableDeclaration(s.b.f.NewIdentifier(name), nil, s.b.serializeTypeForDeclaration(nil, nil, symbol, true), nil),
			}),
		),
	)
	if name != localName {
		s.addResult(stmt, modifierFlags&^ast.ModifierFlagsExport)
		if !isPrivate {
			s.b.ctx.approximateLength += 16 + len(name) + len(localName)
			s.results = append(s.results, s.b.f.NewExportDeclaration(
				nil,
				false,
				s.b.f.NewNamedExports(s.b.f.NewNodeList([]*ast.Node{
					s.b.f.NewExportSpecifier(false, s.b.f.NewIdentifier(name), s.b.f.NewIdentifier(localName)),
				})),
				nil,
				nil,
			))
		}
	} else {
		s.addResult(stmt, modifierFlags)
	}
}

func (s *symbolTableSerializationState) serializeAsFunctionNamespaceMerge(t *Type, symbol *ast.Symbol, localName string, modifierFlags ast.ModifierFlags) {
	signatures := s.b.ch.getSignaturesOfType(t, SignatureKindCall)
	for _, sig := range signatures {
		s.b.ctx.approximateLength += 1 // ;
		// Each overload becomes a separate function declaration, in order
		decl := s.b.signatureToSignatureDeclarationHelper(sig, ast.KindFunctionDeclaration, &SignatureToSignatureDeclarationOptions{
			name: s.b.f.NewIdentifier(localName),
		})
		s.addResult(decl, modifierFlags)
	}
	// Module symbol emit will take care of module-y members, provided it has exports
	if symbol.Flags&(ast.SymbolFlagsValueModule|ast.SymbolFlagsNamespaceModule) != 0 && symbol.Exports != nil && len(symbol.Exports) != 0 {
		return // module emit will handle it
	}
	props := core.Filter(s.b.ch.getPropertiesOfType(t), func(p *ast.Symbol) bool { return s.isNamespaceMember(p) })
	s.b.ctx.approximateLength += len(localName)
	s.serializeAsNamespaceDeclaration(props, s.b.f.NewIdentifier(localName), modifierFlags, true)
}

func (s *symbolTableSerializationState) serializeAsAlias(symbol *ast.Symbol, localName string, modifierFlags ast.ModifierFlags) {
	node := s.b.ch.getDeclarationOfAliasSymbol(symbol)
	if node == nil {
		return
	}
	target := s.b.ch.getMergedSymbol(s.b.ch.getTargetOfAliasDeclaration(node))
	if target == nil {
		return
	}
	targetName := s.getInternalSymbolName(target, target.Name)
	s.includePrivateSymbol(target)

	switch node.Kind {
	case ast.KindExportSpecifier:
		specifier := node.Parent.Parent.AsExportDeclaration().ModuleSpecifier
		var specifierExpr *ast.Node
		if specifier != nil && ast.IsStringLiteralLike(specifier) {
			specifierExpr = s.b.f.NewStringLiteral(specifier.Text(), 0)
		}
		s.serializeExportSpecifier(symbol.Name, core.IfElse(specifier != nil, target.Name, targetName), specifierExpr)
	case ast.KindExportAssignment:
		s.serializeMaybeAliasAssignment(symbol)
	default:
		// For other alias kinds, emit an import = or export specifier
		s.serializeExportSpecifier(localName, targetName, nil)
	}
}

func (s *symbolTableSerializationState) serializeExportSpecifier(localName string, targetName string, specifier *ast.Node) {
	s.b.ctx.approximateLength += 16 + len(localName)
	var propertyName *ast.Node
	if localName != targetName {
		propertyName = s.b.f.NewIdentifier(targetName)
		s.b.ctx.approximateLength += len(targetName)
	}
	s.results = append(s.results, s.b.f.NewExportDeclaration(
		nil,
		false,
		s.b.f.NewNamedExports(s.b.f.NewNodeList([]*ast.Node{
			s.b.f.NewExportSpecifier(false, propertyName, s.b.f.NewIdentifier(localName)),
		})),
		specifier,
		nil,
	))
}

func (s *symbolTableSerializationState) serializeMaybeAliasAssignment(symbol *ast.Symbol) bool {
	if symbol.Flags&ast.SymbolFlagsPrototype != 0 {
		return false
	}
	name := symbol.Name
	isExportEquals := name == ast.InternalSymbolNameExportEquals
	isDefault := name == ast.InternalSymbolNameDefault
	isExportAssignmentCompatibleSymbolName := isExportEquals || isDefault

	// serialize as an anonymous property declaration
	varName := s.getUnusedName(name, symbol)
	// We have to use `getWidenedType` here since the object within a json file is unwidened within the file
	// (Unwidened types can only exist in expression contexts and should never be serialized)
	typeToSerialize := s.b.ch.getWidenedType(s.b.ch.getTypeOfSymbol(s.b.ch.getMergedSymbol(symbol)))

	// Inside a module/namespace declaration, use `let` for non-getter-only accessors (matching Strada behavior).
	// Otherwise use `const`.
	flags := ast.NodeFlagsConst
	if s.b.ctx.enclosingDeclaration != nil && s.b.ctx.enclosingDeclaration.Kind == ast.KindModuleDeclaration &&
		(symbol.Flags&ast.SymbolFlagsAccessor == 0 || symbol.Flags&ast.SymbolFlagsSetAccessor != 0) {
		flags = ast.NodeFlagsLet
	}
	s.b.ctx.approximateLength += len(varName) + 5
	stmt := s.b.f.NewVariableStatement(
		nil,
		s.b.f.NewVariableDeclarationList(
			flags,
			s.b.f.NewNodeList([]*ast.Node{
				s.b.f.NewVariableDeclaration(s.b.f.NewIdentifier(varName), nil, s.b.serializeTypeForDeclaration(nil, typeToSerialize, symbol, true), nil),
			}),
		),
	)
	s.addResult(stmt, core.IfElse(name == varName, ast.ModifierFlagsExport, ast.ModifierFlagsNone))

	if isExportAssignmentCompatibleSymbolName {
		s.b.ctx.approximateLength += len(varName) + 10
		s.results = append(s.results, s.b.f.NewExportAssignment(nil, isExportEquals, nil, s.b.f.NewIdentifier(varName)))
	}
	return true
}

// --- Helper functions ---

func isExpanding(ctx *NodeBuilderContext) bool {
	return ctx.maxExpansionDepth != -1
}

func isHashPrivate(s *ast.Symbol) bool {
	return s.ValueDeclaration != nil && s.ValueDeclaration.Name() != nil && ast.IsPrivateIdentifier(s.ValueDeclaration.Name())
}

func (s *symbolTableSerializationState) getNamespaceMembersForSerialization(symbol *ast.Symbol) []*ast.Symbol {
	var exports []*ast.Symbol
	for _, sym := range s.b.ch.getExportsOfSymbol(symbol) {
		exports = append(exports, sym)
	}
	merged := s.b.ch.getMergedSymbol(symbol)
	if merged != symbol {
		membersSet := make(map[ast.SymbolId]*ast.Symbol)
		for _, e := range exports {
			membersSet[ast.GetSymbolId(e)] = e
		}
		for _, exported := range s.b.ch.getExportsOfSymbol(merged) {
			resolved := s.b.ch.resolveSymbol(exported)
			if s.b.ch.getSymbolFlags(resolved)&ast.SymbolFlagsValue == 0 {
				if _, ok := membersSet[ast.GetSymbolId(exported)]; !ok {
					membersSet[ast.GetSymbolId(exported)] = exported
					exports = append(exports, exported)
				}
			}
		}
	}
	return core.Filter(exports, func(m *ast.Symbol) bool {
		return s.isNamespaceMember(m) && scanner.IsIdentifierText(m.Name, core.LanguageVariantStandard)
	})
}

func (s *symbolTableSerializationState) isTypeOnlyNamespace(symbol *ast.Symbol) bool {
	return core.Every(s.getNamespaceMembersForSerialization(symbol), func(m *ast.Symbol) bool {
		resolved := s.b.ch.resolveSymbol(m)
		return s.b.ch.getSymbolFlags(resolved)&ast.SymbolFlagsValue == 0
	})
}

func (s *symbolTableSerializationState) isNamespaceMember(p *ast.Symbol) bool {
	return p.Flags&(ast.SymbolFlagsType|ast.SymbolFlagsNamespace|ast.SymbolFlagsAlias) != 0 ||
		!(p.Flags&ast.SymbolFlagsPrototype != 0 || p.Name == "prototype" || (p.ValueDeclaration != nil && ast.HasStaticModifier(p.ValueDeclaration) && ast.IsClassLike(p.ValueDeclaration.Parent)))
}

func (s *symbolTableSerializationState) isTypeRepresentableAsFunctionNamespaceMerge(t *Type, symbol *ast.Symbol) bool {
	// Simplified check: can this type be represented as a function + namespace merge?
	signatures := s.b.ch.getSignaturesOfType(t, SignatureKindCall)
	if len(signatures) == 0 {
		return false
	}
	// Check if the type has only call signatures and properties (no construct signatures, no string/number index)
	constructSignatures := s.b.ch.getSignaturesOfType(t, SignatureKindConstruct)
	if len(constructSignatures) > 0 {
		return false
	}
	indexInfos := s.b.ch.getIndexInfosOfType(t)
	if len(indexInfos) > 0 {
		return false
	}
	return true
}

func (s *symbolTableSerializationState) getNonInheritedProperties(t *Type, baseTypes []*Type, properties []*ast.Symbol) []*ast.Symbol {
	if len(baseTypes) == 0 {
		return properties
	}
	seen := make(map[string]*ast.Symbol)
	for _, p := range properties {
		seen[p.Name] = p
	}
	for _, base := range baseTypes {
		baseWithThis := s.b.ch.getTypeWithThisArgument(base, s.b.ch.getTargetType(t).AsInterfaceType().thisType, false)
		baseProps := s.b.ch.getPropertiesOfType(baseWithThis)
		for _, prop := range baseProps {
			if existing, ok := seen[prop.Name]; ok && prop.Parent == existing.Parent {
				delete(seen, prop.Name)
			}
		}
	}
	return core.Filter(properties, func(p *ast.Symbol) bool {
		_, ok := seen[p.Name]
		return ok
	})
}

func (s *symbolTableSerializationState) getImplementsTypes(classType *Type) []*Type {
	var result []*Type
	if classType.symbol == nil {
		return result
	}
	for _, declaration := range classType.symbol.Declarations {
		implementsTypeNodes := ast.GetImplementsTypeNodes(declaration)
		if implementsTypeNodes == nil {
			continue
		}
		for _, node := range implementsTypeNodes {
			implementsType := s.b.ch.getTypeFromTypeNode(node)
			if !s.b.ch.isErrorType(implementsType) {
				result = append(result, implementsType)
			}
		}
	}
	return result
}

func (s *symbolTableSerializationState) getUnusedName(input string, symbol *ast.Symbol) string {
	if symbol != nil {
		id := ast.GetSymbolId(symbol)
		if name, ok := s.remappedSymbolNames[id]; ok {
			return name
		}
	}
	if symbol != nil {
		input = s.getNameCandidateWorker(symbol, input)
	}
	i := 0
	original := input
	for s.usedSymbolNames.Has(input) {
		i++
		input = original + "_" + strconv.Itoa(i)
	}
	s.usedSymbolNames.Add(input)
	if symbol != nil {
		s.remappedSymbolNames[ast.GetSymbolId(symbol)] = input
	}
	return input
}

func (s *symbolTableSerializationState) getNameCandidateWorker(symbol *ast.Symbol, localName string) string {
	if localName == ast.InternalSymbolNameDefault || localName == ast.InternalSymbolNameClass || localName == ast.InternalSymbolNameFunction {
		restoreFlags := s.b.saveRestoreFlags()
		s.b.ctx.flags |= nodebuilder.FlagsInInitialEntityName
		nameCandidate := s.b.getNameOfSymbolAsWritten(symbol)
		restoreFlags()
		if len(nameCandidate) > 0 && (nameCandidate[0] == '\'' || nameCandidate[0] == '"') {
			localName = strings.Trim(nameCandidate, "'\"")
		} else {
			localName = nameCandidate
		}
	}
	if localName == ast.InternalSymbolNameDefault {
		localName = "_default"
	} else if localName == ast.InternalSymbolNameExportEquals {
		localName = "_exports"
	}
	if !scanner.IsIdentifierText(localName, core.LanguageVariantStandard) || ast.IsNonContextualKeyword(scanner.StringToToken(localName)) {
		localName = "_" + strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return '_'
		}, localName)
	}
	return localName
}

func (s *symbolTableSerializationState) getInternalSymbolName(symbol *ast.Symbol, localName string) string {
	id := ast.GetSymbolId(symbol)
	if name, ok := s.remappedSymbolNames[id]; ok {
		return name
	}
	localName = s.getNameCandidateWorker(symbol, localName)
	s.remappedSymbolNames[id] = localName
	return localName
}

func isConstantVariable(symbol *ast.Symbol) bool {
	return symbol.Flags&ast.SymbolFlagsBlockScopedVariable != 0 &&
		symbol.ValueDeclaration != nil &&
		ast.IsVarConst(symbol.ValueDeclaration)
}
