#!/usr/bin/env node

import * as ts from 'typescript';
import * as fs from 'fs';
import * as url from 'url';

function main() {
    const args = process.argv.slice(2);
    
    if (args.length < 1) {
        console.log('Usage: node parse_ts_file.ts <path-to-ts-file>');
        process.exit(1);
    }

    const filename = args[0];

    if (!filename.endsWith('.ts')) {
        console.error(`File must have .ts extension: ${filename}`);
        process.exit(1);
    }

    if (!fs.existsSync(filename)) {
        console.error(`File does not exist: ${filename}`);
        process.exit(1);
    }

    const sourceText = fs.readFileSync(filename, 'utf8');

    const sourceFile = ts.createSourceFile(
        filename,
        sourceText,
        ts.ScriptTarget.Latest,
        /*setParentNodes*/ true,
    );
    const modifiedSourceFile = addLogStatementsToFile(sourceFile);

    const printer = ts.createPrinter({
        newLine: ts.NewLineKind.LineFeed,
    });

    const modifiedCode = printer.printFile(modifiedSourceFile);

    fs.writeFileSync(filename, modifiedCode);
    console.log(`Modified file written to: ${filename}`);
}

function addLogStatementsToFile(sourceFile: ts.SourceFile): ts.SourceFile {
    const newFile = ts.visitNode(sourceFile, visitNode) as ts.SourceFile;
    return ts.factory.updateSourceFile(newFile, [...newFile.statements, createLogHelper()]);
}

function visitNode<T extends ts.Node>(node: T): T {
    if (isFunctionLikeDeclaration(node)) {
        return addLogToFunction(node, ts.getNameOfDeclaration(node)?.getText() || '<anonymous>');
    }
    if (ts.isReturnStatement(node)) {
        const fn = ts.findAncestor(node, isFunctionLikeDeclaration);
        const functionName = fn ? (ts.getNameOfDeclaration(fn)?.getText() || '<anonymous>') : '<unknown>';
        return ts.factory.updateReturnStatement(
            node,
            ts.factory.createCallExpression(
                ts.factory.createIdentifier("logReturn"),
                /*typeArguments*/ undefined,
                /*arguments*/ [
                    ts.factory.createStringLiteral(functionName),
                    node.expression || ts.factory.createIdentifier("undefined"),
                ]),
        ) as ts.Node as T;
    }
    return ts.visitEachChild(node, visitNode, undefined);
}

function isFunctionLikeDeclaration(node: ts.Node): node is ts.FunctionLikeDeclaration {
    switch (node.kind) {
        case ts.SyntaxKind.FunctionDeclaration:
        case ts.SyntaxKind.MethodDeclaration:
        case ts.SyntaxKind.Constructor:
        case ts.SyntaxKind.GetAccessor:
        case ts.SyntaxKind.SetAccessor:
        case ts.SyntaxKind.FunctionExpression:
        case ts.SyntaxKind.ArrowFunction:
            return true;
        default:
            return false;
    }
}

function addLogToFunction<T extends ts.FunctionLikeDeclaration>(node: T, functionName: string): T {
    if (!node.body || !ts.isBlock(node.body)) {
        return node;
    }

    const logStatement = createLogStatement(functionName, /*enter*/ true);
    const newStatements = [logStatement, ...node.body.statements.map(stmt => visitNode(stmt))];
    
    const newBody = ts.factory.updateBlock(node.body, newStatements);

    return updateFunctionLikeDeclaration(node, newBody)
}

