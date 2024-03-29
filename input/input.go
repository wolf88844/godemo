package input

import (
	"bufio"
	"godemo/file"
	"log"
	"os"
	"time"
)

const (
	timeTemplate      = "2006-01-02 15:04:05"
	currentTimeTemple = "2006-01-02"
)

//输入原地址
func InputSrcPath() string {
	log.Println("请输入文件路径：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	srcPath := input.Text()
	log.Printf("原路径为：%s\n", srcPath)
	log.Println("=======")
	return srcPath
}

//输入编译后地址
func InputTargetPath() (string, map[string][]string) {
	log.Println("请输入编译后路径：")
	targetInput := bufio.NewScanner(os.Stdin)
	targetInput.Scan()
	targetPath := targetInput.Text()
	log.Printf("编译后路径为：%s\n", targetPath)
	var namesPath = make(map[string][]string)
	if targetPath != "" {
		//在编译后的路径里查找有一个java对应多个class文件的名称
		namesPath = file.FindMultiClassFile(targetPath)
	}
	log.Println("=======")
	return targetPath, namesPath
}

//输入导出文件名称
func InputOutFileDirPath() string {
	log.Println("请输入复制后的路径：")
	targetDirInput := bufio.NewScanner(os.Stdin)
	targetDirInput.Scan()
	var targetDirName = targetDirInput.Text()
	if targetDirName == "" {
		//获取当前路径
		currentPath, _ := os.Getwd()
		targetDirName = currentPath + "\\updateFiles"
	}
	log.Printf("复制路径为：%s\n", targetDirName)
	log.Println("=======")
	return targetDirName
}

//输入开始时间
func InputBeginTime() time.Time {
	log.Println("请输入开始时间：")
	timeInput := bufio.NewScanner(os.Stdin)
	timeInput.Scan()
	text := timeInput.Text()
	if text == "" {
		now := time.Now()
		text = now.Format(currentTimeTemple)
	}
	stamp, _ := time.ParseInLocation(timeTemplate, text, time.Local)
	beginTime := stamp.Local()
	log.Printf("开始时间为：%v\n", beginTime)
	log.Println("=======")
	return beginTime
}

//输入结束时间
func InputEndTime() time.Time {
	log.Println("请输入结束时间：")
	timeInput := bufio.NewScanner(os.Stdin)
	timeInput.Scan()
	text := timeInput.Text()
	if text == "" {
		now := time.Now()
		text = now.Format(currentTimeTemple)
	}
	stamp, _ := time.ParseInLocation(timeTemplate, text, time.Local)
	beginTime := stamp.Local()
	log.Printf("结束时间为：%v\n", beginTime)
	log.Println("=======")
	return beginTime
}
