package file

import (
	"os"
	"testing"
)

func TestWithUserHomedir(t *testing.T) {
	// Patch env
	os.Setenv("HOME", "foo")
	defer os.Unsetenv("HOME")

	testCases := []struct {
		paths    []string
		fullPath string
	}{
		// When paths is one string, concat the path
		{
			paths:    []string{"bar"},
			fullPath: "foo/bar",
		},
		// When multiple paths, concat the path
		{
			paths:    []string{"bar", "baz"},
			fullPath: "foo/bar/baz",
		},
		// When paths is nil, just return homedir
		{
			paths:    nil,
			fullPath: "foo",
		},
	}

	for _, testCase := range testCases {
		out, err := WithUserHomeDir(testCase.paths...)
		if err != nil {
			t.Fatalf("Got error fetching home directory\n")
		}

		if out != testCase.fullPath {
			t.Errorf("Expected %s, got %s\n", testCase.fullPath, out)
		}
	}
}
