package input

import (
	"bufio"
	"log"
	"os"
)

func InputSrcPath() string {
	log.Println("请输入文件路径：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	srcPath := input.Text()
	return srcPath
}
