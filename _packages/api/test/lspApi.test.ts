import assert from "node:assert"
import { Buffer } from "node:buffer"
import fs from "node:fs"
import subprocess from "node:child_process"
import path, { resolve } from "node:path"
import { afterEach, before, beforeEach, suite, test, type TestContext } from "node:test"
import url from "node:url"
import type { Range } from "vscode-languageserver-types"
import type { Type, Symbol, Node, TypeReference, GenericType, UnionType, LiteralType, IndexType, IndexedAccessType, ConditionalType, SubstitutionType, ObjectType, PseudoBigInt, BigIntLiteralType, TemplateLiteral, TemplateLiteralType, TupleType, Signature, IndexInfo, __String, Declaration, SignatureDeclaration, SourceFile, LineAndCharacter, IndexSignatureDeclaration, server } from "typescript"
// @ts-expect-error
import { TypeFlags } from "../src/typeFlags.ts"
// @ts-expect-error
import { SymbolFlags } from "../src/symbolFlags.ts"


suite("TypeScriptGoServiceGetElementTypeTest", {}, () => {
  const logFile = path.join(import.meta.dirname, "..", "tsgo.log")
  const testProjectDir = path.join(import.meta.dirname, "testProject")

  before(async () => {
    if (fs.existsSync(logFile))
      await fs.promises.unlink(logFile)
  })

  type MyContext = TestContext & {
    client?: LspClient
    openFilesContent?: Map<string, string>
  }

  beforeEach(async ctx => {
    const myCtx = ctx as MyContext

    myCtx.openFilesContent ??= new Map()

    try {
      const client = new LspClient(ctx.name, logFile, testProjectDir)
      await client.spawnAndInit();
      myCtx.client = client
    } catch (e) {
      await myCtx.client?.kill();
      myCtx.client = undefined
      throw e
    }
  })

  afterEach(async ctx => {
    const myCtx = ctx as MyContext

    myCtx.openFilesContent?.clear()

    try {
      await myCtx.client?.shutdown()
    } catch (e) {
      await myCtx.client?.kill()
      throw e
    } finally {
      myCtx.client = undefined
    }
  })

  async function getElementType(ctx: TestContext, content: string, elementText: string) {
    const myCtx = ctx as MyContext
    
    const fileName = "a.ts"
    const alreadyOpenedContent = myCtx.openFilesContent!.get(fileName)
    if (alreadyOpenedContent == null) {
      myCtx.client!.didOpen("a.ts", content)
      myCtx.openFilesContent!.set(fileName, content)
    } else if (alreadyOpenedContent !== content) {
      throw new Error(`File ${fileName} already opened with the different content`)
    }

    return await myCtx.client!.getElementType(fileName, getRange(content, elementText))
  }

  // *** Tests ***

  test("any", async ctx => {
    const type = await getElementType(ctx, "let foo", "foo")
    assert.strictEqual(type.flags, TypeFlags.Any)
  })

  test("unknown", async ctx => {
    const type = await getElementType(ctx, "let foo: unknown", "foo")
    assert.strictEqual(type.flags, TypeFlags.Unknown)
  })

  test("undefined", async ctx => {
    const type = await getElementType(ctx, "let foo: undefined", "foo")
    assert.strictEqual(type.flags, TypeFlags.Undefined)
  })

  test("number", async ctx => {
    const type = await getElementType(ctx, "declare const foo: number", "foo")
    // 8 is a converted flag, use direct instead
    assert.strictEqual(type.flags, TypeFlags.Number)
  })

  test("numberLiteral", async ctx  => {
    const type = await getElementType(ctx, "const foo = 123", "foo")
    // 256 is a converted flag
    assert.strictEqual(type.flags, TypeFlags.NumberLiteral)
    assert.strictEqual((type as LiteralType).value, 123)
  })

  test("string", async ctx => {
    const type = await getElementType(ctx, "declare const foo: string", "foo")
    assert.strictEqual(type.flags, TypeFlags.String)
  })

  test("stringLiteral", async ctx => {
    const type = await getElementType(ctx, "declare const foo: '1'", "foo")
    assert.strictEqual(type.flags, TypeFlags.StringLiteral)
    assert.strictEqual((type as LiteralType).value, "1")
  })

  test("function", async ctx => {
    const content = "function foo(x: number): string {}"
    const type = await getElementType(ctx, content, "foo")
    assert.strictEqual(type.flags, 1<<20 /* Object */)
    assert.strictEqual((type as ObjectType).objectFlags, 16 /* anonymous */)
    assert.strictEqual(getFragment(content, type.symbol.declarations![0]), content)

    const [callSignature] = await (type as TypeEx).getCallSignaturesEx()
    assert.strictEqual(callSignature.parameters[0].name, "x")
    assert.strictEqual((callSignature as SignatureEx).resolvedReturnType.flags, TypeFlags.String)
  })

  test("objectOptionalProperty", async ctx => {
    const content = "type Foo = {x?:\n  123}"

    const type = await getElementType(ctx, content, "Foo")
    assert.strictEqual(type.flags, 1<<20)
    assert.strictEqual((type as ObjectType).objectFlags, 16)

    const properties = await (type as TypeEx).getPropertiesEx()
    assert.strictEqual(properties.length, 1)

    const [xProperty] = properties
    assert.strictEqual(xProperty.escapedName, "x")
    assert.strictEqual(xProperty.name, "x")
    assert.strictEqual(xProperty.flags, SymbolFlags.Property | SymbolFlags.Optional)

    assert.strictEqual(xProperty.declarations?.length, 1)
    const [xDeclaration] = xProperty.declarations
    assert.strictEqual(getFragment(content, xDeclaration), "x?:\n  123")

    assert.deepStrictEqual(await (type as TypeEx).getCallSignaturesEx(), [])
    assert.deepStrictEqual(await (type as TypeEx).getConstructSignaturesEx(), [])
  })

  test("genericArguments", async ctx => {
    const content = "type Foo<T> = {x: T}"
    const type = await getElementType(ctx, content, "Foo<T>")

    const [tArg] = type.aliasTypeArguments!
    assert.strictEqual(tArg.flags, TypeFlags.Object)

    const {symbol} = tArg
    assert.strictEqual(symbol.name, "T")
    assert.strictEqual(symbol.escapedName, "T")
    assert.strictEqual(symbol.flags, SymbolFlags.TypeParameter)
    assert.strictEqual(getFragment(content, symbol.declarations![0]), "T")
  })

  test("conditionalType", async ctx => {
    const content = `type Foo<T> = T extends string ? '1' : 1`

    const type = await getElementType(ctx, content, "Foo<T>")
    assert.strictEqual(type.flags, 1<<26 /* Conditional */)

    const {aliasSymbol} = type
    assert.strictEqual(aliasSymbol!.flags, TypeFlags.Object)
    assert.strictEqual(aliasSymbol!.name, "Foo")
    assert.strictEqual(aliasSymbol!.escapedName, "Foo")

    const {extendsType} = type as ConditionalType
    assert.strictEqual(extendsType.flags, TypeFlags.String)

    const resolvedTrueType = await (type as TypeEx).getResolvedTrueTypeEx()
    assert.strictEqual(resolvedTrueType.flags, TypeFlags.StringLiteral)
    assert.strictEqual((resolvedTrueType as LiteralType).value, "1")

    const resolvedFalseType = await (type as TypeEx).getResolvedFalseTypeEx()
    assert.strictEqual(resolvedFalseType.flags, TypeFlags.NumberLiteral)
    assert.strictEqual((resolvedFalseType as LiteralType).value, 1)
  })

  test("unionType", async ctx => {
    const content = "type Foo = 1 | 2"
    const type = await getElementType(ctx, content, "Foo")
    assert.strictEqual(type.flags, 1<<27 /* Union */)
    const {types} = (type as UnionType)
    assert.strictEqual(types.length, 2)
    const [one, two] = types

    assert.strictEqual(one.flags, TypeFlags.NumberLiteral)
    assert.strictEqual((one as LiteralType).value, 1)

    assert.strictEqual(two.flags, TypeFlags.NumberLiteral)
    assert.strictEqual((two as LiteralType).value, 2)
  })

  test("symbolType", async ctx => {
    const content = "type Foo = {a: 5}"
    const type = await getElementType(ctx, content, "Foo")
    assert.strictEqual(type.flags, 1<<20)

    const [aProperty] = await (type as TypeEx).getPropertiesEx()
    assert.strictEqual(aProperty.name, "a")

    const aType = await (aProperty as SymbolEx).getTypeEx()
    assert.strictEqual(aType.flags, TypeFlags.NumberLiteral)
    assert.strictEqual((aType as LiteralType).value, 5)
  })

  test("cachedTypes", async ctx => {
    const content = "type Foo<T> = T extends 1 ? 1 : 1"
    const type = await getElementType(ctx, content, "Foo")

    const trueType = await (type as TypeEx).getResolvedTrueTypeEx()
    const falseType = await (type as TypeEx).getResolvedFalseTypeEx()

     // === because the link lspApiObjectIdRef existing within the same response
    assert.strictEqual(trueType, falseType)
  })

  test("repeatedRequest", async ctx => {
    const content = "type Foo = {a: 12}"
    const type0 = await getElementType(ctx, content, "Foo")
    const type1 = await getElementType(ctx, content, "Foo")

    // === because caching between requests
    assert.strictEqual(type0, type1)
  })

  test("areMutuallyAssignable", async ctx => {
    const content = "type Foo = {a: 12}\ntype Bar = {a: 12}"
    const typeFoo = await getElementType(ctx, content, "Foo")
    const typeBar = await getElementType(ctx, content, "Bar")
    assert.notStrictEqual(typeFoo, typeBar)
    assert.strictEqual(await (typeFoo as TypeEx).isMutuallyAssignableWith(typeBar), true)
  })

  test("getTypeWithoutOpening", async ctx => {
    // This uses self_managed_projects
    const content = String(await fs.promises.readFile(path.join(testProjectDir, "b.ts")))
    const type = await (ctx as MyContext).client!.getElementType("b.ts", getRange(content, "Foo"))
    assert.strictEqual(type.flags, 1<<20)
    assert.strictEqual(type.aliasSymbol?.escapedName, "Foo")
  })
})

