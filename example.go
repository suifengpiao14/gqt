package gqt

import (
	"embed"
	"sync"
)

// 实际使用时，需要初始化该变量
var RepositoryFS *embed.FS
var TemplateDir = "template"

var repository *RepositorySQL
var repositoryOnce sync.Once

func GetRepositorySQL() *RepositorySQL {
	if repository == nil {
		InitRepositorySQL()
	}
	return repository
}

func InitRepositorySQL() {
	repositoryOnce.Do(func() {
		repository = NewRepositorySQL()
		err := repository.AddByFS(RepositoryFS, TemplateDir, TemplatefuncMap)
		if err != nil {
			panic(err)
		}
	})
}
