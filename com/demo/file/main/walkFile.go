package main

import (
	"bufio"
	"godemo/com/demo/file"
	"io"
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
	targetDirName, targetPath, srcPath, currentPath, javaName, classPath string
	checkTime                                                            time.Time
	count, dircount                                                      int16
	names                                                                map[string][]string
)

//func main(){
//	log.Println("请输入编译后路径：")
//		targetInput := bufio.NewScanner(os.Stdin)
//		targetInput.Scan()
//		targetPath = targetInput.Text()
//		log.Printf("编译后路径为：%s\n", targetPath)
//		//在编译后的路径里查找有一个java对应多个class文件的名称
//		result := file.FindMultiClassFile(targetPath)
//		for k,v:=range result{
//			fmt.Printf("%s %s\n",k,v)
//		}
//}

func main() {
	log.Println("请输入文件路径：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	srcPath = input.Text()
	log.Printf("原路径为：%s\n", srcPath)
	log.Println("=======")
	log.Println("请输入编译后路径：")
	targetInput := bufio.NewScanner(os.Stdin)
	targetInput.Scan()
	targetPath = targetInput.Text()
	log.Printf("编译后路径为：%s\n", targetPath)
	if targetPath != "" {
		//在编译后的路径里查找有一个java对应多个class文件的名称
		names = file.FindMultiClassFile(targetPath)
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
			if paths, ok := names[javaName]; ok {
				for _, value := range paths {
					classPath = value
					s := classPath[len(targetPath)+1:]
					newPath := currentPath + "\\" + targetDirName + "\\" + s
					//处理class
					copyFile(newPath, classPath)
				}
			} else {
				//找到java类对应的class文件地址
				filepath.Walk(targetPath, wolfTargetFunc)
				s := classPath[len(targetPath)+1:]
				newPath := currentPath + "\\" + targetDirName + "\\" + s
				//处理class
				copyFile(newPath, classPath)
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

func wolfTargetFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.EqualFold(info.Name(), javaName+".class") {
		classPath = path
		return nil
	}
	return nil
}

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
