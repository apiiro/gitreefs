package fs

import (
	"github.com/orcaman/concurrent-map"
)

type Root struct {
	clonesPath         string
	repositoriesByName cmap.ConcurrentMap
}

func NewRoot(clonesPath string) (root *Root, err error) {
	return &Root{
		clonesPath:         clonesPath,
		repositoriesByName: cmap.New(),
	}, nil
}

func (root *Root) getOrAddRepository(name string) (repository *Repository, err error) {
	wrapped :=
		root.repositoriesByName.Upsert(name, nil, func(found bool, existingValue interface{}, _ interface{}) interface{} {
			if found && existingValue != nil && existingValue.(*Repository) != nil {
				return existingValue
			}
			var repository *Repository
			repository, err = NewRepository(root.clonesPath, name)
			return repository
		})
	if wrapped.(*Repository) == nil {
		return nil, err
	}
	return wrapped.(*Repository), err
}
