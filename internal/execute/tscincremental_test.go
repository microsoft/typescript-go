package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
)

func TestIncremental(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	testCases := []*tscInput{
		{
			subScenario: "serializing error chain",
			sys: newTestSys(FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{
                    "compilerOptions": {
                        "incremental": true,
                        "strict": true,
                        "jsx": "react",
                        "module": "esnext",
                    },
                }`,
				"/home/src/workspaces/project/index.tsx": `
                    declare namespace JSX {
                        interface ElementChildrenAttribute { children: {}; }
                        interface IntrinsicElements { div: {} }
                    }

                    declare var React: any;

                    declare function Component(props: never): any;
                    declare function Component(props: { children?: number }): any;
                    (<Component>
                        <div />
                        <div />
                    </Component>)`,
			}, "/home/src/workspaces/project"),
		},
		{
			subScenario: "change to modifier of class expression field with declaration emit enabled",
			sys: newTestSys(FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "module": "esnext", "declaration": true  } }`,
				"/home/src/workspaces/project/main.ts": `
                        import MessageablePerson from './MessageablePerson.js';
                        function logMessage( person: MessageablePerson ) {
                            console.log( person.message );
                        }`,
				"/home/src/workspaces/project/MessageablePerson.ts": `
                        const Messageable = () => {
                            return class MessageableClass {
                                public message = 'hello';
                            }
                        };
                        const wrapper = () => Messageable();
                        type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
                        export default MessageablePerson;`,
			}, "/home/src/workspaces/project"),
			commandLineArgs: []string{"--incremental"},
			// edits: [
			//     noChangeRun,
			//     {
			//         caption: "modify public to protected",
			//         edit: sys => sys.replaceFileText("/home/src/workspaces/project/MessageablePerson.ts", "public", "protected"),
			//     },
			//     noChangeRun,
			//     {
			//         caption: "modify protected to public",
			//         edit: sys => sys.replaceFileText("/home/src/workspaces/project/MessageablePerson.ts", "protected", "public"),
			//     },
			//     noChangeRun,
			// ],
		},
		{
			subScenario: "change to modifier of class expression field",
			sys: newTestSys(FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "module": "esnext" } }`,
				"/home/src/workspaces/project/main.ts": `
                        import MessageablePerson from './MessageablePerson.js';
                        function logMessage( person: MessageablePerson ) {
                            console.log( person.message );
                        }`,
				"/home/src/workspaces/project/MessageablePerson.ts": `
                        const Messageable = () => {
                            return class MessageableClass {
                                public message = 'hello';
                            }
                        };
                        const wrapper = () => Messageable();
                        type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
                        export default MessageablePerson;`,
			}, "/home/src/workspaces/project"),
			commandLineArgs: []string{"--incremental"},
			// edits: [
			//     noChangeRun,
			//     {
			//         caption: "modify public to protected",
			//         edit: sys => sys.replaceFileText("/home/src/workspaces/project/MessageablePerson.ts", "public", "protected"),
			//     },
			//     noChangeRun,
			//     {
			//         caption: "modify protected to public",
			//         edit: sys => sys.replaceFileText("/home/src/workspaces/project/MessageablePerson.ts", "protected", "public"),
			//     },
			//     noChangeRun,
			// ],
		},
	}

	for _, test := range testCases {
		test.verify(t, "incremental")
	}
}
