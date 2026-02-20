//// [tests/cases/compiler/distributeControlFlowCorrelatedUnion.ts] ////

//// [distributeControlFlowCorrelatedUnion.ts]
type A = { type: "a", value: boolean }
type B = { type: "b", value: string }

const handlers = {
  a: (action: A) => {},
  b: (action: B) => {},
}

const handle = (action: A | B) => {
  handlers[action.type](action)

  distribute (action) {
    handlers[action.type](action)
  }

  handlers[action.type]({ type: "a", value: true })

  distribute (action) {
    handlers[action.type]({ type: "a", value: true })
  }
}






//// [distributeControlFlowCorrelatedUnion.js]
const handlers = {
    a: (action) => { },
    b: (action) => { },
};
const handle = (action) => {
    handlers[action.type](action);
    {
        handlers[action.type](action);
    }
    handlers[action.type]({ type: "a", value: true });
    {
        handlers[action.type]({ type: "a", value: true });
    }
};
