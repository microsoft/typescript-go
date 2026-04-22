// @strict: true
// @target: es2022
// @module: esnext
// @moduleResolution: bundler
// @noEmit: true

interface TypeData {
  ref?: string;
}

interface UsageData {
  Type_Data?: Array<TypeData>;
}

interface EmailAddressData {
  Email_Address: string;
  Usage_Data?: Array<UsageData>;
}

declare const cond: boolean;
declare const emails: { value: string }[];
declare function makeRef(s: string): string;
declare function take(xs: EmailAddressData[]): EmailAddressData[];

function f(): EmailAddressData[] {
  let v;
  if (cond) {
    v = [
      {
        Email_Address: emails[0].value,
        Usage_Data: [
          {
            attributes: { p: false },
            Type_Data: [{ ref: makeRef("x") }],
          },
        ],
      },
    ];
  } else {
    v = take([
      {
        Email_Address: emails[0].value,
        Usage_Data: [
          {
            attributes: { p: false },
            Type_Data: [{ ref: makeRef("x") }],
          },
        ],
      },
    ]);
  }
  return v;
}
