package project

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Config struct {
	js *lsutil.UserPreferences
	ts *lsutil.UserPreferences
	// tsserverOptions
}

// if `userPreferences` is nil, this function will return a config with default userPreferences
func NewConfig(userPreferences *lsutil.UserPreferences) *Config {
	return &Config{
		js: userPreferences.CopyOrDefault(),
		ts: userPreferences.CopyOrDefault(),
	}
}

func (c *Config) Copy() *Config {
	return &Config{
		ts: c.ts.CopyOrDefault(),
		js: c.js.CopyOrDefault(),
	}
}

// any non-nil field in b is copied into a
func (a *Config) CopyInto(b *Config) *Config {
	newConfig := &Config{}

	if b.ts != nil {
		newConfig.ts = b.ts
	} else {
		newConfig.ts = a.ts
	}

	if b.js != nil {
		newConfig.js = b.js
	} else {
		newConfig.js = a.js
	}

	return newConfig
}

func (c *Config) Ts() *lsutil.UserPreferences {
	return c.ts
}

func (c *Config) Js() *lsutil.UserPreferences {
	return c.js
}

func (c *Config) GetPreference(activeFile string) *lsutil.UserPreferences {
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
	return lsutil.NewDefaultUserPreferences()
}

func ParseConfiguration(items []any) *Config {
	defaultConfig := NewConfig(nil)
	c := NewConfig(nil)
	for i, item := range items {
		if item == nil {
			// continue
		} else if config, ok := item.(map[string]any); ok {
			newConfig := &Config{}
			if i < 2 {
				newConfig.ts = defaultConfig.ts.ParseWorker(config)
			} else {
				newConfig.js = defaultConfig.js.ParseWorker(config)
			}
			c = c.CopyInto(newConfig)
		} else if item, ok := item.(*lsutil.UserPreferences); ok {
			// case for fourslash -- fourslash sends the entire userPreferences over
			// !!! support format and js/ts distinction?
			return NewConfig(item)
		}
	}
	return c
}
