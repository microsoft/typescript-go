package semver

import (
	"testing"
)

func FuzzTryParseVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.2.3-pre.4+build.5")
	f.Add("0.0.0")
	f.Add("999.999.999")
	f.Add("1.0.0-alpha-beta")
	f.Add("1.0.0-0.3.7")
	f.Add("1.0.0+20130313144700")
	f.Add("")
	f.Add("v1.2.3")
	f.Add("01.2.3")
	f.Add("1.2.3-")
	f.Add("1.2.3+")
	f.Add("1.2.3-.pre")
	f.Add("1.2.3-pre..rel")
	f.Add("1.2.3-01")

	f.Fuzz(func(t *testing.T, s string) {
		v, err := TryParseVersion(s)
		if err != nil {
			return
		}
		// String() must not panic
		str := v.String()
		// Parsed version must round-trip through String() and re-parse to the same value
		v2, err := TryParseVersion(str)
		if err != nil {
			t.Fatalf("round-trip failed: %q -> %q -> error: %v", s, str, err)
		}
		if v.Compare(&v2) != 0 {
			t.Fatalf("round-trip compare mismatch: %q -> %q -> %q", s, str, v2.String())
		}
	})
}

func FuzzTryParseVersionRange(f *testing.F) {
	f.Add(">=1.0.0 <2.0.0")
	f.Add(">=1.0.0 <2.0.0 || >=3.0.0")
	f.Add("1.0.0 - 2.0.0")
	f.Add("~1.2.3")
	f.Add("^1.2.3")
	f.Add("^0.0.1")
	f.Add("*")
	f.Add("")
	f.Add("1")
	f.Add("1.2")
	f.Add("1.2.3")
	f.Add("<*")
	f.Add(">*")
	f.Add(">=1.2.3-0")
	f.Add("abc")
	f.Add(">>1.0.0")
	f.Add("!1.0.0")

	f.Fuzz(func(t *testing.T, s string) {
		vr, ok := TryParseVersionRange(s)
		if !ok {
			return
		}
		// String() must not panic
		_ = vr.String()
	})
}

func FuzzVersionRangeTest(f *testing.F) {
	f.Add(">=1.0.0 <2.0.0", "1.5.0")
	f.Add("^1.2.3", "1.9.0")
	f.Add("~1.2.3", "1.2.9")
	f.Add("1.0.0 - 2.0.0", "1.5.0")
	f.Add("*", "0.0.0")
	f.Add(">=1.0.0", "2.0.0-beta")
	f.Add("<1.0.0", "0.9.9")

	f.Fuzz(func(t *testing.T, rangeStr, versionStr string) {
		vr, ok := TryParseVersionRange(rangeStr)
		if !ok {
			return
		}
		v, err := TryParseVersion(versionStr)
		if err != nil {
			return
		}
		// Test must not panic
		_ = vr.Test(&v)
	})
}

func FuzzVersionCompare(f *testing.F) {
	f.Add("1.0.0", "2.0.0")
	f.Add("1.0.0-alpha", "1.0.0-beta")
	f.Add("1.0.0+build", "1.0.0")
	f.Add("0.0.0-0", "0.0.0")

	f.Fuzz(func(t *testing.T, a, b string) {
		va, errA := TryParseVersion(a)
		if errA != nil {
			return
		}
		vb, errB := TryParseVersion(b)
		if errB != nil {
			return
		}
		// Compare must not panic
		cmpAB := va.Compare(&vb)
		cmpBA := vb.Compare(&va)
		// Antisymmetry: if a < b then b > a
		if (cmpAB < 0 && cmpBA <= 0) || (cmpAB > 0 && cmpBA >= 0) || (cmpAB == 0 && cmpBA != 0) {
			t.Fatalf("antisymmetry violation: Compare(%q,%q)=%d but Compare(%q,%q)=%d", a, b, cmpAB, b, a, cmpBA)
		}
		// Reflexivity
		if va.Compare(&va) != 0 {
			t.Fatalf("reflexivity violation: Compare(%q,%q)!=0", a, a)
		}
	})
}
