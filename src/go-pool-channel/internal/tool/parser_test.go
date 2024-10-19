package tool_test

import (
	"lvciot/go-pool-channel/internal/tool"
	"lvciot/shared/test"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
)

func Test_ParserJustMillion(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	src := filepath.Join(b, "../../../../../data/testcases/measurements-1m.txt")
	dst := filepath.Join(b, "../../../../../data/testcases/measurements-1m.out")
	actual := "measurements.out"

	var c atomic.Uint32

	tool.Parser(src, actual, &c)
	if !test.FileTextCompare(dst, actual) {
		t.Errorf("File 'measurements-1m' is not as expected")
	}
}

func Test_Parser(t *testing.T) {
	var c atomic.Uint32
	tcs, err := test.GetCases()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, tc := range tcs {
		actualFile := "./out/" + filepath.Base(tc.ExpectedFile)
		tool.Parser(tc.SourceFile, actualFile, &c)

		if !test.FileDeepCompare(tc.ExpectedFile, actualFile) {
			t.Errorf("file '%s' is not as expected", tc.Name)
		}
	}
}
