package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryLoaderReadFile(t *testing.T) {
	ds := new(MemoryLoader)
	for _, c := range []struct {
		caseName string
		input    string
		setup    func()
		err      bool
		want     []byte
	}{
		{
			caseName: "file_exist",
			input:    "foo.bar",
			setup: func() {
				ds.SetFile("foo.bar", []byte(`abcde`))
			},
			want: []byte(`abcde`),
		},
		{
			caseName: "file_not_exist",
			input:    "foo.barz",
			setup:    func() {},
			err:      true,
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := ds.ReadFile(c.input)
			if c.err == true {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.want, out)
			}
		})
	}
}

func TestMemoryLoaderReadDir(t *testing.T) {
	ds := new(MemoryLoader)
	for _, c := range []struct {
		caseName string
		input    string
		setup    func()
		err      bool
		want     []string
	}{
		{
			caseName: "dir_exist",
			input:    "foo",
			setup: func() {
				ds.SetFile("foo/bar", []byte(``))
				ds.SetFile("foo/barz", []byte(``))
			},
			want: []string{"bar", "barz"},
		},
		{
			caseName: "dir_not_exist",
			input:    "foo",
			setup:    func() {},
			err:      true,
		},
		{
			caseName: "unexpected_error",
			input:    "fooz",
			setup: func() {
				ds.SetFile("fooz", []byte(`unexpected_error`))
			},
			err: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := ds.ReadDir(c.input)
			if c.err == true {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.want, out)
			}
		})
	}
}

func TestMemoryLoaderIsExist(t *testing.T) {
	ds := new(MemoryLoader)
	for _, c := range []struct {
		caseName string
		input    string
		setup    func()
		err      bool
		want     bool
	}{
		{
			caseName: "check_with_file_exist",
			input:    "foo",
			setup: func() {
				ds.SetFile("foo", []byte(``))
			},
			want: true,
		},
		{
			caseName: "check_with_file_not_exist",
			input:    "bar",
			setup:    func() {},
			want:     false,
		},
		{
			caseName: "check_with_unexpected_error_raised",
			input:    "fooz",
			setup: func() {
				ds.SetFile("fooz", []byte(`unexpected_error`))
			},
			err: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := ds.IsExist(c.input)
			if c.err == true {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.want, out)
			}
		})
	}
}

func TestFileLoader(t *testing.T) {
	ds := new(FileLoader)
	tmpfile, err := ioutil.TempFile("", "zeno-test-file-loader")
	require.NoError(t, err)
	defer func() {
		rErr := os.Remove(tmpfile.Name())
		require.NoError(t, rErr)
	}()

	tmpDir := filepath.Dir(tmpfile.Name())
	tmpContent := []byte("temporary file's content")
	_, err = tmpfile.Write(tmpContent)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	// read file
	out, err := ds.ReadFile(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, tmpContent, out)

	out, err = ds.ReadFile("abcde")
	assert.Error(t, err)

	// read dir
	lst, err := ds.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.True(t, len(lst) >= 1)

	_, err = ds.ReadDir("abcde")
	assert.Error(t, err)

	// check is exist
	ok, err := ds.IsExist(tmpDir)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = ds.IsExist("abcde")
	assert.False(t, ok)
}
