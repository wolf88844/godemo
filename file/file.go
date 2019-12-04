package file

import (
	"godemo/variable"
	"io"
	"io/ioutil"
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
	TargetPath, classPath, javaName string
)

//获取一个java文件对应多个class文件的数据
func FindMultiClassFile(targetPath string) map[string][]string {
	var names = make(map[string][]string)
	filepath.Walk(targetPath, walkMultiClassFunc)
	return names
}

/**
遍历class文件，查询多个class文件，将文件名与路径放入map中
*/
func walkMultiClassFunc(path string, info os.FileInfo, err error) error {
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
			Names[name] = getFileList(dirPath, name)
		}

	}
	return nil
}

//遍历原地址文件
func ErgodicFiles(srcPath string) {
	filepath.Walk(srcPath, walkFilesFunc)
}

func walkFilesFunc(path string, info os.FileInfo, err error) error {
	if info.ModTime().After(CheckTime) {
		//计数
		Count++
		//判断是否是java文件
		name := info.Name()
		if strings.Contains(name, ".java") && TargetPath != "" {
			javaName = path[strings.LastIndex(path, "\\")+1 : strings.LastIndex(path, ".java")]
			//先找是否属于一对多的情况（一个java文件对多个class文件）
			if paths, ok := Names[javaName]; ok {
				for _, value := range paths {
					classPath = value
					copyClassFile()
				}
			} else {
				//找到java类对应的class文件地址
				filepath.Walk(TargetPath, wolkTargetFunc)
				copyClassFile()
			}
		}
		//处理原文件
		if !strings.EqualFold(path, variable.SrcPath) && !info.IsDir() {
			length := len(variable.SrcPath)
			newPath := variable.TargetDirPath + "\\" + path[length+1:]
			copyFile(newPath, path)
			log.Println(newPath)
		}

	}
	return nil
}

/**
查询一对一文件格式的遍历方法
*/
func wolkTargetFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.EqualFold(info.Name(), javaName+".class") {
		classPath = path
		return nil
	}
	return nil
}

/**
递归查找包含有name的文件名称
*/
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

//复制class文件
func copyClassFile() {
	s := classPath[len(TargetPath)+1:]
	newPath := variable.TargetDirPath + "\\" + s
	//处理class
	copyFile(newPath, classPath)
}

/**
执行文件复制操作
*/
func copyFile(dstFileName string, srcFileName string) (err error) {
	srcFile, err := os.OpenFile(srcFileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("open srcfile err=%v\n", err)
		return
	}
	defer srcFile.Close()

	//判断如果是.java文件，则去获取对应的.class文件

	if ok, _ := pathExists(dstFileName); !ok {
		s := dstFileName[:strings.LastIndex(dstFileName, "\\")]
		os.MkdirAll(s, os.ModePerm)
	}

	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("open dstfile err= %v\n", err)
		return
	}
	defer dstFile.Close()

	buf := make([]byte, 1024*4)
	for {
		n, err2 := srcFile.Read(buf)
		if err2 != nil && err2 != io.EOF {
			log.Fatalln(err2)
		}
		if n == 0 {
			break
		}
		if _, err := dstFile.Write(buf[:n]); err != nil {
			log.Fatalln(err2)
		}
	}
	return nil
}

/**
判断文件是否存在
*/
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
