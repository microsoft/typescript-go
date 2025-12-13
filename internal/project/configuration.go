package project

import (
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
)

type Config struct {
	js *lsutil.UserPreferences
	Ts *lsutil.UserPreferences
	// tsserverOptions
}

func NewConfig(userPreferences *lsutil.UserPreferences) *Config {
	// use default userPreferences if nil
	return &Config{
		js: userPreferences.CopyOrDefault(),
		Ts: userPreferences.CopyOrDefault(),
	}
}

func (c *Config) Copy() *Config {
	return &Config{
		Ts: c.Ts.CopyOrDefault(),
		js: c.js.CopyOrDefault(),
	}
}

// any non-nil field in b is copied into a
func (a *Config) CopyInto(b *Config) *Config {
	if b.Ts != nil {
		a.Ts = b.Ts.Copy()
	}
	if b.js != nil {
		a.js = b.js.Copy()
	}
	return a
}

func ParseConfiguration(items []any) *Config {
	defaultConfig := NewConfig(nil)
	c := &Config{}
	for i, item := range items {
		if item == nil {
			// continue
		} else if config, ok := item.(map[string]any); ok {
			newConfig := &Config{}
			if i < 2 {
				newConfig.Ts = defaultConfig.Ts.Copy().ParseWorker(config)
			} else {
				newConfig.js = defaultConfig.js.Copy().ParseWorker(config)
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
