package EmailMiner

import (
	"EmailCollector/Const"
	"EmailCollector/Module/ContactDetector"
	"EmailCollector/Module/PhantomJS"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

var(
	suffixList = []string{".webp",".gif",".htm",".html",".jpeg",".jpg",".php",".png", ".ace",".ani",".arc",".arj",".avi",".bmp",".cab",".class",".css",".exe",".ico",".jar",".mid",".mov",".mp2",".mp3",".mpeg",".mpg",".pdf",".wix"}
	regex_Email = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}`)
)

type EmailMiner struct {
	mutex sync.Mutex
	emailList []string
	emailFilterMap map[string]struct{}
}

func filterEmail(emailAddr string)bool  {
	str := strings.ToLower(emailAddr)
	for _,eSuffix := range suffixList{
		if strings.HasSuffix(str,eSuffix) == true{
			return true
		}
	}
	if strings.Index(str,"@mail.com") != -1{
		return true
	}
	if strings.Index(str,"example") != -1{
		return true
	}
	if strings.Index(str,"sentry.io") != -1{
		return true
	}
	return false
}

func GuessEmail(content string)(retList []string)  {
	if content == ""{
		return retList
	}
	content = strings.ReplaceAll(content,"&#64;","@")
	content = strings.ReplaceAll(content,"\\u0040","@")
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

func (this *EmailMiner)exploreUrl(webUrl string)string  {
	log.Println("检测网址:",webUrl)
	pageContent := PhantomJS.GetPageHtml(webUrl)
	emailList := GuessEmail(pageContent)
	this.mutex.Lock()
	for _,eMail := range emailList{
		if _,bExists := this.emailFilterMap[eMail];bExists == false{
			this.emailList = append(this.emailList, eMail)
			this.emailFilterMap[eMail] = struct{}{}
		}
	}
	this.mutex.Unlock()
	return pageContent
}

func (this *EmailMiner)DetectEmail(webUrl string)[]string {
	this.emailFilterMap = make(map[string]struct{})
	eUrl,err := url.Parse(webUrl)
	if err != nil{
		return nil
	}
	eDomain,err := publicsuffix.Domain(eUrl.Host)
	if err != nil{
		return nil
	}
	if _,bExists := Const.BlackDomainMap[eDomain];bExists==true{
		log.Println("网址域名过滤:",webUrl)
		return nil
	}
	_,bLocation := Const.LocationDomainMap[eDomain]
	pageContent := this.exploreUrl(webUrl)
	if len(this.emailList) > 0{
		return this.emailList
	}
	if bLocation == false{
		var contactDetector ContactDetector.ContactDetector
		pathList := contactDetector.ExtractContactUrl(webUrl,pageContent)
		for _,ePathUrl := range pathList {
			this.exploreUrl(ePathUrl)
		}
	}
	return this.emailList
}