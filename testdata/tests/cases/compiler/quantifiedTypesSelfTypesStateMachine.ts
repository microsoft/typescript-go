type ParseStateMachine<Self> = { [S in keyof Self]: { [E in keyof Self[S]]: keyof Self } }
type StateMachine = <Self extends ParseStateMachine<Self>> Self

let trafficLights: StateMachine = {
  off: {
    ON: "red"
  },
  red: {
    TICK: "yellow",
    OFF: "off"
  },
  yellow: {
    TICK: "green",
    OFF: "off"
  },
  green: {
    TICK: "red",
    OFF: "off"
  }
}

let trafficLightsInvalid: StateMachine = {
  off: {
    ON: "red"
  },
  red: {
    TICK: "yellow",
    OFF: "off"
  },
  yellow: {
    TICK: "green",
    OFF: "off"
  },
  green: {
    TICK: "reddd",
    OFF: "off"
  }
}