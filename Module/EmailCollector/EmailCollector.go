package EmailCollector

import (
	"EmailCollector/Module/EmailMiner"
	"EmailCollector/Module/PhantomJS"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type EmailCollector struct {
	hCsvOut *csv.Writer
	txtFilePath string
	threadCount int
	hPool *ants.Pool
}

func NewEmailCollector(hOutCsv *csv.Writer,filepath string,thCount int)*EmailCollector  {
	return &EmailCollector{
		hCsvOut: hOutCsv,
		txtFilePath: filepath,
		threadCount: thCount,
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
		line := hScanner.Text()
		line = strings.TrimSpace(line)
		if line == ""{
			continue
		}
		tmpIndex := index
		index = index + 1
		this.hPool.Submit(func() {


			fmt.Println("爬取网址:",tmpIndex,line)
			EmailMiner.DetectEmail(hPool,tmpIndex,line)
			pageContent := PhantomJS.GetPageHtml(line)
			emailList := GuessEmail(pageContent)
			strEmail, _ := json.Marshal(emailList)
			tmpMutex.Lock()
			if len(emailList) == 0{
				hCsvOut.Write([]string{strconv.Itoa(tmpIndex),line,""})
			}else{
				hCsvOut.Write([]string{strconv.Itoa(tmpIndex),line,string(strEmail)})
			}
			hCsvOut.Flush()
			tmpMutex.Unlock()
			time.Sleep(10 * time.Second)
		})
	}
	return nil
}