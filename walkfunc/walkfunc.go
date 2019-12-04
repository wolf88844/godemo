package walkfunc

import (
	"godemo/file"
	"godemo/variable"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Names                           map[string][]string
	CheckTime                       time.Time
	Count                           int16
	TargetPath, ClassPath, JavaName string
)

/**
遍历class文件，查询多个class文件，将文件名与路径放入map中
*/
func WalkMultiClassFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if strings.Contains(path, "$") {
		index := strings.LastIndex(path, "$") - 1
		lastIndex := strings.LastIndex(path, "\\") + 1
		name := path[lastIndex : index+1]
		dirPath := path[:lastIndex-1]
		if _, ok := Names[name]; !ok {
			log.Printf("有多个class文件的java文件名称为：%s,文件夹路径为： %s\n", name, dirPath)
			//查找对应所有的class文件
			Names[name] = file.GetFileList(dirPath, name)
		}

	}
	return nil
}

func WalkFilesFunc(path string, info os.FileInfo, err error) error {
	if info.ModTime().After(CheckTime) {
		//计数
		Count++
		//判断是否是java文件
		name := info.Name()
		if strings.Contains(name, ".java") && TargetPath != "" {
			JavaName = path[strings.LastIndex(path, "\\")+1 : strings.LastIndex(path, ".java")]
			//先找是否属于一对多的情况（一个java文件对多个class文件）
			if paths, ok := Names[JavaName]; ok {
				for _, value := range paths {
					ClassPath = value
					file.CopyClassFile()
				}
			} else {
				//找到java类对应的class文件地址
				filepath.Walk(TargetPath, wolkTargetFunc)
				file.CopyClassFile()
			}
		}
		//处理原文件
		if !strings.EqualFold(path, variable.SrcPath) && !info.IsDir() {
			length := len(variable.SrcPath)
			newPath := variable.TargetDirPath + "\\" + path[length+1:]
			file.CopyFile(newPath, path)
			log.Println(newPath)
		}

	}
	return nil
}

/**
查询一对一文件格式的遍历方法
*/
func wolkTargetFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.EqualFold(info.Name(), JavaName+".class") {
		ClassPath = path
		return nil
	}
	return nil
}
