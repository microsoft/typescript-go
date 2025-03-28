//// [tests/cases/compiler/genericIndexedAccessMethodIntersectionCanBeAccessed.ts] ////

//// [genericIndexedAccessMethodIntersectionCanBeAccessed.ts]
type ExtendedService<T> = {
    [K in keyof T]: T[K] & {
        __$daemonMode?: string;
        __$action?: string;
    };
};

type Service<T> = {
    [K in keyof T]: T[K] & {id?: string};
};

export const createService = <T>(
    ServiceCtr: ExtendedService<T> & Service<T>
) => {
    Object.keys(ServiceCtr).forEach(key => {
        const method = (ServiceCtr)[key as keyof T];
        const {__$daemonMode, __$action, id} = method;
    })
}


//// [genericIndexedAccessMethodIntersectionCanBeAccessed.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createService = void 0;
const createService = (ServiceCtr) => {
    Object.keys(ServiceCtr).forEach(key => {
        const method = (ServiceCtr)[key];
        const { __$daemonMode, __$action, id } = method;
    });
};
exports.createService = createService;
