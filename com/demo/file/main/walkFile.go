package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var checkTime time.Time
const timeTemplate ="2006-01-02 15:04:05"
const currentTimeTemple  = "2006-01-02"
var count,dircount int16
var srcPath,currentPath string


func main() {
	log.Println("请输入文件路径：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	srcPath = input.Text()
	log.Printf("原路径为：%s\n", srcPath)
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
	log.Printf("截止时间为：%v\n",checkTime)
	//获取当前路径
	currentPath,_ =os.Getwd()
	log.Printf("当前路径%s\n",currentPath)
	//遍历文件夹
	filepath.Walk(srcPath,walkFunc)

	log.Printf("共有 %d个文件夹, %d 个文件需要更新\n",dircount,count)
}

func walkFunc(path string, info os.FileInfo,err error) error{
	if info.ModTime().After(checkTime){
		//计数
		if info.IsDir(){
			dircount++
		}else{
				count++
		}
		//处理
		if !strings.EqualFold(path,srcPath) && !info.IsDir(){
			length := len(srcPath)
			newPath := currentPath+"\\updateFiles"+"\\"+path[length+1:]
			copyFile(newPath,path)
			log.Println(newPath)
		}
	}

	return nil
}

func copyFile(dstFileName string,srcFileName string)(err error){
	srcFile,err := os.OpenFile(srcFileName,os.O_RDONLY,0666)
	if err != nil{
		log.Fatalf("open srcfile err=%v\n",err)
		return
	}
	defer srcFile.Close()

	if ok,_:=PathExists(dstFileName);!ok{
		s := dstFileName[:strings.LastIndex(dstFileName, "\\")]
		os.MkdirAll(s, os.ModePerm)
	}

	dstFile,err:= os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil{
		log.Fatalf("open dstfile err= %v\n",err)
		return
	}
	defer dstFile.Close()

	buf :=make([]byte,1024*4)
	for{
		n,err2:=srcFile.Read(buf)
		if err2 != nil && err2 !=io.EOF{
			log.Fatalln(err2)
		}
		if n==0{
			break
		}
		if _,err:=dstFile.Write(buf[:n]);err!=nil{
			log.Fatalln(err2)
		}
	}
	return nil
}

func PathExists(path string)(bool ,error){
	_,err:=os.Stat(path)
	if err ==nil{
		return true,nil
	}
	if os.IsNotExist(err){
		return false,nil
	}
	return false,nil
}
