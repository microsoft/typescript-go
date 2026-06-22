// https://github.com/microsoft/typescript-go/issues/4391
// @noEmit: true
// @strict: true

type RedirectOptions<TTo extends string> = {
    from?: string;
    to: TTo;
    params: TTo extends "/agents/$agentName/$sessionId" ? { agentName: string; sessionId: string } : { agentName: string };
};

declare function redirect<const TTo extends string>(opts: RedirectOptions<TTo>): void;

interface RouteRedirect {
    <const TTo extends string>(opts: Omit<RedirectOptions<TTo>, "from">): void;
}

declare const routeRedirect: RouteRedirect;

redirect({
    from: "/agents/$agentName/",
    to: "/agents/$agentName/$sessionId",
    params: { agentName: "a", sessionId: "s" },
});

routeRedirect({
    to: "/agents/$agentName/$sessionId",
    params: { agentName: "a", sessionId: "s" },
});
