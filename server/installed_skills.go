package server

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/Autumn-27/norma/skill"
)

// installedSkillNames uses the same parser as the runtime, not merely folder
// existence. A broken/empty directory does not close a capability gap.
func (s *Server) installedSkillNames() ([]string, error) {
	if _, err := os.Stat(s.skillDir); errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	reg, err := skill.LoadDir(s.skillDir)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	if reg != nil {
		for _, sk := range reg.List() {
			if sk.Description == "" {
				continue
			}
			names[sk.Name] = true
			if sk.Dir != "" {
				names[filepath.Base(sk.Dir)] = true
			}
		}
	}
	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}
