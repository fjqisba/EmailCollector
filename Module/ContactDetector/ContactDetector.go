package ContactDetector

import (
	"EmailCollector/Const"
	"github.com/PuerkitoBio/goquery"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"net/url"
	"strings"
)

type ContactDetector struct {
	//原始网址域名
	rawDomain string
}

func (this *ContactDetector)FilterUrl(urlPath string)bool {
	if strings.HasSuffix(urlPath,".png"){
		return true
	}
	return false
}

//拿到网址中所有的href地址

func (this *ContactDetector)GetAllHRefList(doc *goquery.Document,webUrl string)(retList []string){
	filterMap := make(map[string]struct{})
	eUrl,err := url.Parse(webUrl)
	if err != nil{
		return nil
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, bExists := s.Attr("href")
		if bExists == false {
			return
		}
		if href == "" || href == "/" || href == "#"{
			return
		}
		var linkAddr string
		if strings.HasPrefix(href,"http") == false{
			linkAddr = eUrl.Scheme + "://" + eUrl.Host + href
		}else{
			linkAddr = href
		}
		_,bExists = filterMap[linkAddr]
		if bExists == false{
			filterMap[linkAddr] = struct{}{}
			retList = append(retList, linkAddr)
		}
	})
	return retList
}

//联系Url链接的提取

func (this *ContactDetector)ExtractContactUrl(webUrl string,pageContent string)[]string  {
	eUrl,err := url.Parse(webUrl)
	if err != nil{
		return nil
	}
	this.rawDomain,err = publicsuffix.Domain(eUrl.Host)
	if err != nil{
		return nil
	}
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(pageContent))
	if err != nil{
		return nil
	}
	hRefList := this.GetAllHRefList(doc,webUrl)
	var retContactList []string
	for _,eHRefUrl := range hRefList{
		newUrl,err := url.Parse(eHRefUrl)
		if err != nil{
			continue
		}
		newDomain,err := publicsuffix.Domain(newUrl.Host)
		if err != nil{
			continue
		}
		//存在联系相关地址
		if _,bLocate := Const.LocationDomainMap[newDomain];bLocate == true{
			retContactList = append(retContactList, eHRefUrl)
			continue
		}
		//跨域网址忽略掉
		if newDomain != this.rawDomain{
			continue
		}
		newUrl.Path = strings.ToLower(newUrl.Path)
		if this.FilterUrl(newUrl.Path){
			continue
		}
		//寻找是否包含有联系关键词
		for _,eContactPath := range Const.ContactKeyWords{
			if strings.Contains(newUrl.Path,eContactPath) == true{
				retContactList = append(retContactList, eHRefUrl)
				break
			}
		}
	}
	return retContactList
}
