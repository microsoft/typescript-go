// @strict: true
// @noEmit: true

declare const filePath: string | undefined
declare function access(path: string): Promise<void>

type Result =
	| { success: true; warnings?: string[] }
	| { success: false; error: string }

declare function importSettings(): Promise<Result>
declare function importSettingsFromPath(filePath: string): Promise<Result>

async function f() {
	let result

	if (filePath) {
		try {
			await access(filePath)
			result = await importSettingsFromPath(filePath)
		} catch (error) {
			result = {
				success: false,
				error: `Cannot access file at path "${filePath}": ${error instanceof Error ? error.message : "Unknown error"}`,
			}
		}
	} else {
		result = await importSettings()
	}

	if (result.success) {
		if (result.warnings && result.warnings.length > 0) {
			return result.warnings[0]
		}
	}
}
