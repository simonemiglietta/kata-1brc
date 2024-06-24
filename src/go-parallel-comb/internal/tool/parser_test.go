package tool_test

import (
	"bytes"
	"io"
	"log"
	"lvciot/go-seq/internal/tool"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type testCase struct {
	name     string
	source   string
	expected string
}

func TestParser(t *testing.T) {
	actual := "measurements.out"
	i := 0
	tcs := newTestCases()

	for _, tc := range tcs {
		tool.Parser(tc.source, actual, &i)

		if !deepCompare(tc.expected, actual) {
			t.Fatalf("File %s is not as expected", tc.name)
		}
	}
}

func newTestCases() []testCase {
	_, b, _, _ := runtime.Caller(0)
	pattern := filepath.Join(b, "../../../../data/testcases/*.txt")

	testFiles, _ := filepath.Glob(pattern)

	tc := make([]testCase, len(testFiles))

	for i, tf := range testFiles {
		tc[i] = testCase{
			name:     filepath.Base(tf),
			source:   tf,
			expected: strings.Replace(tf, ".txt", ".out", 1),
		}
	}

	return tc
}

func deepCompare(file1, file2 string) bool {
	const chunkSize = 64000

	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
