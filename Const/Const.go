package Const

var(
	//遇到这些网站不用猜其它路径
	LocationDomainMap = map[string]struct{}{
		"dhl.com":{},
	}
	//遇到这些网站就过滤
	BlackDomainMap = map[string]struct{}{
		"twitter.com":{},
		"huawei.com":{},
		"youtube.com":{},
		"snapchat.com":{},
		"instagram.com":{},
		"business.site":{},
		"tiktok.com":{},
		"facebook.com":{},
		"whatsapp.com":{},
	}
	//联系关键词
	ContactKeyWords = []string{
		"aloqa",
		"contat",
		"contact",
		"cysylltwch",
		"feedback",
		"feed-back",
		"fidbak",
		"fifandraisana",
		"hakkimizda",
		"hubungan",
		"hubungi",
		"ikopanya",
		"impressum",
		"kapcsolatba",
		"kukhudzana",
		"kontak",
		"kuntatt",
		"oxhumana",
	}
)