// @strict: true
// @noEmit: true

class Session {
    user: string | null = null;

    get hasUser(): this is { user: string } {
        return this.user !== null;
    }

    hasUserMethod(): this is { user: string } {
        return this.user !== null;
    }
}

const session = new Session();

if (session.hasUser) {
    session.user.toUpperCase();
}

if (session.hasUserMethod()) {
    session.user.toUpperCase();
}

if (!session.hasUser) {
    session.user; // string | null
}