// *** Util ***

function getRange(text: string, element: string) {
  if (element.includes("\n"))
    throw new Error(`Multiline element not implemented: ${element}`)

  const lines = text.split("\n")
  const posIndices = lines
    .map(line => line.indexOf(element))

  const lineIndex = posIndices.findIndex(pos => pos !== -1)
  if (lineIndex === -1)
    throw new Error(`Element ${element} not found in ${text}`)

  return {
    start: { line: lineIndex, character: posIndices[lineIndex] },
    end: { line: lineIndex, character: posIndices[lineIndex] + element.length }
  }
}

function getFragment(content: string, node: Node) {
  const range = (node as NodeEx).range
  const lines = content.split("\n")

  if (range.start.line === range.end.line)
    return lines[range.start.line].slice(range.start.character, range.end.character)

  const startLineFragment = lines[range.start.line].slice(range.start.character)
  const endLineFragment = lines[range.end.line].slice(0, range.end.character)
  const midLines = lines.slice(range.start.line + 1, range.end.line)
  return [startLineFragment, ...midLines, endLineFragment].join("\n")
}

function delay(ms: number) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms)
  })
}

// *** Protocol ***

interface Message {
	jsonrpc: "2.0"
}

interface RequestMessage extends Message {
  id: number | string
  method: string
  params?: unknown[] | object
}

