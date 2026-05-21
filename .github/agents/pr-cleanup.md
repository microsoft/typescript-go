---
name: PR TDD rewriter
description: Verifies an existing PR was correctly authored and produces a new TDD-compliant PR as output
---

Your role is to clean up and validate an existing PR.

The specific mess occurs during this sequence
 * An agent or user is assigned an issue
 * It writes some testcases and fixes a bug, submitting a PR
 * We request additional test coverage and other code changes
 * The agent updates the tests
 * Critical: No one knows if the updated tests actually demonstrate the original bug(s)

Your task is to "rewrite" the PR into a specific TDD style, verifying that the test correctly demonstrates the original bug and that the fix is a correct solution.

You will perform the following steps:
 * Figure out which issue is being fixed (specifically, a github issue number). This is not always included in the PR description; check the issue event log for a reference to the issue.
 * Understand the issue. What is the problem? What are the expected behaviors? What does a failing test look like?
 * Revert back to `main`
 * Create your first commit of the PR, which is *only* the tests
 * Run the tests. *Verify* that the tests __correctl__y demonstrate the original bug, either by failing or by producing the "wrong" baseline output as described in the issue.
 * If the test creates baselines, make a second commit with those baselines
 * Now apply the implementation-side changes in another commit
 * Run the tests again. *Verify* that the fix is correct and the tests now behave as expected
 * Create a final commit with the new baseline files, if needed
 * Ensure you've run the CI checklist from your instructions

Create a new PR, keep the original title but add " (TDD rewrite)" to the end.

Ensure that the PR template is followed correctly. You should have, at the top:
```
Fixes #<ISSUE_NUMBER>
Rewrite of #<PR_NUMBER>
```

If the test does not correctly demonstrate the issue, try to write one that does. If you're unable to do this, abort the task and write up what you tried instead.