package main

import (
	"EmailCollector/Module/EmailCollector"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var(
	tmpMutex sync.Mutex
	regex_Email = regexp.MustCompile(`[a-zA-Z0-9._%+-]+(@|\\u0040)[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}`)
)

func filterEmail(emailAddr string)bool  {
	str := strings.ToLower(emailAddr)
	if strings.HasSuffix(str,".jpg") == true{
		return true
	}
	if strings.HasSuffix(str,".gif") == true{
		return true
	}
	if strings.Index(str,"@mail.com") != -1{
		return true
	}
	if strings.Index(str,"example") != -1{
		return true
	}
	if strings.Index(str,".wix") != -1{
		return true
	}
	if strings.Index(str,".png") != -1{
		return true
	}
	if strings.Index(str,"sentry.io") != -1{
		return true
	}else if strings.Contains(str,".webp"){
		return true
	}
	return false
}

func GuessEmail(content string)(retList []string)  {
	if content == ""{
		return retList
	}
	filterMap := make(map[string]struct{})
	matchList := regex_Email.FindAllStringSubmatch(content,-1)
	if len(matchList) == 0{
		return retList
	}
	for _,eMatch := range matchList{
		emailAddr := eMatch[0]
		if filterEmail(emailAddr) == true{
			continue
		}
		if _,bExists := filterMap[eMatch[0]];bExists == false{
			filterMap[eMatch[0]] = struct{}{}
			retList = append(retList, eMatch[0])
		}
	}
	return retList
}

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
	StartQuery(hCsvOut,filepath,threadCount)
}
