package lsutil

func ParseUserConfig(items map[string]any) UserPreferences {
	prefs := NewDefaultUserPreferences()
	if editorItem, ok := items["editor"]; ok && editorItem != nil {
		if editorSettings, ok := editorItem.(map[string]any); ok {
			prefs.FormatCodeSettings.ParseEditorSettings(editorSettings)
		}
	}
	if jsTsItem, ok := items["js/ts"]; ok && jsTsItem != nil {
		switch jsTsSettings := jsTsItem.(type) {
		case map[string]any:
			prefs.ParseWorker(jsTsSettings)
		case UserPreferences:
			prefs.MergeNonDefaults(&jsTsSettings)
		case *UserPreferences:
			prefs.MergeNonDefaults(jsTsSettings)
		}
	}
	return prefs
}
