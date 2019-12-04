package file

import (
	"godemo/variable"
	"godemo/walkfunc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//获取一个java文件对应多个class文件的数据
func FindMultiClassFile(targetPath string) map[string][]string {
	var names = make(map[string][]string)
	filepath.Walk(targetPath, walkfunc.WalkMultiClassFunc)
	return names
}

//遍历原地址文件
func ErgodicFiles(srcPath string) {
	filepath.Walk(srcPath, walkfunc.WalkFilesFunc)
}

/**
递归查找包含有name的文件名称
*/
func GetFileList(path string, name string) []string {
	var paths []string
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if file.IsDir() {
			GetFileList(path+file.Name()+"\\", name)
		} else {
			if strings.HasPrefix(file.Name(), name) {
				paths = append(paths, path+"\\"+file.Name())
			}
		}
	}
	return paths
}

//复制class文件
func CopyClassFile() {
	s := walkfunc.ClassPath[len(walkfunc.TargetPath)+1:]
	newPath := variable.TargetDirPath + "\\" + s
	//处理class
	CopyFile(newPath, walkfunc.ClassPath)
}

/**
执行文件复制操作
*/
func CopyFile(dstFileName string, srcFileName string) (err error) {
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
