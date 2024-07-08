package main

import (
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type testCase struct {
	name          string
	numCores      int
	countStations int64
	duration      int64
}

func TestCores(t *testing.T) {
	maxCores := runtime.NumCPU()
	var tests []testCase
	tests = make([]testCase, maxCores)
	for i := 0; i < maxCores; i++ {
		stringCore := strconv.Itoa(i + 1)
		tests[i] = testCase{
			name:          stringCore,
			numCores:      i + 1,
			countStations: 0,
		}
	}

	_, b, _, _ := runtime.Caller(0)
	srcFile := filepath.Join(b, SrcFile)

	for i, test := range tests {

		start := time.Now()

		tests[i].countStations = int64(len(StationStats(srcFile, test.numCores)))

		duration := time.Since(start)
		tests[i].duration = duration.Milliseconds()

		if tests[i].countStations == 0 {
			t.Errorf("For args %v, got %q", test.numCores, test.countStations)
		}
	}
	for i := 0; i < len(tests)-2; i++ {
		if tests[i].countStations != tests[i+1].countStations {
			t.Errorf("tests produce different output")
		}
	}

	// Restore the original command-line arguments
	t.Logf("\n")
	for _, test := range tests {
		t.Logf("Running test on %s cores in total time: %d ms\n", test.name, test.duration)
	}
}
