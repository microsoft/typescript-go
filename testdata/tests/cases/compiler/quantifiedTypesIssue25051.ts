// https://github.com/microsoft/TypeScript/issues/25051

type NumberRecord = { kind: "n", v: number, f: (v: number) => void };
type StringRecord = { kind: "s", v: string, f: (v: string) => void };
type BooleanRecord = { kind: "b", v: boolean, f: (v: boolean) => void };
type GenericRecord = <T> { kind: string, v: T, f: (v: T) => void };

function processRecord(record: GenericRecord) {
  record.f(record.v);
}
processRecord({} as NumberRecord)
processRecord({} as StringRecord)
processRecord({} as BooleanRecord)
processRecord({} as NumberRecord | StringRecord | BooleanRecord)