package kernels

import "testing"

func TestDefaultConfigGroups(t *testing.T) {
	for _, g := range DefaultConfigGroups {
		if _, ok := ConfigOptGroups[g]; !ok {
			t.Fatalf("default config group '%s' does not exist", g)
		}
	}
}
