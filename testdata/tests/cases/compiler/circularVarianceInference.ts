// @noEmit: true
// @strict: true

declare namespace Router {
    interface Params<T, U> {
        router: Router<U, {}>;
    }

    type Context<State, CustomContext> =
        & { state: State }
        & Params<State, CustomContext>;

    type Middleware<State, CustomContext> = (
        context:
            & { state: State }
            & CustomContext
            & Params<CustomContext, CustomContext>,
    ) => any;
}

declare class Router<State, CustomContext> {
    post<AddedState, AddedContext>(
        path: string,
        ...middleware: Router.Middleware<AddedState, AddedContext>[]
    ): Router.Middleware<State, {}>;
}

interface MyState {
    foo: string;
}

interface MyContext {}

declare const router: Router<MyState, MyContext>;
declare const routeHandler: (
    context: Router.Context<MyState, MyContext>,
) => void;

router.post("/foo", routeHandler);
