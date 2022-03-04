package gqt

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

func FileName2Namespace(filename string, dir string, suffix string) (namespace string) {
	prefix := strings.ReplaceAll(dir, "\\", ".")
	prefix = strings.ReplaceAll(prefix, "/", ".")
	namespace = strings.TrimSuffix(filename, suffix)
	namespace = strings.ReplaceAll(namespace, "\\", ".")
	namespace = strings.ReplaceAll(namespace, "/", ".")
	namespace = strings.TrimPrefix(namespace, prefix)
	namespace = strings.Trim(namespace, ".")
	return
}

func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetTplFilesByDir get current and reverse dir tpl file
func GetTplFilesByDir(dir string) (allFileList []string, err error) {
	pattern := fmt.Sprintf("%s/*%s", strings.TrimRight(dir, "/"), Suffix)
	allFileList, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	pattern = fmt.Sprintf("%s/**/*%s", strings.TrimRight(dir, "/"), Suffix)
	subDirFileList, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	allFileList = append(allFileList, subDirFileList...)
	return
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
