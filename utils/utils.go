package utils

import (
	"os"
	"strings"
)

// FileExists checks whether given file is present.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StringSorter implements sort.Interface for []string.
type StringSorter []string

func (a StringSorter) Len() int           { return len(a) }
func (a StringSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StringSorter) Less(i, j int) bool { return a[i] < a[j] }

// Contains tests whether string e is in slice s.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// UnderTest returns true when gscloud is running within 'Go test'.
func UnderTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
