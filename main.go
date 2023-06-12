package main

import (
	"EmailCollector/Module/EmailCollector"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var(
	tmpMutex sync.Mutex
)

func main() {
	var filepath string
	var threadCount int
	flag.StringVar(&filepath,"f","","网址文件")
	flag.IntVar(&threadCount,"t",0,"线程个数")
	flag.Parse()
	if filepath == "" {
		fmt.Println("用法: EmailCollector.exe -f 文件路径 -t 线程个数")
		return
	}
	if threadCount == 0{
		threadCount = 3
	}
	hFile,err := os.Create("./" + time.Now().Format("200601021504") + ".csv")
	if err != nil{
		fmt.Println("创建输出文件失败")
		return
	}
	defer hFile.Close()
	hCsvOut := csv.NewWriter(hFile)
	coreCollector := EmailCollector.NewEmailCollector(hCsvOut,filepath,threadCount)
	err = coreCollector.StartWork()
	if err != nil{
		fmt.Println(err)
	}
}