interface ResponseMessage extends Message {
  id: number | string | null
  result?: any
  error?: any
}

interface NotificationMessage extends Message {
  method: string;
  params?: unknown[] | object
}

// *** Types ***

type TypeRequestKind = "Default" | "Contextual" | "ContextualCompletions"

type ServerObjectDef = {lspApiObjectId: number}
type ServerObjectRef = {lspApiObjectIdRef: number}

type ServerType = ServerObjectDef & {
  lspApiObjectType: "TypeObject"
  id: number
  lspApiProjectId: number
  lspApiTypeCheckerId: number
}
type ServerSymbol = ServerObjectDef & {
  lspApiObjectType: "SymbolObject"
  id: number
  lspApiProjectId: number
  lspApiTypeCheckerId: number
}

type ServerSignature = ServerObjectDef & {
  lspApiObjectType: "SignatureObject"
}

type ServerIndexInfo = ServerObjectDef & {
  lspApiObjectType: "IndexInfo"
}

type ServerNode = ServerObjectDef & {
  lspApiObjectType: "NodeObject"
}

type ServerOtherObject = ServerObjectDef & {}

type ServerObject = ServerObjectRef | ServerType | ServerSymbol | ServerIndexInfo | ServerNode | ServerOtherObject

/**
 * These are properties which were formerly received via Checker, Type etc methods and properties.
 * Some of them were internal.
 * 
 * They can be moved now to the `handleCustomTsServerCommand` handler.
 */
type TypeEx = Type & {
  /**
   * Type.id
   */
  id: number
  /**
   * Checker.getResolvedTypeArguments()
   */
  resolvedTypeArguments?: Type[]
  /**
   * Checker.getBaseConstraintsOfType()
   */
  constraint?: Type
  /**
   * This is enumQualifiedName received from Nodes structure
   */
  nameType?: string
  /**
   * Property of TypeParameter
   */
  isThisType?: boolean
  /**
   * Property of IntrinsicType
   */
  intrinsicName?: string

  /**
   * Async versions of Type.get...() methods
   */
  getCallSignaturesEx(): Promise<Signature[]>
  getConstructSignaturesEx(): Promise<Signature[]>
  getPropertiesEx(): Promise<Symbol[]>
  getResolvedFalseTypeEx(): Promise<Type>
  getResolvedTrueTypeEx(): Promise<Type>
 
  /**
   * Results of 'getTypeProperties'
   */
  callSignatures?: Signature[]
  constructSignatures?: Signature[]
  indexInfos?: IndexInfo[]
  properties?: Symbol[]
  resolvedFalseType?: Type
  resolvedProperties?: Symbol[]
  resolvedTrueType?: Type

  /**
   * Types equivalence via t.isAssignableTo(u) && u.isAssignableTo(t)
   */
  isMutuallyAssignableWith(t: Type): Promise<boolean>
}

