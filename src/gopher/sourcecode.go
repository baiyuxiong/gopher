package gopher

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"bytes"

	"github.com/gorilla/mux"
)

// 列出所有版本源码
func sourcecodeHandler(handler Handler) {
	renderTemplate(handler, "sourcecode/index.html", BASE, map[string]interface{}{
		"action": "/sourcecode/index",
		"active": "sourcecode",
	})
}

//指定版本的源码
func sourceVersionHandler(handler Handler){
	sourcecodeRoot := Config.SourcecodeRoot

	vars := mux.Vars(handler.Request)
	version := vars["version"]
	path := vars["path"]

	fileFullPath := sourcecodeRoot + "/" + version
	parentPath := ""
	if(len(path) > 0){
		if(strings.LastIndexAny(path, "-") > 0){
			parentPath = "/"+path[0:strings.LastIndexAny(path, "-")];
		}
		newPath := strings.Replace(path,"-","/",-1);
		fileFullPath = fileFullPath + "/" +newPath
	}

	fmt.Println(fileFullPath);

	fi, err := os.Stat(fileFullPath)

	if err != nil {
		panic("open file "+ fileFullPath +"failed!")
	}

	if fi.IsDir() {
		//List dir

//		fileList := make([]os.FileInfo, 0)
		//
		//		filepath.Walk(fileFullPath,func(path string, info os.FileInfo, err error) error {
		//				fileList = append(fileList,info)
		//				return nil
		//			})

		rd, err := ioutil.ReadDir(fileFullPath)
		if err != nil {
			panic("read dir"+ fileFullPath +"failed!")
		}

		renderTemplate(handler, "sourcecode/show.html", BASE, map[string]interface{}{
			"fileList" : rd,
			"action": "/sourcecode/show",
			"active": "sourcecode",
			"isDir":true,
			"version":version,
			"path":path,
			"parentPath":parentPath,
		})
	}else{
		buff, err := ioutil.ReadFile(fileFullPath)
		if err != nil {
			panic("open file failed!")
		}
		buff = bytes.Trim(buff, "\xef\xbb\xbf") //strip BOM of utf-8 file
		renderTemplate(handler, "sourcecode/show.html", BASE, map[string]interface{}{
			"fileContent" : string(buff),
			"action": "/sourcecode/show",
			"active": "sourcecode",
			"isDir":false,
			"version":version,
			"path":path,
			"parentPath":parentPath,
		})
	}
}
