package test

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func FileTextCompare(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	s1 := bufio.NewScanner(f1)
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	s2 := bufio.NewScanner(f2)
	defer f2.Close()

	same := false

	for s1.Scan() && s2.Scan() {
		t1 := s1.Text()
		t2 := s2.Text()

		same = t1 == t2
	}

	return same
}

func FileDeepCompare(file1, file2 string) bool {
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
