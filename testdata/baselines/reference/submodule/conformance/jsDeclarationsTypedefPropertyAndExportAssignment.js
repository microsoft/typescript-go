//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypedefPropertyAndExportAssignment.ts] ////

//// [module.js]
/** @typedef {'parseHTML'|'styleLayout'} TaskGroupIds */

/**
 * @typedef TaskGroup
 * @property {TaskGroupIds} id
 * @property {string} label
 * @property {string[]} traceEventNames
 */

/**
 * @type {{[P in TaskGroupIds]: {id: P, label: string}}}
 */
const taskGroups = {
    parseHTML: {
        id: 'parseHTML',
        label: 'Parse HTML & CSS'
    },
    styleLayout: {
        id: 'styleLayout',
        label: 'Style & Layout'
    },
}

/** @type {Object<string, TaskGroup>} */
const taskNameToGroup = {};

module.exports = {
    taskGroups,
    taskNameToGroup,
};
//// [index.js]
const {taskGroups, taskNameToGroup} = require('./module.js');

/** @typedef {import('./module.js').TaskGroup} TaskGroup */

/**
 * @typedef TaskNode
 * @prop {TaskNode[]} children
 * @prop {TaskNode|undefined} parent
 * @prop {TaskGroup} group
 */

/** @typedef {{timers: Map<string, TaskNode>}} PriorTaskData */
class MainThreadTasks {
    /**
     * @param {TaskGroup} x
     * @param {TaskNode} y
     */
    constructor(x, y){}
}

module.exports = MainThreadTasks;

//// [module.js]
"use strict";
/** @typedef {'parseHTML'|'styleLayout'} TaskGroupIds */
/**
 * @typedef TaskGroup
 * @property {TaskGroupIds} id
 * @property {string} label
 * @property {string[]} traceEventNames
 */
/**
 * @type {{[P in TaskGroupIds]: {id: P, label: string}}}
 */
const taskGroups = {
    parseHTML: {
        id: 'parseHTML',
        label: 'Parse HTML & CSS'
    },
    styleLayout: {
        id: 'styleLayout',
        label: 'Style & Layout'
    },
};
/** @type {Object<string, TaskGroup>} */
const taskNameToGroup = {};
module.exports = {
    taskGroups,
    taskNameToGroup,
};
//// [index.js]
"use strict";
const { taskGroups, taskNameToGroup } = require('./module.js');
/** @typedef {import('./module.js').TaskGroup} TaskGroup */
/**
 * @typedef TaskNode
 * @prop {TaskNode[]} children
 * @prop {TaskNode|undefined} parent
 * @prop {TaskGroup} group
 */
/** @typedef {{timers: Map<string, TaskNode>}} PriorTaskData */
class MainThreadTasks {
    /**
     * @param {TaskGroup} x
     * @param {TaskNode} y
     */
    constructor(x, y) { }
}
module.exports = MainThreadTasks;


//// [module.d.ts]
/** @typedef {'parseHTML'|'styleLayout'} TaskGroupIds */
export type TaskGroupIds = 'parseHTML' | 'styleLayout';
export type TaskGroup = {
    id: TaskGroupIds;
    label: string;
    traceEventNames: string[];
};
const _default: {
    taskGroups: {
        parseHTML: {
            id: "parseHTML";
            label: string;
        };
        styleLayout: {
            id: "styleLayout";
            label: string;
        };
    };
    taskNameToGroup: Record<string, TaskGroup>;
};
export = _default;
//// [index.d.ts]
export type TaskGroup = import('./module.js').TaskGroup;
export type TaskNode = {
    children: TaskNode[];
    parent: TaskNode | undefined;
    group: TaskGroup;
};
export type PriorTaskData = {
    timers: Map<string, TaskNode>;
};
/** @typedef {import('./module.js').TaskGroup} TaskGroup */
/**
 * @typedef TaskNode
 * @prop {TaskNode[]} children
 * @prop {TaskNode|undefined} parent
 * @prop {TaskGroup} group
 */
/** @typedef {{timers: Map<string, TaskNode>}} PriorTaskData */
class MainThreadTasks {
    /**
     * @param {TaskGroup} x
     * @param {TaskNode} y
     */
    constructor(x: TaskGroup, y: TaskNode);
}
export = MainThreadTasks;


//// [DtsFileErrors]


out/index.d.ts(18,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/module.d.ts(8,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    export type TaskGroup = import('./module.js').TaskGroup;
    export type TaskNode = {
        children: TaskNode[];
        parent: TaskNode | undefined;
        group: TaskGroup;
    };
    export type PriorTaskData = {
        timers: Map<string, TaskNode>;
    };
    /** @typedef {import('./module.js').TaskGroup} TaskGroup */
    /**
     * @typedef TaskNode
     * @prop {TaskNode[]} children
     * @prop {TaskNode|undefined} parent
     * @prop {TaskGroup} group
     */
    /** @typedef {{timers: Map<string, TaskNode>}} PriorTaskData */
    class MainThreadTasks {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /**
         * @param {TaskGroup} x
         * @param {TaskNode} y
         */
        constructor(x: TaskGroup, y: TaskNode);
    }
    export = MainThreadTasks;
    
==== out/module.d.ts (1 errors) ====
    /** @typedef {'parseHTML'|'styleLayout'} TaskGroupIds */
    export type TaskGroupIds = 'parseHTML' | 'styleLayout';
    export type TaskGroup = {
        id: TaskGroupIds;
        label: string;
        traceEventNames: string[];
    };
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        taskGroups: {
            parseHTML: {
                id: "parseHTML";
                label: string;
            };
            styleLayout: {
                id: "styleLayout";
                label: string;
            };
        };
        taskNameToGroup: Record<string, TaskGroup>;
    };
    export = _default;
    