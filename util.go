package gqt

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"strings"
)

func FileName2Namespace(filename string, dir string, suffix string) (namespace string) {
	prefix := strings.ReplaceAll(dir, "\\", ".")
	prefix = strings.ReplaceAll(prefix, "/", ".")
	namespace = strings.TrimSuffix(filename, suffix)
	namespace = strings.ReplaceAll(namespace, "\\", ".")
	namespace = strings.ReplaceAll(namespace, "/", ".")
	namespace = strings.Trim(namespace, ".")
	namespace = strings.TrimPrefix(namespace, prefix)
	return
}

func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
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
