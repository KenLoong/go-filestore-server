package test

import (
	"filestore-server/util"
	"fmt"
	"os"
	"testing"
)

func TestSth(t *testing.T) {
	t.Run("测试文件sha1值",testSha1)
}

func testSha1(t *testing.T) {
	file, err := os.Open("F:\\GoProjects\\src\\filestore-server\\tmp\\操作系统.txt")
	if err != nil {
		fmt.Println("read file err : ",err.Error())
		return
	}
	defer file.Close()
	fmt.Println(util.FileSha1(file)) //716312f237f7d6580ebbd1e216a09f3f33ae5845
}


