package bfs

import (
	"github.com/orcaman/concurrent-map"
	"gitreefs/core/common"
	"gitreefs/core/git"
	"gitreefs/core/logger"
	"path"
)

type Repository struct {
	name            string
	provider        *git.RepositoryProvider
	commitishByName cmap.ConcurrentMap
}

func NewRepository(clonesPath string, name string) (repository *Repository, err error) {
	clonePath := path.Join(clonesPath, name)
	err = common.ValidateDirectory(clonePath, false)
	if err != nil {
		return
	}
	var provider *git.RepositoryProvider
	provider, err = git.NewRepositoryProvider(clonePath)
	if err != nil || provider == nil {
		return nil, err
	}
	repository = &Repository{
		name:            name,
		provider:        provider,
		commitishByName: cmap.New(),
	}
	logger.Debug("NewRepository: %v", clonePath)
	return
}

func (repository *Repository) getOrAddCommitish(name string) (commitish *Commitish, err error) {
	wrapped :=
		repository.commitishByName.Upsert(name, nil, func(found bool, existingValue interface{}, _ interface{}) interface{} {
			if found && existingValue != nil && existingValue.(*Commitish) != nil {
				return existingValue
			}
			var commitish *Commitish
			commitish, err = NewCommitish(name, repository.provider)
			return commitish
		})
	if wrapped.(*Commitish) == nil {
		return nil, err
	}
	return wrapped.(*Commitish), err
}
