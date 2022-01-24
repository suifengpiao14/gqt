// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gqt is a template engine for SQL queries.

It helps to separate SQL code from Go code and permits to compose the queries
with a simple syntax.

The template engine is the standard package "text/template".

Usage

Create a template directory tree of .sql files. Here an example template with
the definition of three blocks:

	// File /path/to/sql/repository/dir/example.sql
	{{define "allUsers"}}
	SELECT *
	FROM users
	WHERE 1=1
	{{end}}

	{{define "getUser"}}
	SELECT *
	FROM users
	WHERE id=?
	{{end}}

	{{define "allPosts"}}
	SELECT *
	FROM posts
	WHERE date>=?
	{{if ne .Order ""}}ORDER BY date {{.Order}}{{end}}
	{{end}}

Then, with Go, add the directory to the default repository and execute the
queries:

	// Setup
	gqt.Add("/path/to/sql/repository/dir", "*.sql")

	// Simple query without parameters
	db.Query(gqt.Get("allUsers"))
	// Query with parameters
	db.QueryRow(gqt.Get("getuser"), 1)
	// Query with context and parameters
	db.Query(gqt.Exec("allPosts", map[string]interface{
		"Order": "DESC",
	}), date)

The templates are parsed immediately and recursively.

Namespaces

The templates can be organized in namespaces and stored in multiple root
directories.

	templates1/
	|-- roles/
	|	|-- queries.sql
	|-- users/
	|	|-- queries.sql
	|	|-- commands.sql

	templates2/
	|-- posts/
	|	|-- queries.sql
	|	|-- commands.sql
	|-- users/
	|	|-- queries.sql
	|-- queries.sql

The blocks inside the sql files are merged, the blocks with the same namespace
and name will be overridden following the alphabetical order.

The sub-directories are used as namespaces and accessed like:

	gqt.Add("../templates1", "*.sql")
	gqt.Add("../templates2", "*.sql")

	// Will search inside templates1/users/*.sql and templates2/users/*.sql
	gqt.Get("users/allUsers")

Multiple databases

When dealing with multiple databases at the same time, like PostgreSQL and
MySQL, just create two repositories:

	// Use a common directory
	dir := "/path/to/sql/repository/dir"

	// Create the PostgreSQL repository
	pgsql := gqt.NewRepository()
	pgsql.Add(dir, "*.pg.sql")

	// Create a separated MySQL repository
	mysql := gqt.NewRepository()
	mysql.Add(dir, "*.my.sql")

	// Then execute
	pgsql.Get("queryName")
	mysql.Get("queryName")
*/
package gqt

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

// Repository stores SQL templates.
type Repository struct {
	templates map[string]*template.Template // namespace: template
}

// NewRepository creates a new Repository.
func NewRepository() *Repository {
	return &Repository{
		templates: make(map[string]*template.Template),
	}
}

var suffix = ".sql.tpl"

// Add adds a root directory to the repository, recursively. Match only the
// given file extension. Blocks on the same namespace will be overridden. Does
// not follow symbolic links.
func (r *Repository) Add(root string, funcMap template.FuncMap) (err error) {
	// List the directories
	pattern := fmt.Sprintf("%s/*%s", strings.TrimRight(root, "/"), suffix)
	allFileList, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	for _, filename := range allFileList {
		relativeName := strings.TrimPrefix(filename, root)
		namespace := r.GetNamespace(relativeName)
		t, err := template.New(namespace).Funcs(funcMap).ParseFiles(filename)
		if err != nil {
			return err
		}
		r.templates[namespace] = t
	}

	return
}

func (r *Repository) GetNamespace(filename string) (namespace string) {
	namespace = strings.TrimSuffix(filename, suffix)
	namespace = strings.ReplaceAll(namespace, "\\", ".")
	namespace = strings.ReplaceAll(namespace, "/", ".")
	namespace = strings.Trim(namespace, ".")
	return
}

func (r *Repository) AddFromContent(filename string, content string, funcMap template.FuncMap) (err error) {
	namespace := r.GetNamespace(filename)
	t, err := template.New(namespace).Funcs(funcMap).Parse(content)
	if err != nil {
		return err
	}
	r.templates[namespace] = t
	return
}

