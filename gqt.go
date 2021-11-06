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
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
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
		t, err := template.New("").Funcs(funcMap).ParseFiles(filename)
		if err != nil {
			return err
		}
		relativeName := strings.TrimPrefix(filename, root)
		namespace := r.GetNamespace(relativeName)
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

// Add adds a root directory to the repository, recursively. Match only the
// given file extension. Blocks on the same namespace will be overridden. Does
// not follow symbolic links.
func (r *Repository) AddFromPackrBox(box *packr.Box, funcMap template.FuncMap) (err error) {
	// List the directories
	allFileList := box.List()

	var content string
	var filename string
	for _, filename = range allFileList {
		if !strings.Contains(filename, suffix) {
			continue
		}
		content, err = box.FindString(filename)
		if err != nil {
			return
		}
		t, err := template.New("").Funcs(funcMap).Parse(content)
		if err != nil {
			return err
		}
		namespace := r.GetNamespace(filename)
		r.templates[namespace] = t
	}
	return
}

func (r *Repository) AddFromContent(filename string, content string, funcMap template.FuncMap) (err error) {
	t, err := template.New("").Funcs(funcMap).Parse(content)
	if err != nil {
		return err
	}
	namespace := r.GetNamespace(filename)
	r.templates[namespace] = t
	return
}

// addDir parses a directory.
func (r *Repository) addDir(path, namespace, pattern string, funcMap template.FuncMap) error {
	// Parse the template
	t, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(path, pattern))
	if err != nil {
		return err
	}
	r.templates[namespace] = t
	return nil
}

// Get is a shortcut for r.Exec(), passing nil as data.
func (r *Repository) Get(name string) (s string, err error) {
	return r.Exec(name, nil)
}

// Exec is a shortcut for r.Parse(), but panics if an error occur.
func (r *Repository) Exec(name string, data interface{}) (s string, err error) {
	s, err = r.Parse(name, data)
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

var defaultRepository = NewRepository()

// Add method for the default repository.
func Add(root string, funcMap template.FuncMap) error {
	return defaultRepository.Add(root, funcMap)
}

// Add method for the default repository.
func AddFromContent(filename string, content string, funcMap template.FuncMap) error {
	return defaultRepository.AddFromContent(filename, content, funcMap)
}

// Add method for the default repository.
func AddFromPackrBox(box *packr.Box, funcMap template.FuncMap) error {
	return defaultRepository.AddFromPackrBox(box, funcMap)
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
