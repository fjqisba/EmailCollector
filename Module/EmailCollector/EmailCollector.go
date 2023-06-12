package EmailCollector

import (
	"EmailCollector/Module/EmailMiner"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"github.com/panjf2000/ants/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type EmailCollector struct {
	//输出文件
	hCsvOut *csv.Writer
	fileMutex sync.Mutex
	//输入网址列表
	txtFilePath string
	//线程个数
	threadCount int
	//线程池
	hPool *ants.Pool
	wg sync.WaitGroup
}

func NewEmailCollector(hOutCsv *csv.Writer,filepath string,thCount int)*EmailCollector  {
	return &EmailCollector{
		hCsvOut: hOutCsv,
		txtFilePath: filepath,
		threadCount: thCount,
	}
}

func (this *EmailCollector)queryTaskWrapper(id int,webUrl string)func()  {
	return func() {
		var emailMiner EmailMiner.EmailMiner
		emailList := emailMiner.DetectEmail(webUrl)
		strEmail,_ := json.Marshal(emailList)
		this.fileMutex.Lock()
		if len(emailList) == 0{
			this.hCsvOut.Write([]string{strconv.Itoa(id),webUrl,""})
		}else{
			log.Println(webUrl,"检测出邮箱:",string(strEmail))
			this.hCsvOut.Write([]string{strconv.Itoa(id),webUrl,string(strEmail)})
		}
		this.hCsvOut.Flush()
		this.fileMutex.Unlock()
		this.wg.Done()
	}
}

func (this *EmailCollector)StartWork()error  {
	hFile,err := os.Open(this.txtFilePath)
	if err != nil{
		return err
	}
	defer hFile.Close()
	this.hPool,err = ants.NewPool(this.threadCount)
	if err != nil{
		return err
	}
	defer this.hPool.Release()
	//开始循环解析
	hScanner := bufio.NewScanner(hFile)
	index := 0
	for hScanner.Scan() {
		tmpIndex := index
		index = index + 1
		line := hScanner.Text()
		line = strings.TrimSpace(line)
		this.wg.Add(1)
		this.hPool.Submit(this.queryTaskWrapper(tmpIndex,line))
	}
	this.wg.Wait()
	return nil
}