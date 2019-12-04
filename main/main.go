package main

import (
	"godemo/file"
	"godemo/input"
	"godemo/variable"
	"godemo/walkfunc"
	"log"
)

func main() {
	//输入原地址
	variable.SrcPath = input.InputSrcPath()
	//输入复制路径
	variable.TargetDirPath = input.InputOutFileDirPath()
	//输入编译后地址，并查询一对多文件到map中
	walkfunc.TargetPath, walkfunc.Names = input.InputTargetPath()
	//输入判断时间
	walkfunc.CheckTime = input.InputCheckTime()
	//遍历原地址文件
	file.ErgodicFiles(variable.SrcPath)
	log.Printf("共有 %d 个文件需要更新\n", walkfunc.Count)
}