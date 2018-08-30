package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meomap/zeno/loader"
)

func TestMatchPlaybook(t *testing.T) {
	ds := new(loader.MemoryLoader)
	for _, c := range []struct {
		caseName string
		playbook string
		diffs    []string
		setup    func()
		want     bool
	}{
		{
			caseName: "playbook_changed",
			playbook: "changed.yml",
			diffs:    []string{"roles/r1/t1.yml", "roles/r2/t2.yml"},
			setup: func() {
				ds.SetFile("changed.yml", []byte(`
- name: Test playbook changed
  hosts: all
  roles:
  - role: r1`))
				ds.SetFile("roles/r1/t1.yml", []byte(""))
				ds.SetFile("roles/r2/t2.yml", []byte(""))
			},
			want: true,
		},
		{
			caseName: "playbook_unchanged",
			playbook: "unchanged.yml",
			diffs:    []string{"roles/r1/t1.yml", "roles/r2/t2.yml"},
			setup: func() {
				ds.SetFile("unchanged.yml", []byte(`
- name: Test playbook changed
  hosts: all
  roles:
  - role: r3`))
				ds.SetFile("roles/r1/t1.yml", []byte(""))
				ds.SetFile("roles/r2/t2.yml", []byte(""))
				ds.SetFile("roles/r3/t3.yml", []byte(""))
			},
			want: false,
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.caseName), func(t *testing.T) {
			ds.Clear()
			c.setup()
			out, err := MatchPlaybook(c.playbook, c.diffs, ".", ds)
			require.NoError(t, err)
			assert.Equal(t, c.want, out)
		})
	}
}
