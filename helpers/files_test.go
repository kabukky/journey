package helpers

import "testing"

var filenameTests = []struct {
	in  string
	out string
}{
	{
		in:  "test.txt",
		out: "test",
	},
	{
		in:  "another.test.ext",
		out: "another.test",
	},
	{
		in:  "no_ext",
		out: "no_ext",
	},
	{
		in:  "",
		out: ".",
	},
}

var directoryTests = []struct {
	in  string
	out bool
}{
	{
		in:  ".",
		out: true,
	},
	{
		in:  "..",
		out: true,
	},
	{
		in:  "nowaythisdirectoryactuallyexists",
		out: false,
	},
}

func TestGetFilenameWithoutExtension(t *testing.T) {
	for _, test := range filenameTests {
		if actual := GetFilenameWithoutExtension(test.in); actual != test.out {
			t.Errorf("Expected '%s', received '%s' for '%s'", test.out, actual, test.in)
		}
	}
}

func TestIsDirectory(t *testing.T) {
	for _, test := range directoryTests {
		if actual := IsDirectory(test.in); actual != test.out {
			t.Errorf("Expected %t, received %t for '%s'", test.out, actual, test.in)
		}
	}
}
