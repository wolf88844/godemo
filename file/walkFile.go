package main

import (
	"bufio"
	"godemo/input"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeTemplate      = "2006-01-02 15:04:05"
	currentTimeTemple = "2006-01-02"
)

var (
	targetDirName, targetPath, srcPath, currentPath, javaName, classPath, name, classpath string    //编译后目录名称 编译后路径 原路径 当前路径 java文件名称 class文件名称
	checkTime                                                                             time.Time //判断时间
	count, dircount                                                                       int16     //文件数量 文件夹数量
	namesPath, names                                                                      map[string][]string
)

func main() {
	srcPath = input.InputSrcPath()
	log.Printf("原路径为：%s\n", srcPath)
	log.Println("=======")

	log.Println("请输入编译后路径：")
	targetInput := bufio.NewScanner(os.Stdin)
	targetInput.Scan()
	targetPath = targetInput.Text()
	log.Printf("编译后路径为：%s\n", targetPath)
	if targetPath != "" {
		//在编译后的路径里查找有一个java对应多个class文件的名称
		namesPath = FindMultiClassFile(targetPath)
	}
	log.Println("=======")

	log.Println("请输入复制后的文件夹名称：")
	targetDirInput := bufio.NewScanner(os.Stdin)
	targetDirInput.Scan()
	targetDirName = targetDirInput.Text()
	if targetDirName == "" {
		targetDirName = "updateFiles"
	}
	log.Println("=======")

	log.Println("请输入比较时间：")
	timeInput := bufio.NewScanner(os.Stdin)
	timeInput.Scan()
	text := timeInput.Text()
	if text == "" {
		now := time.Now()
		format := now.Format(currentTimeTemple)
		stamp, _ := time.ParseInLocation(currentTimeTemple, format, time.Local)
		checkTime = stamp.Local()
	} else {
		stamp, _ := time.ParseInLocation(timeTemplate, text, time.Local)
		checkTime = stamp.Local()
	}
	log.Printf("截止时间为：%v\n", checkTime)
	log.Println("=======")

	//获取当前路径
	currentPath, _ = os.Getwd()
	log.Printf("当前路径%s\n", currentPath)
	log.Println("=======")

	//遍历文件夹
	filepath.Walk(srcPath, walkFunc)

	log.Printf("共有 %d个文件夹, %d 个文件需要更新\n", dircount, count)
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if info.ModTime().After(checkTime) {
		//计数
		if info.IsDir() {
			dircount++
		} else {
			count++
		}

		//判断是否是java文件
		name := info.Name()
		if strings.Contains(name, ".java") && targetPath != "" {
			javaName = path[strings.LastIndex(path, "\\")+1 : strings.LastIndex(path, ".java")]
			//先找是否属于一对多的情况（一个java文件对多个class文件）
			if paths, ok := namesPath[javaName]; ok {
				for _, value := range paths {
					classPath = value
					copyClassFile()
				}
			} else {
				//找到java类对应的class文件地址
				filepath.Walk(targetPath, wolfTargetFunc)
				copyClassFile()
			}
		}
		//处理原文件
		if !strings.EqualFold(path, srcPath) && !info.IsDir() {
			length := len(srcPath)
			newPath := currentPath + "\\" + targetDirName + "\\" + path[length+1:]
			copyFile(newPath, path)
			log.Println(newPath)
		}

	}
	return nil
}

func copyClassFile() {
	s := classPath[len(targetPath)+1:]
	newPath := currentPath + "\\" + targetDirName + "\\" + s
	//处理class
	copyFile(newPath, classPath)
}

/**
查询一对一文件格式的遍历方法
*/
func wolfTargetFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.EqualFold(info.Name(), javaName+".class") {
		classPath = path
		return nil
	}
	return nil
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

	if ok, _ := PathExists(dstFileName); !ok {
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
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

//获取一个java文件对应多个class文件的数据
func FindMultiClassFile(targetPath string) map[string][]string {
	names = make(map[string][]string)
	filepath.Walk(targetPath, walkMultiFunc)
	return names
}

/**
遍历文件，查询多个class文件，将文件名与路径放入map中
*/
func walkMultiFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if strings.Contains(path, "$") {
		index := strings.LastIndex(path, "$") - 1
		lastIndex := strings.LastIndex(path, "\\") + 1
		name := path[lastIndex : index+1]
		dirPath := path[:lastIndex-1]
		if _, ok := names[name]; !ok {
			log.Printf("有多个class文件的java文件名称为：%s,文件夹路径为： %s\n", name, dirPath)
			//查找对应所有的class文件
			names[name] = getFileList(dirPath, name)
		}

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
