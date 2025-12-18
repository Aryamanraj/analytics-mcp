package update

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b   string
		expect int
		ok     bool
	}{
		{"1.2.3", "1.2.3", 0, true},
		{"1.2.3", "1.2.4", -1, true},
		{"1.3.0", "1.2.9", 1, true},
		{"2.0.0", "1.9.9", 1, true},
		{"1.0", "1.0.0", 0, false},
		{"1.a.0", "1.0.0", 0, false},
	}

	for _, tc := range cases {
		got, err := CompareVersions(tc.a, tc.b)
		if tc.ok && err != nil {
			t.Fatalf("compare %s vs %s unexpected error: %v", tc.a, tc.b, err)
		}
		if !tc.ok {
			if err == nil {
				t.Fatalf("compare %s vs %s expected error", tc.a, tc.b)
			}
			continue
		}
		if got != tc.expect {
			t.Fatalf("compare %s vs %s expected %d got %d", tc.a, tc.b, tc.expect, got)
		}
	}
}

func TestMatchesMax(t *testing.T) {
	cases := []struct {
		version string
		max     string
		ok      bool
		allow   bool
	}{
		{"1.13.4", "1.13.x", true, true},
		{"1.14.0", "1.13.x", true, false},
		{"1.13.0", "1.13.5", true, true},
		{"1.13.6", "1.13.5", true, false},
		{"bad", "1.13.5", false, false},
		{"1.13.4", "bad", false, false},
	}

	for _, tc := range cases {
		allowed, err := MatchesMax(tc.version, tc.max)
		if tc.ok && err != nil {
			t.Fatalf("matches %s %s unexpected error: %v", tc.version, tc.max, err)
		}
		if !tc.ok {
			if err == nil {
				t.Fatalf("matches %s %s expected error", tc.version, tc.max)
			}
			continue
		}
		if allowed != tc.allow {
			t.Fatalf("matches %s %s expected %v got %v", tc.version, tc.max, tc.allow, allowed)
		}
	}
}

func TestIsCompatible(t *testing.T) {
	cases := []struct {
		core, min, max string
		compatible     bool
	}{
		{"1.12.3", "1.12.0", "1.13.x", true},
		{"1.11.9", "1.12.0", "1.13.x", false},
		{"1.14.0", "1.12.0", "1.13.x", false},
		{"1.13.5", "", "1.13.5", true},
		{"1.13.6", "", "1.13.5", false},
		{"bad", "1.12.0", "1.13.x", false},
	}

	for _, tc := range cases {
		ok, _ := IsCompatible(tc.core, tc.min, tc.max)
		if ok != tc.compatible {
			t.Fatalf("compat %s min %s max %s expected %v got %v", tc.core, tc.min, tc.max, tc.compatible, ok)
		}
	}
}
