package core

import (
	"testing"
)

func TestPatternMatchesWithZeroStarIndex(t *testing.T) {
	// Create a pattern with StarIndex = 0 and an empty Text
	pattern := Pattern{
		Text:      "",
		StarIndex: 0,
	}

	// This should no longer panic after the fix
	result := pattern.Matches("some-candidate")
	
	// Verify the result is as expected
	// With our fix, this should return true if the candidate ends with the suffix
	// which is everything after the star (p.Text[p.StarIndex+1:])
	// In this case, the suffix is an empty string, so any candidate should match
	if !result {
		t.Errorf("Expected pattern to match, but it didn't")
	}
}

func TestTryParsePatternWithLeadingStar(t *testing.T) {
	// This is similar to the pattern that might be created when parsing "*?url"
	pattern := TryParsePattern("*?url")
	
	if pattern.StarIndex != 0 {
		t.Errorf("Expected StarIndex to be 0, got %d", pattern.StarIndex)
	}
	
	// Verify the pattern is considered valid
	if !pattern.IsValid() {
		t.Errorf("Pattern should be valid")
	}
}

func TestCircularModuleReference(t *testing.T) {
	// This test simulates the scenario with the eslint.config.js file
	// that contains: export { default } from '../eslint.config.js'
	
	// Create a pattern that might be used in module resolution
	// The key issue is when Text is empty but StarIndex is 0
	pattern := Pattern{
		Text:      "",
		StarIndex: 0,
	}
	
	// This should no longer panic after the fix
	result := pattern.Matches("../eslint.config.js")
	
	// Verify the result is as expected
	// With our fix, this should return true if the candidate ends with the suffix
	// which is everything after the star (p.Text[p.StarIndex+1:])
	// In this case, the suffix is an empty string, so any candidate should match
	if !result {
		t.Errorf("Expected pattern to match, but it didn't")
	}
}