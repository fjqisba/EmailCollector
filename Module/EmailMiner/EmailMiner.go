package EmailMiner

import (
	"github.com/panjf2000/ants/v2"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"net/url"
	"strings"
)

var(
	blackDomainMap = map[string]struct{}{
		"123":{},
	}
)

type EmailMiner struct {

}

func (this *EmailMiner)DetectEmail(id int,webUrl string)error {
	eUrl,err := url.Parse(webUrl)
	if err != nil {
		return err
	}
	host := eUrl.Host
	tIndex := strings.IndexByte(eUrl.Host,':')
	if tIndex != -1{
		host = eUrl.Host[0:tIndex]
	}
	strDomain,err := publicsuffix.Domain(host)
	if _,bExists := blackDomainMap[strDomain];bExists == true{
		return nil
	}
	return nil
}

func DetectEmail(hPool *ants.Pool,id int,webUrl string)[]string  {
	var miner EmailMiner
	miner.DetectEmail(id,webUrl)
	return nil
}