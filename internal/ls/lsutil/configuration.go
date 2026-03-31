package lsutil

import (
	"github.com/microsoft/typescript-go/internal/tspath"
)

type UserConfig struct {
	js UserPreferences
	ts UserPreferences
}

func NewUserConfig(userPreferences UserPreferences) *UserConfig {
	return &UserConfig{
		js: userPreferences,
		ts: userPreferences,
	}
}

func (c *UserConfig) Copy() *UserConfig {
	return &UserConfig{
		ts: c.ts,
		js: c.js,
	}
}

func (a *UserConfig) Merge(b *UserConfig) *UserConfig {
	return &UserConfig{
		ts: b.ts,
		js: b.js,
	}
}

func (c *UserConfig) TS() *UserPreferences {
	return &c.ts
}

func (c *UserConfig) JS() *UserPreferences {
	return &c.js
}

func (c *UserConfig) GetPreferences(activeFile string) *UserPreferences {
	if activeFile == "" || tspath.ExtensionIsTs(tspath.GetAnyExtensionFromPath(activeFile, nil, true)) {
		return &c.ts
	}
	return &c.js
}

func ParseNewUserConfig(items map[string]any) *UserConfig {
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
	return NewUserConfig(prefs)
}
