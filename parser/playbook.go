package parser

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	"github.com/meomap/zeno/loader"
)

// Task with file includes
type Task struct {
	Name         string `yaml:"name"`
	IncludeTasks string `yaml:"include_tasks"`
	ImportTasks  string `yaml:"import_tasks"`
	Include      string `yaml:"include"`
}

// Role may define tasks include/import
type Role struct {
	Name string `yaml:"role"`
}

// Play composites of multiple roles & tasks
type Play struct {
	Roles []Role `yaml:"roles"`
}

// ParsePlaybook returns list of dirs/files used by current playbook
func ParsePlaybook(filePath string, repoDir string, ds loader.DataSource) ([]string, error) {
	log.Printf("Parse playbook '%s'", filePath)
	content, err := ds.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "dataSource file_path=%s", filePath)
	}
	playbook := []Play{}
	if err = yaml.Unmarshal(content, &playbook); err != nil {
		return nil, errors.Wrapf(err, "yaml.Unmarshal file_path=%s", filePath)
	}
	playbookRoot := filepath.Dir(path.Join(repoDir, filePath))
	deps := []string{}
	for _, play := range playbook {
		for _, role := range play.Roles {
			roleDeps, rErr := parseRole(role.Name, playbookRoot, ds)
			if rErr != nil {
				return nil, errors.Wrapf(rErr, "parseRole name=%s", role.Name)
			}
			deps = append(deps, roleDeps...)
		}
	}
	log.Printf("Dependencies: %+v", deps)
	return deps, nil
}

func parseRole(name string, playbookRoot string, ds loader.DataSource) ([]string, error) {
	// log.Printf("Parse role '%s' root=%s", name, playbookRoot)
	rPath, err := searchRolePath(name, playbookRoot, ds)
	if err != nil {
		return nil, errors.Wrapf(err, "searchRolePath name=%s", name)
	}
	// all files containing path prefix that matched
	deps := []string{rPath}

	// fetch all task includes/imports
	taskRoot := path.Join(rPath, "tasks")
	includeFiles, err := ds.ReadDir(taskRoot)
	if err != nil {
		if os.IsNotExist(err) {
			// no need explore more
			return deps, nil
		}
		return nil, errors.Wrapf(err, "ds.ReadDir dir_path=%s", taskRoot)
	}
	// fetch task file content
	for _, incPath := range includeFiles {
		tDeps, tErr := parseTask(incPath, taskRoot, ds)
		if tErr != nil {
			return nil, errors.Wrapf(tErr, "parseTask path=%s", incPath)
		}
		deps = append(deps, tDeps...)
	}
	return deps, nil
}

// looking for files from other than current root only
func parseTask(name string, root string, ds loader.DataSource) ([]string, error) {
	// log.Printf("Parse task '%s' root=%s", name, root)
	filePath := path.Join(root, name)
	deps := []string{}
	baseDir := path.Dir(filePath)
	if baseDir != root {
		deps = append(deps, filePath)
	}

	content, err := ds.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "dataSource file_path=%s", filePath)
	}
	taskList := []Task{}
	if err = yaml.Unmarshal(content, &taskList); err != nil {
		return nil, errors.Wrapf(err, "yaml.Unmarshal file_path=%s", filePath)
	}

	parseInclude := func(name string) error {
		iDeps, iErr := parseTask(name, root, ds)
		if iErr != nil {
			return errors.Wrapf(iErr, "parseTask name=%s", name)
		}
		deps = append(deps, iDeps...)
		return nil
	}

	for _, task := range taskList {
		if task.IncludeTasks != "" {
			if err = parseInclude(task.IncludeTasks); err != nil {
				return nil, errors.Wrapf(err, "parseInclude include_tasks=%s", task.IncludeTasks)
			}
		}
		if task.ImportTasks != "" {
			if err = parseInclude(task.ImportTasks); err != nil {
				return nil, errors.Wrapf(err, "parseInclude import_tasks=%s", task.ImportTasks)
			}
		}
		if task.Include != "" {
			if err = parseInclude(task.Include); err != nil {
				return nil, errors.Wrapf(err, "parseInclude include=%s", task.Include)
			}
		}
	}
	return deps, nil
}

// role name could be directory path relative to playbook base dir `roles`,
// or without `roles/` dir. note that paths from DEFAULT_ROLES_PATH are ignored
func searchRolePath(name string, baseDir string, ds loader.DataSource) (string, error) {
	searchPaths := []string{baseDir, path.Join(baseDir, "roles")}
	for _, p := range searchPaths {
		rPath := path.Join(p, name)
		if exist, err := ds.IsExist(rPath); err != nil {
			return "", errors.Wrapf(err, "ds.IsExist path=%s", rPath)
		} else if exist {
			return rPath, nil
		}
	}
	return "", errors.Errorf("role %s was not found in %+v", name, searchPaths)
}
