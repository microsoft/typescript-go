// @noEmit: true
// @strict: true

declare namespace Router {
    interface Params<T, U> {
        router: Router<U, {}>;
    }

    type Context<T, U> =
        & { state: T }
        & Params<T, U>;

    type Middleware<T, U> = (context: Context<T, U>) => void;
}

declare class Router<T, U> {
    routes(): Router.Middleware<T, U>;
}

declare let a: Router.Params<{ value: string }, {}>;
declare let b: Router.Params<{ value: number }, {}>;

a = b; // Error
