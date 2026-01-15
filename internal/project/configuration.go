package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ls/lsutil"
)

type Config struct {
	mu sync.Mutex
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
	c.mu.Lock()
	defer c.mu.Unlock()
	return &Config{
		ts: c.ts.CopyOrDefault(),
		js: c.js.CopyOrDefault(),
	}
}

// any non-nil field in b is copied into a
func (a *Config) CopyInto(b *Config) *Config {
	a.mu.Lock()
	defer a.mu.Unlock()
	if b.ts != nil {
		a.ts = b.ts.Copy()
	}
	if b.js != nil {
		a.js = b.js.Copy()
	}
	return a
}

func (c *Config) Ts() *lsutil.UserPreferences {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ts.CopyOrDefault()
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
