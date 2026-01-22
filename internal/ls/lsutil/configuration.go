package lsutil

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type UserPreferenceConfig struct {
	js *UserPreferences
	ts *UserPreferences
	// tsserverOptions
}

// if `userPreferences` is nil, this function will return a config with default userPreferences
func NewUserPreferenceConfig(userPreferences *UserPreferences) *UserPreferenceConfig {
	return &UserPreferenceConfig{
		js: userPreferences.CopyOrDefault(),
		ts: userPreferences.CopyOrDefault(),
	}
}

func (c *UserPreferenceConfig) Copy() *UserPreferenceConfig {
	return &UserPreferenceConfig{
		ts: c.ts.CopyOrDefault(),
		js: c.js.CopyOrDefault(),
	}
}

// any non-nil field in b is copied into a
func (a *UserPreferenceConfig) CopyInto(b *UserPreferenceConfig) *UserPreferenceConfig {
	newAllPreferences := &UserPreferenceConfig{}

	if b.ts != nil {
		newAllPreferences.ts = b.ts
	} else {
		newAllPreferences.ts = a.ts
	}

	if b.js != nil {
		newAllPreferences.js = b.js
	} else {
		newAllPreferences.js = a.js
	}

	return newAllPreferences
}

func (c *UserPreferenceConfig) Ts() *UserPreferences {
	return c.ts
}

func (c *UserPreferenceConfig) Js() *UserPreferences {
	return c.js
}

func (c *UserPreferenceConfig) GetPreferences(activeFile string) *UserPreferences {
	fileEnding := strings.TrimPrefix(tspath.GetAnyExtensionFromPath(activeFile, nil, true), ".")
	if tspath.ExtensionIsTs(fileEnding) {
		if c.ts != nil {
			return c.ts
		} else if c.js != nil {
			return c.js
		}
	} else {
		if c.js != nil {
			return c.js
		} else if c.ts != nil {
			return c.ts
		}
	}
	return NewDefaultUserPreferences()
}

func ParseNewUserPreferenceConfiguration(items []any) *UserPreferenceConfig {
	defaultPref := NewUserPreferenceConfig(nil)
	c := NewUserPreferenceConfig(nil)
	for i, item := range items {
		if item == nil {
			// continue
		} else if config, ok := item.(map[string]any); ok {
			newConfig := &UserPreferenceConfig{}
			if i < 2 {
				newConfig.ts = defaultPref.ts.ParseWorker(config)
			} else {
				newConfig.js = defaultPref.js.ParseWorker(config)
			}
			c = c.CopyInto(newConfig)
		} else if item, ok := item.(*UserPreferences); ok {
			// case for fourslash -- fourslash sends the entire userPreferences over
			// !!! support format and js/ts distinction?
			return NewUserPreferenceConfig(item)
		}
	}
	return c
}
