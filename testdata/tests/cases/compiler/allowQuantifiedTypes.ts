type T0 = <T extends object> { values: T[], identifier: (value: T) => string }
type T1 = <T extends string>(t: T) => T
type T2 = <T extends string> <U extends T>(u: U) => U 