import * as cp from "child_process";
import * as fs from "fs";
import path from "path";
import * as readline from "readline";
import which from "which";

const failingTestsPath = path.join(import.meta.dirname, "failingTests.txt");

interface TestEvent {
    Time?: string;
    Action: string;
    Package?: string;
    Test?: string;
    Output?: string;
    Elapsed?: number;
}

async function main() {
    const oldFailingTests = fs.readFileSync(failingTestsPath, "utf-8");
    const go = which.sync("go");

    let testProcess: cp.ChildProcess;
    try {
        // Run tests with TSGO_FOURSLASH_IGNORE_FAILING=1 to run all tests including those in failingTests.txt
        testProcess = cp.spawn(go, ["test", "-json", "./internal/fourslash/tests/gen"], {
            stdio: ["ignore", "pipe", "pipe"],
            env: { ...process.env, TSGO_FOURSLASH_IGNORE_FAILING: "1" },
        });
    }
    catch (error) {
        fs.writeFileSync(failingTestsPath, oldFailingTests, "utf-8");
        throw new Error("Failed to spawn test process: " + error);
    }

    if (!testProcess.stdout || !testProcess.stderr) {
        throw new Error("Test process stdout or stderr is null");
    }

    const failingTests: string[] = [];
    const testOutputs = new Map<string, string[]>();
    const allOutputs: string[] = [];
    let hadPanic = false;

    const rl = readline.createInterface({
        input: testProcess.stdout,
        crlfDelay: Infinity,
    });

    rl.on("line", line => {
        try {
            const event: TestEvent = JSON.parse(line);

            // Collect output for each test
            if (event.Action === "output" && event.Output) {
                allOutputs.push(event.Output);
                if (event.Test) {
                    if (!testOutputs.has(event.Test)) {
                        testOutputs.set(event.Test, []);
                    }
                    testOutputs.get(event.Test)!.push(event.Output);
                }

                // Check for panics
                if (/^panic/m.test(event.Output)) {
                    hadPanic = true;
                }
            }

            // Process failed tests
            if (event.Action === "fail" && event.Test) {
                const outputs = testOutputs.get(event.Test) || [];
                const fullOutput = outputs.join("");

                // Check if this is a baseline-only failure
                // Baseline failures contain specific messages from baseline.go
                const hasBaselineMessage = /new baseline created at/.test(fullOutput) ||
                    /the baseline file .* has changed/.test(fullOutput);

                // Check for non-baseline errors
                // Look for patterns that indicate real test failures
                // We need to filter out baseline messages when checking for errors
                const outputWithoutBaseline = fullOutput
                    .replace(/the baseline file .* has changed\. \(Run `hereby baseline-accept` if the new baseline is correct\.\)/g, "")
                    .replace(/new baseline created at .*\./g, "")
                    .replace(/the baseline file .* does not exist in the TypeScript submodule/g, "")
                    .replace(/the baseline file .* does not match the reference in the TypeScript submodule/g, "");

                const hasNonBaselineError = /^panic/m.test(outputWithoutBaseline) ||
                    /Error|error/i.test(outputWithoutBaseline) ||
                    /fatal|Fatal/.test(outputWithoutBaseline) ||
                    /Unexpected/.test(outputWithoutBaseline);

                // Only mark as failing if it's not a baseline-only failure
                // (i.e., if there's no baseline message, or if there are other errors besides baseline)
                if (!hasBaselineMessage || hasNonBaselineError) {
                    failingTests.push(event.Test);
                }
            }
        }
        catch (e) {
            // Not JSON, possibly stderr or other output - ignore
        }
    });

    testProcess.stderr.on("data", data => {
        // Check stderr for panics too
        const output = data.toString();
        allOutputs.push(output);
        if (/^panic/m.test(output)) {
            hadPanic = true;
        }
    });

    await new Promise<void>((resolve, reject) => {
        testProcess.on("close", code => {
            if (hadPanic) {
                fs.writeFileSync(failingTestsPath, oldFailingTests, "utf-8");
                reject(new Error("Unrecovered panic detected in tests\n" + allOutputs.join("")));
                return;
            }

            fs.writeFileSync(failingTestsPath, failingTests.sort((a, b) => a.localeCompare(b, "en-US")).join("\n") + "\n", "utf-8");
            resolve();
        });

        testProcess.on("error", error => {
            fs.writeFileSync(failingTestsPath, oldFailingTests, "utf-8");
            reject(error);
        });
    });
}

main().catch(error => {
    console.error("Error:", error);
    process.exit(1);
});
