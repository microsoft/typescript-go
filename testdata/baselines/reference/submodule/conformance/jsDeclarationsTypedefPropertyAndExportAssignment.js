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

//// [index.js]
const { taskGroups, taskNameToGroup } = require('./module.js');
class MainThreadTasks {
    constructor(x, y) { }
}
module.exports = MainThreadTasks;
