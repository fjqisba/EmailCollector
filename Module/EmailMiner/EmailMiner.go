package EmailMiner

import (
	"EmailCollector/Module/PhantomJS"
	"fmt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)


var(
	blackDomainMap = map[string]struct{}{
		"facebook.com":{},
		"twitter.com":{},
		"huawei.com":{},
	}
	pathList = []string{
		"contact","contact-us",
	}
	regex_Email = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}`)
)

type EmailMiner struct {
	wg sync.WaitGroup
	mutex sync.Mutex
	emailList []string
	emailFilterMap map[string]struct{}
}

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

func (this *EmailMiner)exploreUrl(webUrl string)  {
	log.Println("检测网址:",webUrl)
	defer this.wg.Done()
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
	if _,bExists := blackDomainMap[eDomain];bExists==true{
		return nil
	}
	this.wg.Add(1)
	go this.exploreUrl(webUrl)
	for _,ePath := range pathList {
		expUrl := fmt.Sprintf("%s://%s/%s",eUrl.Scheme,eUrl.Host,ePath)
		this.wg.Add(1)
		go this.exploreUrl(expUrl)
	}
	this.wg.Wait()
	return this.emailList
}