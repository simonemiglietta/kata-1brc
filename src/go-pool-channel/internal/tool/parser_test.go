package tool_test

import (
	"lvciot/go-pool-channel/internal/tool"
	"lvciot/shared/test"
	"sync/atomic"
	"testing"
)

type testCase struct {
	name     string
	source   string
	expected string
}

func TestParser(t *testing.T) {
	actualFile := "measurements.out"
	var c atomic.Uint32
	tcs := test.GetCases()

	if len(tcs) == 0 {
		t.Fatalf("no test cases found")
	}

	for _, tc := range tcs {
		tool.Parser(tc.SourceFile, actualFile, &c)

		if !test.FileDeepCompare(tc.ExpectedFile, actualFile) {
			t.Errorf("file %s is not as expected", tc.Name)
		}
	}
}
