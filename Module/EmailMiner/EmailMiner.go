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
	//遇到这些网站就过滤
	blackDomainMap = map[string]struct{}{
		"twitter.com":{},
		"huawei.com":{},
		"youtube.com":{},
		"snapchat.com":{},
		"instagram.com":{},
		"business.site":{},
	}
	//遇到这些网站不用猜其它路径
	locationDomainMap = map[string]struct{}{
		"dhl.com":{},
		"facebook.com":{},
	}
	pathList = []string{
		"contactar",
		"contactez",
		"contacter",
		"aloqa",
		"contatti",
		"contact",
		"contacte",
		"contacto",
		"contato",
		"contatto",
		"cysylltwch",
		"fifandraisana",
		"hubungan",
		"hubungi",
		"ikopanya",
		"kapcsolatba",
		"Kontak",
		"kontakt",
		"kontakta",
		"kontaktas",
		"kontakts",
		"kontaktua",
		"kukhudzana",
		"kuntatt",
		"lamba",
		"makipag-ugnay",
		"oxhumana",
		"stik",
		"temas",
		"wasiliana",
		"contact-us",
		"xiriir",
		"iletisim",
		"bizeulasin",
		"adborth",
		"aiseolas",
		"atsauksmes",
		"Atsiliepimas",
		"balik",
		"bildirim",
		"comentarios",
		"comentaris",
		"feedback",
		"feed-back",
		"fidbak",
		"informacije",
		"informacje",
		"information",
		"iritzia",
		"maklum",
		"maoni",
		"ndemanga",
		"nzaghachi",
		"palaute",
		"parere",
		"reagim",
		"rispons",
		"risposta",
		"tagasiside",
		"Tanggepan",
		"terugkoppeling",
		"terugvoer",
		"tilbakemeldinger",
		"tlhahiso",
		"urupare",
		"afdruk",
		"aftryk",
		"akara",
		"alama",
		"anprent",
		"avtryck",
		"avtrykk",
		"aztarna",
		"bosma",
		"damga",
		"empremta",
		"Impressum",
		"impresszum",
		"imprima",
		"imprimer",
		"imprimir",
		"imprint",
		"impronta",
		"Isamisi",
		"istampar",
		"izdruka",
		"jejak",
		"kaluaran",
		"marika",
		"mongolo",
		"odcisk",
		"odtis",
		"otisak",
		"otisk",
		"printiad",
		"riix",
		"shafi",
		"spaudas",
		"tohu",
		"zolemba",
		"hakkimizda",
	}
	suffixList = []string{".webp",".gif",".htm",".html",".jpeg",".jpg",".php",".png", ".ace",".ani",".arc",".arj",".avi",".bmp",".cab",".class",".css",".exe",".ico",".jar",".mid",".mov",".mp2",".mp3",".mpeg",".mpg",".pdf",".wix"}
	regex_Email = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}`)
)

type EmailMiner struct {
	wg sync.WaitGroup
	mutex sync.Mutex
	emailList []string
	emailFilterMap map[string]struct{}
	concurrencyCh chan struct{}
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

func (this *EmailMiner)exploreUrl(webUrl string)  {
	defer this.wg.Done()
	this.concurrencyCh <- struct{}{}
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
	<- this.concurrencyCh
}


func (this *EmailMiner)DetectEmail(webUrl string)[]string {
	this.emailFilterMap = make(map[string]struct{})
	this.concurrencyCh = make(chan struct{},3)
	eUrl,err := url.Parse(webUrl)
	if err != nil{
		return nil
	}
	eDomain,err := publicsuffix.Domain(eUrl.Host)
	if err != nil{
		return nil
	}
	if _,bExists := blackDomainMap[eDomain];bExists==true{
		log.Println("网址域名过滤:",webUrl)
		return nil
	}
	_,bLocation := locationDomainMap[eDomain]
	this.wg.Add(1)
	go this.exploreUrl(webUrl)
	if bLocation == false{
		for _,ePath := range pathList {
			expUrl := fmt.Sprintf("%s://%s/%s",eUrl.Scheme,eUrl.Host,ePath)
			this.wg.Add(1)
			go this.exploreUrl(expUrl)
		}
	}
	this.wg.Wait()
	return this.emailList
}