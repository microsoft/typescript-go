type PasreCaseInsensitive<Self, T extends string> =
  Self extends string
    ? Lowercase<Self> extends Lowercase<T>
        ? Self
        : `Error: Type '${Self}' is not assignable to type 'CaseInsensitive<${T}>'`
    : T

type CaseInsensitive<T extends string> = <Self> PasreCaseInsensitive<Self, T>

declare const setHeader: 
  (key: CaseInsensitive<"Set-Cookie" | "Accept">, value: string) => void

setHeader("Set-Cookie", "test")
setHeader("Accept", "test2")
setHeader("sEt-cOoKiE", "stop writing headers like this but ok")
setHeader("Acept", "nah this has a typo")