import {
    Project,
    Symbol,
    Type,
} from "./api.ts";
import type { Client } from "./client.ts";
import type {
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "./proto.ts";

export class ObjectRegistry {
    private client: Client;
    private projects: Map<string, WeakRef<Project>> = new Map();
    private symbols: Map<string, WeakRef<Symbol>> = new Map();
    private types: Map<string, WeakRef<Type>> = new Map();

    private pendingReleases: Set<string> = new Set();

    private finalizationRegistry = new FinalizationRegistry((handle: string) => {
        this.pendingReleases.add(handle);
        queueMicrotask(() => {
            this.releaseObjects();
        });

        switch (handle[0]) {
            case "p":
                this.projects.delete(handle);
                break;
            case "s":
                this.symbols.delete(handle);
                break;
            case "t":
                this.types.delete(handle);
                break;
            default:
                throw new Error(`Unknown handle type: ${handle[0]}`);
        }
    });

    constructor(client: Client) {
        this.client = client;
    }

    private releaseObjects() {
        this.client.request("release", Array.from(this.pendingReleases));
        this.pendingReleases.clear();
    }

    getProject(data: ProjectResponse): Project {
        const projectRef = this.projects.get(data.id);
        let project = projectRef?.deref();
        if (project) {
            return project;
        }

        project = new Project(this.client, this, data);
        this.projects.set(data.id, new WeakRef(project));
        this.finalizationRegistry.register(project, data.id, project);
        return project;
    }

    getSymbol(data: SymbolResponse): Symbol {
        const symbolRef = this.symbols.get(data.id);
        let symbol = symbolRef?.deref();
        if (symbol) {
            return symbol;
        }

        symbol = new Symbol(this.client, data);
        this.symbols.set(data.id, new WeakRef(symbol));
        this.finalizationRegistry.register(symbol, data.id, symbol);
        return symbol;
    }

    getType(data: TypeResponse): Type {
        const typeRef = this.types.get(data.id);
        let type = typeRef?.deref();
        if (type) {
            return type;
        }

        type = new Type(this.client, data);
        this.types.set(data.id, new WeakRef(type));
        this.finalizationRegistry.register(type, data.id, type);
        return type;
    }

    releaseProject(project: Project): void {
        this.finalizationRegistry.unregister(project);
        this.pendingReleases.delete(project.id);
        this.projects.delete(project.id);
        this.client.request("release", project.id);
    }

    releaseSymbol(symbol: Symbol): void {
        this.finalizationRegistry.unregister(symbol);
        this.pendingReleases.delete(symbol.id);
        this.symbols.delete(symbol.id);
        this.client.request("release", symbol.id);
    }

    releaseType(type: Type): void {
        this.finalizationRegistry.unregister(type);
        this.pendingReleases.delete(type.id);
        this.types.delete(type.id);
        this.client.request("release", type.id);
    }
}
