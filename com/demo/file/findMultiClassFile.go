package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	name, classpath string
	names           map[string][]string
)

//获取一个java文件对应多个class文件的数据
func FindMultiClassFile(targetPath string) map[string][]string {
	names = make(map[string][]string)
	filepath.Walk(targetPath, walkFunc)
	return names
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if strings.Contains(path, "$") {
		index := strings.LastIndex(path, "$") - 1
		lastIndex := strings.LastIndex(path, "\\") + 1
		name := path[lastIndex : index+1]
		dirPath := path[:lastIndex-1]
		if _, ok := names[name]; !ok {
			fmt.Printf("有多个class文件的java文件名称为：%s,文件夹路径为： %s\n", name, dirPath)
			//查找对应所有的class文件
			names[name] = getFileList(dirPath, name)
		}

	}
	return nil
}

func getFileList(path string, name string) []string {
	var paths []string
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if file.IsDir() {
			getFileList(path+file.Name()+"\\", name)
		} else {
			if strings.HasPrefix(file.Name(), name) {
				paths = append(paths, path+"\\"+file.Name())
			}
		}
	}
	return paths
}
