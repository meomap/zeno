package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchPath(t *testing.T) {
	for k, c := range []struct {
		path     string
		haystack []string
		ok       bool
	}{
		{path: "foo", haystack: []string{"foo", "bar"}, ok: true},
		{path: "bar", haystack: []string{"foo", "bar"}, ok: true},
		{path: "foo", haystack: []string{"foo/bar"}, ok: true},
		{path: "bar", haystack: []string{"foo/bar"}, ok: false},
		{path: "/foo", haystack: []string{"/foo/bar", "bar"}, ok: true},
		{path: "/bar", haystack: []string{"/foo/bar", "bar"}, ok: false},
		{path: "foo", haystack: []string{"bar"}, ok: false},
		{path: "bar", haystack: []string{"bar"}, ok: true},
		{path: "foo", haystack: []string{}, ok: false},
	} {
		assert.Equal(t, c.ok, matchPath(c.path, c.haystack), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}
