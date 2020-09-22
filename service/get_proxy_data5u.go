package service

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"proxy-collect/component"
	"proxy-collect/component/logger"
	"proxy-collect/config"
	"strings"
)

type GetProxyData5u struct {
}

func (s *GetProxyData5u) GetUrlList() []string {
	list := []string{
		"http://www.data5u.com/",
		"http://www.data5u.com/free/gngn/index.shtml",
		"http://www.data5u.com/free/gnpt/index.shtml",
	}
	return list
}

func (s *GetProxyData5u) GetContentHtml(requestUrl string) string {

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", config.USER_AGENT)
	req.Header.Set("Host", "www.data5u.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	logger.Info("get proxy from data5u", requestUrl)
	return component.WebRequest(req)
}

func (s *GetProxyData5u) ParseHtml(body string) [][]string {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		logger.Error(err)
		return nil
	}

	var proxyList [][]string
	doc.Find("ul.l2").Each(func(i int, selection *goquery.Selection) {
		td := selection.Find("span>li").First()
		proxyHost := td.Text()
		td2 := selection.Find("span>li").Eq(1)
		proxyPort := td2.Text()
		if proxyHost == "" || proxyPort == "" {
			logger.Error("解析html node 失败")
		}
		proxyArr := []string{strings.Trim(proxyHost, " "), strings.Trim(proxyPort, " ")}
		proxyList = append(proxyList, proxyArr)
	})

	return proxyList
}
