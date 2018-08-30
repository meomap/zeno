package loader

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// DataSource defines common interface for data loading
type DataSource interface {
	ReadFile(string) ([]byte, error)
	ReadDir(string) ([]string, error)
	IsExist(string) (bool, error)
}

// MemoryLoader implements IO operations for testing
type MemoryLoader struct {
	files map[string][]byte
}

// ReadFile returns byte content given preload file name
func (ml MemoryLoader) ReadFile(name string) ([]byte, error) {
	content, ok := ml.files[name]
	if !ok {
		return nil, errors.Errorf("file %s not exist", name)
	}
	return content, nil
}

// SetFile update stored file's content
func (ml *MemoryLoader) SetFile(name string, content []byte) {
	if ml.files == nil {
		ml.files = map[string][]byte{}
	}
	pathComps := strings.Split(name, string(filepath.Separator))
	lenComps := len(pathComps)
	var (
		subdir []byte
		ok     bool
		parent string
	)
	for i := 0; i < lenComps-1; i++ {
		if i == 0 {
			parent = pathComps[i]
		} else {
			parent = path.Join(parent, pathComps[i])
		}
		child := pathComps[i+1]
		if subdir, ok = ml.files[parent]; !ok {
			subdir = []byte(child)
		} else {
			updated := strings.Split(string(subdir), ",")
			updated = append(updated, child)
			subdir = []byte(strings.Join(updated, ","))
		}
		ml.files[parent] = subdir
	}
	ml.files[name] = content
}

// ReadDir returns list of files' name under specified directory
func (ml MemoryLoader) ReadDir(name string) ([]string, error) {
	content, ok := ml.files[name]
	if !ok {
		return nil, &os.PathError{Err: os.ErrNotExist}
	}
	// hook to raise unexpected error
	val := string(content)
	if val == "unexpected_error" {
		ok = false
		return nil, errors.New(val)
	}
	children := strings.Split(val, ",")
	return children, nil
}

// IsExist returns true if given file name exists
func (ml MemoryLoader) IsExist(name string) (ok bool, err error) {
	var content []byte
	if content, ok = ml.files[name]; ok {
		// hook to raise unexpected error
		val := string(content)
		if val == "unexpected_error" {
			ok = false
			err = errors.New(val)
		}
	}
	return
}

// Clear reset in-mem data
func (ml *MemoryLoader) Clear() {
	ml.files = nil
}

// FileLoader implements IO operation on local disk file
type FileLoader struct {
}

// ReadFile returns byte content given file name
func (fl FileLoader) ReadFile(name string) ([]byte, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile name=%s", name)
	}
	return content, nil
}

// ReadDir returns list of files' name under specified directory
func (fl FileLoader) ReadDir(name string) (out []string, err error) {
	var stats []os.FileInfo
	if stats, err = ioutil.ReadDir(name); err != nil {
		err = errors.Wrapf(err, "ioutil.ReadDir name=%s", name)
		return
	}
	for _, v := range stats {
		out = append(out, v.Name())
	}
	return
}

// IsExist returns true if given file name exists
func (fl FileLoader) IsExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "os.Stat name=%s", name)
	}
	return true, nil
}
