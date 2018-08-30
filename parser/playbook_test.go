package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meomap/zeno/loader"
)

func TestParsePlaybook(t *testing.T) {
	ds := new(loader.MemoryLoader)
	for _, c := range []struct {
		caseName string
		playbook string
		setup    func()
		want     []string
	}{
		{
			caseName: "playbook_with_empty_role",
			playbook: "empty.yml",
			setup: func() {
				ds.SetFile("empty.yml", []byte(`
- name: Test empty playbook
  hosts: all`))
			},

			want: []string{},
		},
		{
			caseName: "playbook_with_single_role_explicit_path",
			playbook: "single_role_explicit.yml",
			setup: func() {
				ds.SetFile("single_role_explicit.yml", []byte(`
- name: Test single role explicitly
  hosts: all
  roles:
  - role: roles/r1
`))
				ds.SetFile("roles/r1", []byte(""))
			},
			want: []string{"roles/r1"},
		},
		{
			caseName: "playbook_with_single_role_implicit_path",
			playbook: "single_role_implicit.yml",
			setup: func() {
				ds.SetFile("single_role_implicit.yml", []byte(`
- name: Test single role implicit
  hosts: all
  roles:
  - role: r2
`))
				ds.SetFile("roles/r2", []byte(""))
			},
			want: []string{"roles/r2"},
		},
		{
			caseName: "playbook_with_multiple_roles",
			playbook: "multiple_roles.yml",
			setup: func() {
				ds.SetFile("multiple_roles.yml", []byte(`
- name: Test multiple roles
  hosts: all
  roles:
  - role: r1
  - role: r2
`))
				ds.SetFile("roles/r1", []byte(""))
				ds.SetFile("roles/r2", []byte(""))
			},
			want: []string{"roles/r1", "roles/r2"},
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := ParsePlaybook(c.playbook, "", ds)
			require.NoError(t, err)
			assert.Equal(t, c.want, out)
		})
	}
}

func TestSearchRolePath(t *testing.T) {
	ds := new(loader.MemoryLoader)
	for _, c := range []struct {
		caseName string
		role     string
		setup    func()
		want     string
	}{
		{
			caseName: "explicitly_declared_within_roles_dir",
			role:     "roles/test-role",
			setup: func() {
				ds.SetFile("roles/test-role", []byte(""))
			},
			want: "roles/test-role",
		},
		{
			caseName: "implicitly_declared_within_roles_dir",
			role:     "test-role",
			setup: func() {
				ds.SetFile("roles/test-role", []byte(""))
			},
			want: "roles/test-role",
		},
		{
			caseName: "relative_path_to_base_dir",
			role:     "../other/another-role",
			setup: func() {
				ds.SetFile("other/another-role", []byte(""))
			},
			want: "other/another-role",
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := searchRolePath(c.role, "", ds)
			require.NoError(t, err)
			assert.Equal(t, c.want, out)
		})
	}
}

func TestParseTask(t *testing.T) {
	ds := new(loader.MemoryLoader)
	for _, c := range []struct {
		caseName string
		task     string
		baseDir  string
		setup    func()
		want     []string
	}{
		{
			caseName: "file_with_no_includes",
			task:     "no-includes.yml",
			setup: func() {
				ds.SetFile("no-includes.yml", []byte(`
- name: Assert true is not false
  assert:
    that: 1 != 0`))
			},
			want: []string{},
		},
		{
			caseName: "file_with_include_tasks_in_same_dir",
			task:     "with-include_tasks.yml",
			setup: func() {
				ds.SetFile("with-include_tasks.yml", []byte(`
- name: Do something
  include_tasks: something.yml`))
				ds.SetFile("something.yml", []byte(``))
			},
			want: []string{},
		},
		{
			caseName: "file_with_include_tasks_in_relative_dir",
			task:     "with-include_tasks-relative-dir.yml",
			baseDir:  "/tmp/r1/tasks",
			setup: func() {
				ds.SetFile("/tmp/r1/tasks/with-include_tasks-relative-dir.yml", []byte(`
- name: Do more things
  include_tasks: ../../morethings.yml`))
				ds.SetFile("/tmp/morethings.yml", []byte(``))
			},
			want: []string{"/tmp/morethings.yml"},
		},
		{
			caseName: "file_with_import_tasks",
			task:     "with-import_tasks.yml",
			baseDir:  "/tmp/r2/tasks",
			setup: func() {
				ds.SetFile("/tmp/r2/tasks/with-import_tasks.yml", []byte(`
- name: Do import
  import_tasks: ../../staticthing.yml`))
				ds.SetFile("/tmp/staticthing.yml", []byte(``))
			},
			want: []string{"/tmp/staticthing.yml"},
		},
		{
			caseName: "file_with_depricated_include",
			task:     "with-depricated_include.yml",
			baseDir:  "/tmp/r3/tasks",
			setup: func() {
				ds.SetFile("/tmp/r3/tasks/with-depricated_include.yml", []byte(`
- name: Do depricated include
  include: ../../depricatedthing.yml`))
				ds.SetFile("/tmp/depricatedthing.yml", []byte(``))
			},
			want: []string{"/tmp/depricatedthing.yml"},
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			root := c.baseDir
			if root == "" {
				root = "."
			}
			out, err := parseTask(c.task, root, ds)
			require.NoError(t, err)
			assert.Equal(t, c.want, out)
		})
	}
}

func TestParseRole(t *testing.T) {
	ds := new(loader.MemoryLoader)
	for _, c := range []struct {
		caseName string
		role     string
		setup    func()
		want     []string
	}{
		{
			caseName: "role_with_empty_tasks",
			role:     "empty",
			setup: func() {
				ds.SetFile("roles/empty", []byte(""))
			},
			want: []string{"roles/empty"},
		},
		{
			caseName: "role_with_multiple_tasks",
			role:     "multiple-tasks",
			setup: func() {
				ds.SetFile("roles/multiple-tasks/tasks/t1.yml", []byte(""))
				ds.SetFile("roles/multiple-tasks/tasks/t2.yml", []byte(""))
			},
			want: []string{"roles/multiple-tasks"},
		},
		{
			caseName: "role_with_include_tasks",
			role:     "include-tasks",
			setup: func() {
				ds.SetFile("test/something.yml", []byte(""))
				ds.SetFile("roles/include-tasks/tasks/t1.yml", []byte(`
- name: Do include tasks
  include: ../../../test/something.yml`))
			},
			want: []string{"roles/include-tasks", "test/something.yml"},
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := parseRole(c.role, "", ds)
			require.NoError(t, err)
			assert.Equal(t, c.want, out)
		})
	}
}
