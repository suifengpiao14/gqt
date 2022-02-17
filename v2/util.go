package main

import (
	"embed"
	"fmt"
	"strings"
	"time"
)

// : 在sqlx.Named方法中 为标记为变量开始的特殊字符，写入时间字符会有冲突(使用:: 和常规查看不符，所以时间不内置到模板函数中)
func ZeroTime() string {
	return "0000-00-00 00:00:00"
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func PermanentTime() string {
	return "3000-12-31 23:59:59"
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