type SymbolEx = Symbol & {
  /**
   * Checker.getTypeOfSymbolAtLocation()
   */
  getTypeEx(): Promise<Type>
}

type SignatureEx = Signature & {
  /**
   * Internal fields
   */
  flags: number
  /**
   * checker.getReturnTypeOfSignature()
   */
  resolvedReturnType: Type
}

type NodeEx = Node & {
  /**
   * ts.isComputedPropertyName()
   */
  computedProperty: boolean
  /**
   * instead of pos, end
   */
  range: {start: LineAndCharacter, end: LineAndCharacter}
}

type ClientObject = Type | Symbol | Signature | IndexInfo | Node | PseudoBigInt

// *** Client ***

class LspClient {
  #testName: string
  #logFile: string
  #testProjectDir: string

  constructor(testName: string, logFile: string, testProjectDir: string) {
    this.#testName = testName
    this.#logFile = logFile
    this.#testProjectDir = testProjectDir
  }

  async spawnAndInit() {
    await this.#openLog()
    this.#spawnSubprocess()
    this.#addStderrHandler()
    this.#addStdoutHandler()
    await this.#sendRequest({
      jsonrpc: "2.0",
      id: this.#nextId(),
      method: "initialize",
      params: {
        processId: process.pid,
        rootPath: this.#testProjectDir,
        rootUri: url.pathToFileURL(this.#testProjectDir),
        capabilities: {},
        clientInfo: { name: "LSP API Test" },
      },
    })
    this.#sendMessage({jsonrpc: "2.0", method:"initialized", params:{}})
    return this
  }

  #logFileHandle: fs.promises.FileHandle | undefined
  #logFileWritingPromise: Promise<unknown> | undefined
  async #openLog() {
    this.#logFileHandle = await fs.promises.open(this.#logFile, "a")
    this.#log("---", `Starting test: ${this.#testName}`)
  }

  async #log(src: string, msg: string) {
    for ( ; ; ) {
      const promise = this.#logFileWritingPromise
      await promise
      if (promise === this.#logFileWritingPromise) break
    }

    const newline = msg.at(-1) === "\n" ? "" : "\n"
    const writePromise = fs.promises.writeFile(this.#logFileHandle!, Buffer.from(`${src} ${new Date().toISOString()} ${msg}${newline}`))
    this.#logFileWritingPromise = writePromise
    return writePromise
  }

  #childProcess: subprocess.ChildProcess | undefined
  #spawnSubprocess() {
    const tsgoExec = path.join(import.meta.dirname, "..", "..", "..", "built", "local", "tsgo")
    const args = [ "--lsp", "--stdio" ]
    this.#log("---", `Spawning "${tsgoExec}" with args [${args.join(", ")}]`)
    this.#childProcess = subprocess.execFile(tsgoExec, args)
    this.#log("---", `The subprocess started with pid ${this.#childProcess.pid}`)
  }

  #addStderrHandler() {
    this.#childProcess!.stderr!.on("data", data => {
      this.#log("ERR", String(data))
    })
  }

  #inDataPending = ""
  #addStdoutHandler() {
    this.#childProcess!.stdout!.on("data", newData => {
      this.#log("IN ", String(newData))
      this.#inDataPending += newData
      for ( ; ; ) {
        const remainingData = this.#handleInData(this.#inDataPending)
        if (remainingData !== this.#inDataPending)
          this.#inDataPending = remainingData
        else
          break
      }
    })
  }

  #handleInData(inData: string) {
    const sep = "\r\n\r\n"
    const sepPos = inData.indexOf(sep)
    if (sepPos === -1) return inData

    const m = /Content-Length: (\d+)\r\n/.exec(inData)
    if (!m) throw new Error(`Content-Length not found in ${inData}`)

    const inDataPendingBuf = Buffer.from(inData) // non-ascii
    const contentLength = Number(m[1])
    const content = inDataPendingBuf.subarray(sepPos + sep.length, sepPos + sep.length + contentLength)

    this.#handleInputMessage(JSON.parse(String(content)))

    const remainingData = String(inDataPendingBuf.subarray(sepPos + sep.length + contentLength))
    return remainingData
  }

  #pendingIdToResponseConsumer = new Map<string | number, {resolve: (r: ResponseMessage) => void, reject: (e: unknown) => void}>()
  #handleInputMessage(message: RequestMessage | ResponseMessage | NotificationMessage) {
    if ("params" in message && "id" in message) {
      // server -> client request
      this.#sendMessage({jsonrpc: "2.0", id: message.id, result: null})
    } else if ("id" in message && message.id != null) {
      const consumer = this.#pendingIdToResponseConsumer.get(message.id)
      if (consumer) {
        this.#pendingIdToResponseConsumer.delete(message.id)
        if ("error" in message)
          consumer.reject(message)
        else
          consumer.resolve(message)
      }
    } else if ("error" in message) {
      this.#rejectAllPendingRequests(message)
    }
  }

  async #sendRequest(request: RequestMessage) {
    return new Promise((resolve: (r: ResponseMessage) => void, reject: any) => {
      if (this.#pendingIdToResponseConsumer.has(request.id))
        throw new Error(`Duplicate request id ${request.id}`)
      this.#pendingIdToResponseConsumer.set(request.id, {resolve, reject})
      this.#sendMessage(request)
    })
  }

  #sendMessage(request: RequestMessage | ResponseMessage | NotificationMessage) {
    const requestStr = JSON.stringify(request)
    const requestBuf = Buffer.from(requestStr)
    const header = `Content-Length: ${requestBuf.length}\r\n\r\n`
    this.#log("OUT", header)
    this.#log("OUT", requestStr)
    this.#childProcess!.stdin!.write(header)
    this.#childProcess!.stdin!.write(requestBuf)
  }

  #getFileUri(fileRelName: string) {
    return url.pathToFileURL(path.join(this.#testProjectDir, fileRelName))
  }

  #_nextId = 0
  #nextId() {
    return this.#_nextId++
  }

  didOpen(fileRelName: string, text: string) {
    this.#sendMessage({
      jsonrpc: "2.0",
      method:"textDocument/didOpen",
      params: {
        textDocument: {
          uri: this.#getFileUri(fileRelName),
          languageId: "typescript",
          version: 0,
          text,
        }
      }
    })
  }

  async getElementType(fileRelName: string, range: Range, forceReturnType: boolean = false, typeRequestKind: TypeRequestKind = "Default") {
    const file = this.#getFileUri(fileRelName)
    const response = await this.#sendRequest({
      jsonrpc: "2.0",
      id: this.#nextId(),
      method: "$/handleCustomLspApiCommand",
      params: {
        lspApiCommand:"getElementType",
        args: {
          file,
          range,
          forceReturnType,
          typeRequestKind,
        }
      }
    })
    return this.#convertServerType(response.result.response, file) as Type
  }

  async getTypeProperties(
    typeId: number,
    projectId: number,
    typeCheckerId: number,
    originalRequestUri?: URL,
  ) {
    const response = await this.#sendRequest({
      jsonrpc: "2.0",
      id: this.#nextId(),
      method: "$/handleCustomLspApiCommand",
      params: {
        lspApiCommand:"getTypeProperties",
        args: { typeId, projectId, typeCheckerId, originalRequestUri},
      }
    })
    return this.#convertServerType(response.result.response) as Type
  }

  async getSymbolType(
    symbolId: number,
    projectId: number,
    typeCheckerId: number,
  ) {
    const response = await this.#sendRequest({
      jsonrpc: "2.0",
      id: this.#nextId(),
      method: "$/handleCustomLspApiCommand",
      params: {
        lspApiCommand: "getSymbolType",
        args: { symbolId, projectId, typeCheckerId },
      }
    })
    return this.#convertServerType(response.result.response) as Type
  }

  async areTypesMutuallyAssignable(
    projectId: number,
    typeCheckerId: number,
    type1Id: number,
    type2Id: number,
  ) {
    const response = await this.#sendRequest({
      jsonrpc: "2.0",
      id: this.#nextId(),
      method: "$/handleCustomLspApiCommand",
      params: {
        lspApiCommand: "areTypesMutuallyAssignable",
        args: { projectId, typeCheckerId, type1Id, type2Id },
      }
    })
    return response.result.response.areMutuallyAssignable as boolean
  }

  #resolveMapBetweenRequests = new Map<number, ClientObject>()

  #convertServerType(rootServerType: ServerType, fileUri?: URL): ClientObject {
    if (!isType(rootServerType))
      throw new Error(`Root server object must be a type: ${JSON.stringify(rootServerType)}`)

    const ths = this

    const resolveMapWithinSameResponse = new Map<number, {serverObject: ServerObject, clientObject: ClientObject}>()

    function resolveOrConvert<From extends ServerObject, To extends ClientObject>(
      serverObject: From,
      converter: (serverObj: From, clientObj: To) => To,
    ): To {
      const {lspApiObjectIdRef} = serverObject as ServerObjectRef
      if (lspApiObjectIdRef != null) {
        const objs = resolveMapWithinSameResponse.get(lspApiObjectIdRef)
        if (!objs)
          throw new Error(`Could not resolve reference ${lspApiObjectIdRef} in ${JSON.stringify(serverObject)}`)
        return objs.clientObject as To
      }

      const {lspApiObjectId} = serverObject as ServerObjectDef
      if (lspApiObjectId != null) {
        if (resolveMapWithinSameResponse.has(lspApiObjectId))
          throw new Error(`Duplicate lspApiObjectId ${lspApiObjectId} in ${JSON.stringify(serverObject)}`)

        const clientObject = {} as To
        resolveMapWithinSameResponse.set(lspApiObjectId, {serverObject, clientObject})
        if (isType(serverObject))
          ths.#resolveMapBetweenRequests.set(serverObject.id, clientObject)
        converter(serverObject, clientObject)
        return clientObject
      }

      if (isType(serverObject)) {
        const convertedEarlier = ths.#resolveMapBetweenRequests.get(serverObject.id)
        if (convertedEarlier)
          return convertedEarlier as To
      }

      throw new Error(`Could not convert or resolve ${JSON.stringify(serverObject)}`)
    }

    function convertType(typeServerObj: ServerType, target: Type): Type {

      // Properties must be processed in order because definitions come before refs
      for (const [key, value] of Object.entries(typeServerObj)) {

        // Properties returned by both "getElementType" and "getPropertiesOfType"
        if (key === "flags" && typeof value === "string")
          target.flags = Number(value)

        else if (key === "id" && typeof value === "number")
          (target as TypeEx).id = value

        else if (key === "objectFlags" && typeof value === "string")
          (target as ObjectType).objectFlags = Number(value)

        // Returned by "getElementType" (alphabetically)
        else if (key === "aliasSymbol" && isSymbol(value))
          target.aliasSymbol = resolveOrConvert(value, convertSymbol)

        else if (key === "aliasTypeArguments" && isTypes(value))
          target.aliasTypeArguments = value
            .map(arg => resolveOrConvert(arg, convertType))

        else if (key === "baseType" && isType(value))
          (target as SubstitutionType).baseType = resolveOrConvert(value, convertType)

        else if (key === "checkType" && isType(value))
          (target as ConditionalType).checkType = resolveOrConvert(value, convertType)

        else if (key === "constraint" && isType(value))
          (target as TypeEx).constraint = resolveOrConvert(value, convertType)

        else if (key === "elementFlags" && isNumbers(value))
          (target as TupleType).elementFlags = value

        else if (key === "extendsType" && isType(value))
          (target as ConditionalType).extendsType = resolveOrConvert(value, convertType)

        else if (key === "freshType" && isType(value))
          (target as LiteralType).freshType = resolveOrConvert(value, convertType) as LiteralType

        else if (key === "indexType" && isType(value))
          (target as IndexedAccessType).indexType = resolveOrConvert(value, convertType)

        else if (key === "intrinsicName" && typeof value === "string")
          (target as TypeEx).intrinsicName = value

        else if (key === "isThisType" && typeof value === "boolean")
          (target as TypeEx).isThisType = value

        else if (key === "nameType" && typeof value === "string")
          (target as TypeEx).nameType = value

        else if (key === "objectType" && isType(value))
          (target as IndexedAccessType).objectType = resolveOrConvert(value, convertType)

        else if (key === "resolvedTypeArguments" && isTypes(value))
          (target as TypeEx).resolvedTypeArguments = value
            .map(t => resolveOrConvert(t, convertType))

        else if (key === "symbol" && isSymbol(value))
          target.symbol = resolveOrConvert(value, convertSymbol)

        else if (key === "target" && isType(value))
          (target as TypeReference).target = resolveOrConvert(value, convertType) as GenericType

        else if (key === "texts" && isStrings(value))
          (target as TemplateLiteralType).texts = value

        else if (key === "type" && isType(value))
          (target as IndexType).type = resolveOrConvert(value, convertType)

        else if (key === "types"&& isTypes(value))
          (target as UnionType).types = value
            .map(t => resolveOrConvert(t, convertType))

        else if (key === "value") {
          if (isOtherObject(value))
            (target as BigIntLiteralType).value = resolveOrConvert(value, convertPseudoBigInt)
          if (typeof value === "string" || typeof value === "number")
            (target as LiteralType).value = value
        }

        // Returned by "getPropertiesOfType" (alphabetically)

        else if (key === "callSignatures" && isSignatures(value))
          (target as TypeEx).callSignatures = value
            .map(sig => resolveOrConvert(sig, convertSignature))

        else if (key === "constructSignatures" && isSignatures(value))
          (target as TypeEx).constructSignatures = value
            .map(sig => resolveOrConvert(sig, convertSignature))

        else if (key === "indexInfos" && isIndexInfos(value))
          (target as TypeEx).indexInfos = value
            .map(info => resolveOrConvert(info, convertIndexInfo))

        else if (key === "properties" && isSymbols(value))
          (target as TypeEx).properties = value
            .map(sym => resolveOrConvert(sym, convertSymbol))

        else if (key === "resolvedFalseType" && isType(value))
          (target as TypeEx).resolvedFalseType = resolveOrConvert(value, convertType)

        else if (key === "resolvedProperties" && isSymbols(value))
          (target as TypeEx).resolvedProperties = value
            .map(sym => resolveOrConvert(sym, convertSymbol))

        else if (key === "resolvedTrueType" && isType(value))
          (target as TypeEx).resolvedTrueType = resolveOrConvert(value, convertType)
      }

      // Adding methods

      target.getFlags = () => target.flags
      target.getSymbol = () => target.symbol

      let typePropertiesResponse: Type | undefined
      async function getTypeProperties() {
        if (!typePropertiesResponse)
          typePropertiesResponse = await ths.getTypeProperties(
            typeServerObj.id,
            typeServerObj.lspApiProjectId,
            typeServerObj.lspApiTypeCheckerId,
            fileUri,
          )
        return typePropertiesResponse as TypeEx
      }

      (target as TypeEx).getCallSignaturesEx = async () => (await getTypeProperties()).callSignatures!;
      (target as TypeEx).getConstructSignaturesEx = async () => (await getTypeProperties()).constructSignatures!;
      (target as TypeEx).getPropertiesEx = async () => (await getTypeProperties()).properties!;
      (target as TypeEx).getResolvedFalseTypeEx = async () => (await getTypeProperties()).resolvedFalseType!;
      (target as TypeEx).getResolvedTrueTypeEx = async () => (await getTypeProperties()).resolvedTrueType!;

      (target as TypeEx).isMutuallyAssignableWith = async t =>
        await ths.areTypesMutuallyAssignable(
          rootServerType.lspApiProjectId,
          rootServerType.lspApiTypeCheckerId,
          typeServerObj.id,
          (t as TypeEx).id,
        )

      return target
    }

    function convertSymbol(symbolServerObj: ServerSymbol, target: Symbol): Symbol {
      // Processing in order

      for (const [key, value] of Object.entries(symbolServerObj)) {
        if (key === "declarations" && isNodes(value))
          target.declarations = value
            .map(node => resolveOrConvert(node, convertNode) as Declaration)

            else if (key === "escapedName" && typeof value === "string") {
          target.escapedName = value as __String
          (target as {name: string}).name = value
        }

        else if (key === "flags" && typeof value === "string")  
          target.flags = Number(value)

        else if (key === "valueDeclaration" && isNode(value))
          target.valueDeclaration = resolveOrConvert(value, convertNode) as Declaration
      }


      // Method

      let type: Type | undefined
      async function getTypeImpl() {
        type ??= await ths.getSymbolType(symbolServerObj.id, rootServerType.lspApiProjectId, rootServerType.lspApiTypeCheckerId)
        return type
      }

      (target as SymbolEx).getTypeEx = async () => await getTypeImpl()

      return target
    }

    function convertSignature(signatureServerObj: ServerSignature, target: Signature): Signature {
      // Processing in order

      for (const [key, value] of Object.entries(signatureServerObj)) {
        if (key === "declaration" && isNode(value))
          target.declaration = resolveOrConvert(value, convertNode) as SignatureDeclaration

        else if (key === "flags" && typeof value === "string")
          (target as SignatureEx).flags = Number(value)

        else if (key === "parameters" && isSymbols(value))
          target.parameters = value.map(sym => resolveOrConvert(sym, convertSymbol))
    
        else if (key === "resolvedReturnType" && isType(value))
          (target as SignatureEx).resolvedReturnType = resolveOrConvert(value, convertType)

        else if (key === "typeParameters" && isTypes(value))
          target.typeParameters = value.map(type => resolveOrConvert(type, convertType))
      }

      return target
    }

    function convertNode(nodeServerObj: ServerNode, target: Node): Node {
      for (const [key, value] of Object.entries(nodeServerObj)) {
        if (key === "computedProperty" && typeof value === "boolean")
          (target as NodeEx).computedProperty = value

        else if (key === "fileName" && typeof value === "string")
          (target as SourceFile).fileName = value

        else if (key === "parent" && isNode(value))
          (target as {parent: Node}).parent = resolveOrConvert(value, convertNode)

        else if (key === "range" && typeof value === "object")
          (target as NodeEx).range = value as NodeEx["range"]
      }

      return target
    }

    function convertPseudoBigInt(pseudoBigIntServerObj: ServerOtherObject, target: PseudoBigInt): PseudoBigInt {
      for (const [key, value] of Object.entries(pseudoBigIntServerObj)) {
        if (key === "negative" && typeof value === "boolean")
          target.negative = value

        else if (key === "base10Value" && typeof value === "string")
          target.base10Value = value
      }

      return target
    }

    function convertIndexInfo(indexInfoServerObj: ServerIndexInfo, target: IndexInfo): IndexInfo {
      for (const [key, value] of Object.entries(indexInfoServerObj)) {
        if (key === "declaration" && isNode(value))
          target.declaration = resolveOrConvert(value, convertNode) as IndexSignatureDeclaration

        else if (key === "isReadonly" && typeof value === "boolean")
          target.isReadonly = value

        else if (key === "keyType" && isType(value))
          target.keyType = resolveOrConvert(value, convertType)

        else if (key === "type" && isType(value))
          target.type = resolveOrConvert(value, convertType)
      }

      return target
    }


    function isType(serverObject: unknown): serverObject is ServerType {
      const isTypeImpl = (o: unknown): o is ServerType => (o as ServerType).lspApiObjectType === "TypeObject"
      return isTypeImpl(serverObject) || isRefTo(serverObject, isTypeImpl)
    }

    function isRefTo<T>(serverObject: unknown, defChecker: (obj: unknown) => obj is T): serverObject is T {
      const {lspApiObjectIdRef} = serverObject as ServerObjectRef
      if (lspApiObjectIdRef != null) {
        const {serverObject} = resolveMapWithinSameResponse.get(lspApiObjectIdRef) ?? {}
        return !!serverObject && defChecker(serverObject)
      }
      return false
    }

    function isTypes(serverObject: unknown): serverObject is ServerType[] {
      return Array.isArray(serverObject) && serverObject.every(isType)
    }

    function isSymbol(serverObject: unknown): serverObject is ServerSymbol {
      const isSymbolImpl = (o: unknown): o is ServerSymbol => (o as ServerSymbol).lspApiObjectType === "SymbolObject"
      return isSymbolImpl(serverObject) || isRefTo(serverObject, isSymbolImpl)
    }

    function isSymbols(serverObject: unknown): serverObject is ServerSymbol[] {
      return Array.isArray(serverObject) && serverObject.every(isSymbol)
    }

    function isSignature(serverObject: unknown): serverObject is ServerSignature {
      const isSignatureImpl = (o: unknown): o is ServerSignature => (o as ServerSignature).lspApiObjectType === "SignatureObject"
      return isSignatureImpl(serverObject) || isRefTo(serverObject, isSignatureImpl)
    }

    function isSignatures(serverObject: unknown): serverObject is ServerSignature[] {
      return Array.isArray(serverObject) && serverObject.every(isSignature)
    }

    function isNode(serverObject: unknown): serverObject is ServerNode {
      const isNodeImpl = (o: unknown): o is ServerNode => (o as ServerNode).lspApiObjectType === "NodeObject"
      return isNodeImpl(serverObject) || isRefTo(serverObject, isNodeImpl)
    }

    function isNodes(serverObject: unknown): serverObject is ServerNode[] {
      return Array.isArray(serverObject) && serverObject.every(isNode)
    }

    function isIndexInfo(serverObject: unknown): serverObject is ServerIndexInfo {
      const isIndexInfoImpl = (o: unknown): o is ServerIndexInfo => (o as ServerIndexInfo).lspApiObjectType === "IndexInfo"
      return isIndexInfoImpl(serverObject) || isRefTo(serverObject, isIndexInfoImpl)
    }

    function isIndexInfos(serverObject: unknown): serverObject is ServerIndexInfo[] {
      return Array.isArray(serverObject) && serverObject.every(isIndexInfo)
    }

    function isOtherObject(serverObject: unknown): serverObject is ServerOtherObject {
      const isOtherObjectImpl = (o: unknown): o is ServerOtherObject =>
        serverObject != null
          && (serverObject as ServerObjectDef).lspApiObjectId != null
          && typeof serverObject === "object"
          && !("lspApiObjectType" in serverObject)

      return isOtherObjectImpl(serverObject) || isRefTo(serverObject, isOtherObjectImpl)
    }

    function isStrings(serverObject: unknown): serverObject is string[] {
      return Array.isArray(serverObject) && serverObject.every(s => typeof s === "string")
    }

    function isNumbers(serverObject: unknown): serverObject is number[] {
      return Array.isArray(serverObject) && serverObject.every(n => typeof n === "number")
    }


    return resolveOrConvert(rootServerType, convertType)
  }

  async shutdown() {
    await this.#sendRequest({jsonrpc: "2.0", id: this.#nextId(), method: "shutdown"})
    this.#sendMessage({jsonrpc: "2.0", method: "exit"})
    while (this.#childProcess!.exitCode == null) await delay(50)
    await this.#closeLog()
    this.#rejectAllPendingRequests(new Error("The subprocess exited"))
  }

  async kill() {
    this.#childProcess!.kill()
    await this.#closeLog()
    this.#rejectAllPendingRequests(new Error("The subprocess was killed"))
  }

  async #closeLog() {
    await this.#log("---", `The subprocess exited with the code: ${this.#childProcess!.exitCode}`)
    await this.#log("---", `Finished test: ${this.#testName}`)
    await this.#logFileHandle!.close();
  }

  #rejectAllPendingRequests(reason: any) {
    for (const {reject} of this.#pendingIdToResponseConsumer.values())
      reject(reason)
    this.#pendingIdToResponseConsumer.clear()
  }
}