// Get is a shortcut for r.Exec(), passing nil as data.
func (r *Repository) Get(name string) (s string, err error) {
	return r.Exec(name, nil)
}

// GetByNamespace get all template under namespace
func (r *Repository) GetByNamespace(namespace string) (s map[string]string, err error) {
	t, ok := r.templates[namespace]
	if !ok {
		err = fmt.Errorf("not found namespace:%s", namespace)
		return
	}
	s = make(map[string]string, 0)
	templates := t.Templates()
	for _, tpl := range templates {
		name := tpl.Name()
		var b bytes.Buffer
		err = tpl.Execute(&b, nil)
		if err != nil {
			return
		}
		fullname := fmt.Sprintf("%s.%s", namespace, name)
		content := strings.Trim(b.String(), "\r\n")
		if len(content) == 0 {
			continue
		}
		s[fullname] = b.String()
	}
	return
}

// Exec is a shortcut for r.Parse(), but panics if an error occur.
func (r *Repository) Exec(name string, data interface{}) (s string, err error) {
	s, err = r.Parse(name, data)
	return
}

// GetSql is a shortcut for r.Parse(), but panics if an error occur.
func (r *Repository) GetSql(name string, args interface{}, output *string) (err error) {
	if name == "" {
		err = errors.New("name not be empty")
		return err
	}
	*output, err = r.Parse(name, args)
	return
}

// Parse executes the template and returns the resulting SQL or an error.
func (r *Repository) Parse(name string, data interface{}) (string, error) {
	// Prepare namespace and block name
	if name == "" {
		return "", fmt.Errorf("unnamed block")
	}
	path := strings.Split(name, ".")
	namespace := strings.Join(path[0:len(path)-1], ".")
	if namespace == "." {
		namespace = ""
	}
	block := path[len(path)-1]

	// Execute the template
	var b bytes.Buffer
	t, ok := r.templates[namespace]
	if ok == false {
		return "", fmt.Errorf("unknown namespace \"%s\"", namespace)
	}
	err := t.ExecuteTemplate(&b, block, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

var g = singleflight.Group{}

func Flight(sql string, fn func() (interface{}, error)) (err error) {
	if sql == "" {
		err = errors.New("sql must not be empty")
		return
	}
	_, err, _ = g.Do(GetMD5LOWER(sql), fn)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

var defaultRepository = NewRepository()

// Add method for the default repository.
func Add(root string, funcMap template.FuncMap) error {
	return defaultRepository.Add(root, funcMap)
}

// Add method for the default repository.
func AddFromContent(filename string, content string, funcMap template.FuncMap) error {
	return defaultRepository.AddFromContent(filename, content, funcMap)
}

// Get method for the default repository.
func Get(name string) (s string, err error) {
	return defaultRepository.Get(name)
}

// Exec method for the default repository.
func Exec(name string, data interface{}) (s string, e error) {
	return defaultRepository.Exec(name, data)
}

// Parse method for the default repository.
func Parse(name string, data interface{}) (string, error) {
	return defaultRepository.Parse(name, data)
}

// ReadEmbedFS read embed file
func ReadEmbedFS(repositoryFS embed.FS, filename string, fileMap *map[string][]byte) {
	filename = strings.TrimRight(filename, "/")
	if len(filename) >= 2 {
		firstTwoLetter := filename[0:2]
		if firstTwoLetter == "./" { // 切除./ 开头的路径
			filename = filename[2:]
		}
	}
	fsFile, err := repositoryFS.Open(filename)
	if err != nil {
		panic(err)
	}
	fsInfo, err := fsFile.Stat()
	if err != nil {
		panic(err)
	}
	if fsInfo.IsDir() {
		fsList, err := repositoryFS.ReadDir(filename)
		if err != nil {
			panic(err)
		}
		for _, fileInfo := range fsList {
			subFilename := fmt.Sprintf("%s/%s", filename, fileInfo.Name())
			if fileInfo.IsDir() {

				ReadEmbedFS(repositoryFS, subFilename, fileMap)
				continue
			}
			b, err := repositoryFS.ReadFile(subFilename)
			if err != nil {
				panic(err)
			}
			(*fileMap)[subFilename] = b
		}
		return
	}
}