function updateFunctionLikeDeclaration<T extends ts.FunctionLikeDeclaration>(node: T, newBody: ts.Block): T {
    switch (node.kind) {
        case ts.SyntaxKind.FunctionDeclaration:
            return ts.factory.updateFunctionDeclaration(
                node,
                node.modifiers,
                node.asteriskToken,
                node.name,
                node.typeParameters,
                node.parameters,
                node.type,
                newBody,
            ) as T;
        case ts.SyntaxKind.MethodDeclaration:
            return ts.factory.updateMethodDeclaration(
                node,
                node.modifiers,
                node.asteriskToken,
                node.name,
                node.questionToken,
                node.typeParameters,
                node.parameters,
                node.type,
                newBody,
            ) as T;
        case ts.SyntaxKind.Constructor:
            return ts.factory.updateConstructorDeclaration(
                node,
                node.modifiers,
                node.parameters,
                newBody,
            ) as T;
        case ts.SyntaxKind.GetAccessor:
            return ts.factory.updateGetAccessorDeclaration(
                node,
                node.modifiers,
                node.name,
                node.parameters,
                node.type,
                newBody,
            ) as T;
        case ts.SyntaxKind.SetAccessor:
            return ts.factory.updateSetAccessorDeclaration(
                node,
                node.modifiers,
                node.name,
                node.parameters,
                newBody,
            ) as T;
        case ts.SyntaxKind.FunctionExpression:
            return ts.factory.updateFunctionExpression(
                node,
                node.modifiers,
                node.asteriskToken,
                node.name,
                node.typeParameters,
                node.parameters,
                node.type,
                newBody,
            ) as T;
        case ts.SyntaxKind.ArrowFunction:
            return ts.factory.updateArrowFunction(
                node,
                node.modifiers,
                node.typeParameters,
                node.parameters,
                node.type,
                node.equalsGreaterThanToken,
                newBody,
            ) as T;
    }
}

function createLogStatement(functionName: string, enter: boolean): ts.ExpressionStatement {
    const logStmt = ts.factory.createExpressionStatement(
        ts.factory.createCallExpression(
            ts.factory.createPropertyAccessExpression(
                ts.factory.createIdentifier('console'),
                ts.factory.createIdentifier('log')
            ),
            undefined,
            [ts.factory.createStringLiteral(`${enter ? ">" : "<"} ${functionName}`)]
        )
    );
    
    return ts.addSyntheticLeadingComment(logStmt, ts.SyntaxKind.SingleLineCommentTrivia, "@ts-ignore");
}

function createLogHelper(): ts.FunctionDeclaration {
    /*
        function logReturn<T>(fnName: string, e: T): T {
            console.log(`> ${fnName}`);
            return e;
        }
    */
    return ts.factory.createFunctionDeclaration(
        /*modifiers*/ undefined,
        /*asteriskToken*/ undefined,
        ts.factory.createIdentifier('logReturn'),
        /*typeParameters*/ [ts.factory.createTypeParameterDeclaration(/*modifiers*/ undefined, 'T')],
        /*parameters*/ [
            ts.factory.createParameterDeclaration(
                undefined,
                undefined,
                ts.factory.createIdentifier("fnName"),
                undefined,
                ts.factory.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
                undefined,
            ),
            ts.factory.createParameterDeclaration(
                undefined,
                undefined,
                ts.factory.createIdentifier("e"),
                undefined,
                ts.factory.createTypeReferenceNode(ts.factory.createIdentifier("T")),
                undefined
            )
        ],
        /*type*/ ts.factory.createTypeReferenceNode(ts.factory.createIdentifier("T")),
        /*body*/ ts.factory.createBlock(
            [
                ts.addSyntheticLeadingComment(
                    ts.factory.createExpressionStatement(
                        ts.factory.createCallExpression(
                            ts.factory.createPropertyAccessExpression(
                                ts.factory.createIdentifier("console"),
                                ts.factory.createIdentifier("log")
                            ),
                            undefined,
                            [ts.factory.createTemplateExpression(
                                ts.factory.createTemplateHead("> "),
                                [ts.factory.createTemplateSpan(
                                    ts.factory.createIdentifier("fnName"),
                                    ts.factory.createTemplateTail(""),
                                )],
                            )],
                    )),
                    ts.SyntaxKind.SingleLineCommentTrivia,
                    "@ts-ignore"
                ),
                ts.factory.createReturnStatement(ts.factory.createIdentifier("e")),
            ],
            /*multiline*/ true,
        ),
    );
}

if (url.fileURLToPath(import.meta.url) == process.argv[1]) {
    main();
}
