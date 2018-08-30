package search

import (
	"github.com/pkg/errors"

	"github.com/meomap/zeno/loader"
	"github.com/meomap/zeno/parser"
)

// MatchPlaybook reports whether pb appear in affected changes
func MatchPlaybook(pb string, files []string, root string, ds loader.DataSource) (bool, error) {
	deps, err := parser.ParsePlaybook(pb, root, ds)
	if err != nil {
		return false, errors.Wrapf(err, "parser.ParsePlaybook pb=%s root=%s", pb, root)
	}
	for _, v := range deps {
		if matchPath(v, files) {
			return true, nil
		}
	}
	return false, nil
}
