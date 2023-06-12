package TestUnit

import (
	"EmailCollector/Module/EmailMiner"
	"log"
	"testing"
)

func CheckUrl(webUrl string,result string)  {
	var emailEngine EmailMiner.EmailMiner
	emailList := emailEngine.DetectEmail(webUrl)
	if len(emailList) == 0{
		log.Println("检测邮箱错误:",webUrl)
		return
	}
	if emailList[0] != result{
		log.Println("检测邮箱错误:",webUrl)
		return
	}
}

func TestUrl(t *testing.T) {

	CheckUrl("https://airconnect.at/","office@airphone.at")

}