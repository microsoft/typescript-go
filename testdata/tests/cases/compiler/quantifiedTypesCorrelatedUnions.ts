// https://github.com/microsoft/TypeScript/issues/30581

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

function processRecord2(record1: GenericRecord, record2: GenericRecord) {
  record1.f(record2.v); // TODO: better error
}