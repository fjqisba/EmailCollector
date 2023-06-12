package TestUnit

import (
	"EmailCollector/Module/EmailMiner"
	"log"
	"testing"
)

func CheckEmail(t *testing.T,content string,result string)  {
	emailList := EmailMiner.GuessEmail(content)
	if len(emailList) == 0{
		log.Println("检测邮箱失败:",result)
		return
	}
	if emailList[0] != result{
		log.Println("检测邮箱失败:",result)
	}
}

func TestRegex(t *testing.T) {
	CheckEmail(t,"onmousemove=\"xr_mo(this,0)\" >office&#64;computer-corner.at</a","office@computer-corner.at")
	CheckEmail(t,"aggregated_ranges\":[],\"ranges\":[],\"color_ranges\":[],\"text\":\"thefonerepairs\\u0040gmail.com\"}},\"__module_o","thefonerepairs@gmail.com")
}
