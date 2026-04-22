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
  Comments?: string;
}

interface EmailAddressData {
  Email_Address: string;
  Usage_Data?: Array<UsageData>;
}

declare function condition(): boolean;
declare function makeRef(s: string): string;
declare const emails: { value: string; isWork: boolean }[];

async function build(): Promise<EmailAddressData[]> {
  let emailAddressData;
  if (condition()) {
    const primary = emails[0];
    emailAddressData = primary
      ? [
          {
            Email_Address: primary.value,
            Usage_Data: [
              {
                attributes: { "wd:Public": false },
                Type_Data: [
                  {
                    attributes: { "wd:Primary": true },
                    ref: makeRef("HOME"),
                  },
                ],
              },
            ],
          },
        ]
      : [];
  } else {
    emailAddressData = emails.map((email) => ({
      Email_Address: email.value,
      Usage_Data: [
        {
          attributes: { "wd:Public": email.isWork },
          Type_Data: [
            {
              attributes: { "wd:Primary": true },
              ref: makeRef("HOME"),
            },
          ],
        },
      ],
    }));
  }
  return emailAddressData;
}

build().then(console.log);
