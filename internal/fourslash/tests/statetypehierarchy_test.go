package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTypeHierarchyBasic(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /tsconfig.json
{}
// @Filename: /a.ts
interface Base {
    method(): void;
}

interface Middle extends Base {
    anotherMethod(): void;
}

class /*marker*/Derived implements Middle {
    method(): void {}
    anotherMethod(): void {}
}

class SubDerived extends Derived {
    additionalMethod(): void {}
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyClassInheritance(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
class Animal {
    name: string = "";
}

class /*marker*/Mammal extends Animal {
    hasFur: boolean = true;
}

class Dog extends Mammal {
    breed: string = "";
}

class Cat extends Mammal {
    indoor: boolean = true;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyInterface(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface /*marker*/IBase {
    foo(): void;
}

interface IExtended extends IBase {
    bar(): void;
}

class Implementation implements IBase {
    foo(): void {}
}

class ExtendedImpl implements IExtended {
    foo(): void {}
    bar(): void {}
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyTypeAlias(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Nameable {
    name: string;
}

interface Ageable {
    age: number;
}

type /*marker*/Person = Nameable & Ageable;

interface Employee extends Person {
    employeeId: string;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyMultiFile(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /base.ts
export interface IBase {
    baseProp: string;
}
// @Filename: /derived.ts
import { IBase } from './base';

export class /*marker*/Derived implements IBase {
    baseProp: string = "";
    derivedProp: number = 0;
}
// @Filename: /subderived.ts
import { Derived } from './derived';

export class SubDerived extends Derived {
    subProp: boolean = false;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}
func TestTypeHierarchyAbstract(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
abstract class /*marker*/Shape {
    abstract area(): number;
    abstract perimeter(): number;
}

abstract class NamedShape extends Shape {
    constructor(public name: string) {
        super();
    }
    describe(): string {
        return "shape";
    }
}

class Rectangle extends NamedShape {
    constructor(name: string, public width: number, public height: number) {
        super(name);
    }
    area(): number {
        return this.width * this.height;
    }
    perimeter(): number {
        return 2 * (this.width + this.height);
    }
}

class Circle extends NamedShape {
    constructor(name: string, public radius: number) {
        super(name);
    }
    area(): number {
        return Math.PI * this.radius ** 2;
    }
    perimeter(): number {
        return 2 * Math.PI * this.radius;
    }
}

class Square extends Rectangle {
    constructor(name: string, side: number) {
        super(name, side, side);
    }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyGenerics(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Repository<T> {
    find(id: string): T | null;
    save(entity: T): void;
    delete(id: string): void;
}

interface CacheableRepository<T> extends Repository<T> {
    clearCache(): void;
}

class /*marker*/BaseRepository<T> implements Repository<T> {
    find(id: string): T | null { return null; }
    save(entity: T): void {}
    delete(id: string): void {}
}

class CachedRepository<T> extends BaseRepository<T> implements CacheableRepository<T> {
    private cache: Map<string, T> = new Map();
    clearCache(): void { this.cache.clear(); }
}

interface Entity {
    id: string;
}

class User implements Entity {
    id: string = "";
    name: string = "";
}

class UserRepository extends CachedRepository<User> {
    findByName(name: string): User | null { return null; }
}

class AdminRepository extends UserRepository {
    findAdmins(): User[] { return []; }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyComplex(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Serializable {
    serialize(): string;
}

interface Comparable<T> {
    compareTo(other: T): number;
}

interface Named {
    name: string;
}

interface Timestamped {
    createdAt: Date;
    updatedAt: Date;
}

interface Entity extends Named, Timestamped {
    id: string;
}

abstract class /*marker*/BaseModel implements Serializable {
    abstract serialize(): string;
}

class User extends BaseModel implements Entity, Comparable<User> {
    id: string = "";
    name: string = "";
    createdAt: Date = new Date();
    updatedAt: Date = new Date();
    
    serialize(): string {
        return JSON.stringify(this);
    }
    
    compareTo(other: User): number {
        return this.name.localeCompare(other.name);
    }
}

class AdminUser extends User {
    permissions: string[] = [];
}

class SuperAdmin extends AdminUser {
    canManageAdmins: boolean = true;
}

class GuestUser extends User {
    sessionId: string = "";
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyIntersection(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface HasName {
    name: string;
}

interface HasAge {
    age: number;
}

interface HasEmail {
    email: string;
}

type /*marker*/Person = HasName & HasAge;

type Employee = Person & HasEmail & {
    employeeId: string;
};

interface Manager extends Employee {
    department: string;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyDeepInheritance(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
class /*marker*/Level0 {
    prop0: string = "";
}

class Level1 extends Level0 {
    prop1: string = "";
}

class Level2 extends Level1 {
    prop2: string = "";
}

class Level3 extends Level2 {
    prop3: string = "";
}

class Level4 extends Level3 {
    prop4: string = "";
}

class Level5 extends Level4 {
    prop5: string = "";
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyConditionalTypes(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Animal {
    name: string;
}

interface Dog extends Animal {
    bark(): void;
}

interface Cat extends Animal {
    meow(): void;
}

// Basic conditional type
type /*marker*/IsDog<T> = T extends Dog ? true : false;

// Conditional type with different results
type AnimalSound<T> = T extends Dog ? "bark" : T extends Cat ? "meow" : "unknown";

// Extract utility type pattern
type ExtractDog<T> = T extends Dog ? T : never;

// Inferring in conditional types
type ReturnTypeOf<T> = T extends (...args: any[]) => infer R ? R : never;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyMappedTypes(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface User {
    name: string;
    email: string;
    age: number;
}

// Mapped type
type /*marker*/ReadonlyUser = {
    readonly [K in keyof User]: User[K];
};

// Partial-like
type PartialUser = {
    [K in keyof User]?: User[K];
};

// Pick-like
type UserNameAndEmail = {
    [K in "name" | "email"]: User[K];
};`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyClassExpressions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Named class expression assigned to variable
const /*marker*/MyClass = class NamedClass {
    prop: string = "";
};

// Anonymous class expression
const AnotherClass = class {
    prop: number = 0;
};

// Class expression extending another class
class BaseClass {
    baseProp: boolean = true;
}

const DerivedClass = class extends BaseClass {
    derivedProp: string = "";
};`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyNegativeCases(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Simple variable - not a type hierarchy item
const /*marker*/simpleVar = 42;

// Function - not a type hierarchy item  
function simpleFunction() {
    return "hello";
}

// Primitive type alias - should work but have no supertypes
type StringAlias = string;

// Enum - not typically part of type hierarchy
enum Color {
    Red,
    Green,
    Blue
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyEdgeCases(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Self-referencing interface (use unique name to avoid DOM Node conflict)
interface /*marker*/ASTNode {
    value: string;
    children: ASTNode[];
}

// Mutually referencing types
interface TreeASTNode extends ASTNode {
    parent: TreeASTNode | null;
}

class ConcreteASTNode implements ASTNode {
    value: string = "";
    children: ASTNode[] = [];
}

// Multiple levels of self-reference
class TreeASTNodeImpl implements TreeASTNode {
    value: string = "";
    children: ASTNode[] = [];
    parent: TreeASTNode | null = null;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyDuplicates(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Base {
    baseProp: string;
}

interface MiddleA extends Base {
    propA: number;
}

interface MiddleB extends Base {
    propB: boolean;
}

// Diamond inheritance - should deduplicate Base in supertypes
class /*marker*/Diamond implements MiddleA, MiddleB {
    baseProp: string = "";
    propA: number = 0;
    propB: boolean = false;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyIndexedAccess(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
interface Person {
    name: string;
    age: number;
    address: {
        street: string;
        city: string;
    };
}

// Indexed access types
type /*marker*/PersonName = Person["name"];

type PersonAddress = Person["address"];

type PersonKeys = keyof Person;

// Nested indexed access
type PersonCity = Person["address"]["city"];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyTemplateLiterals(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
type Greeting = "hello" | "hi" | "hey";
type Target = "world" | "there";

// Template literal type
type /*marker*/Message = ` + "`${Greeting} ${Target}`" + `;

// Intrinsic string manipulation
type UpperGreeting = Uppercase<Greeting>;
type LowerGreeting = Lowercase<Greeting>;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyEnums(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Basic numeric enum
enum /*marker*/Direction {
    North,
    South,
    East,
    West,
}

// String enum
enum Color {
    Red = 'RED',
    Green = 'GREEN',
    Blue = 'BLUE',
}

// Const enum
const enum LogLevel {
    Debug = 0,
    Info = 1,
    Warning = 2,
    Error = 3,
}

// Enum member types
type DirectionNorth = Direction.North;

// Interface using enum
interface Compass {
    current: Direction;
    history: Direction[];
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyMixins(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Mixin pattern
type Constructor<T = {}> = new (...args: any[]) => T;

function Timestamped<TBase extends Constructor>(Base: TBase) {
    return class extends Base {
        timestamp = Date.now();
    };
}

function Activatable<TBase extends Constructor>(Base: TBase) {
    return class extends Base {
        isActivated = false;
        activate() { this.isActivated = true; }
        deactivate() { this.isActivated = false; }
    };
}

class /*marker*/User {
    name: string = "";
}

// Class using mixins
const TimestampedUser = Timestamped(User);
const ActivatableTimestampedUser = Activatable(Timestamped(User));

class SpecialUser extends ActivatableTimestampedUser {
    special: boolean = true;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyDeclarationMerging(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @Filename: /a.ts
// Interface declaration merging
interface /*marker*/Mergeable {
    firstMethod(): void;
}

interface Mergeable {
    secondMethod(): void;
}

// Class implementing merged interface
class MergedImpl implements Mergeable {
    firstMethod(): void {}
    secondMethod(): void {}
}

// Namespace with interface
namespace MyNamespace {
    export interface NamespaceInterface {
        nsMethod(): void;
    }
    
    export class NamespaceClass implements NamespaceInterface {
        nsMethod(): void {}
    }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}

func TestTypeHierarchyLibExtensions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @stateBaseline: true
// @lib: esnext
// @Filename: /a.ts
// Error hierarchy - extending built-in Error
class /*marker*/ApplicationError extends Error {
    constructor(message: string, public code: number) {
        super(message);
        this.name = 'ApplicationError';
    }
}

class ValidationError extends ApplicationError {
    constructor(message: string, public field: string) {
        super(message, 400);
    }
}

class NetworkError extends ApplicationError {
    constructor(message: string, public statusCode: number) {
        super(message, statusCode);
    }
}

class NotFoundError extends NetworkError {
    constructor(resource: string) {
        super(resource + " not found", 404);
    }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyBaselineTypeHierarchy(t)
}
