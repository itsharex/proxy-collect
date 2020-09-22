package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"proxy-collect/component"
	"proxy-collect/component/logger"
	"proxy-collect/model"
	"regexp"
	"sync"
	"time"
)

func NewProxyService() *proxyService {
	return &proxyService{}
}

type proxyService struct {
}

func (s *proxyService) CheckIpStatusActive(host, port string) bool {
	request_url := "https://www.baidu.com"
	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	proxyServer := fmt.Sprintf("http://%s:%s", host, port)
	proxyUrl, _ := url.Parse(proxyServer)
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http get error", err)
		return false
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("http read error", err)
		return false
	}
	return true
}

func (s *proxyService) CheckIpStatus(host, port string) bool {
	request_url := "https://www.c5game.com/api/product/sale.json?id=2705689&page=1&sort=1&key=1523539522"
	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	proxyServer := fmt.Sprintf("http://%s:%s", host, port)
	proxyUrl, _ := url.Parse(proxyServer)
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("http get error", err)
		return false
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return false
	}
	//fmt.Println(string(body))
	return true
}

func (s *proxyService) CheckProxyAndSave(host string, port string, db *gorm.DB) {
	result := s.CheckIpStatus(host, port)
	if result {
		logger.Success(result, host, port)
	} else {
		logger.Warning(result, host, port)
	}
	var status int8 = 1
	if !result {
		status = 0
		return
	}
	var proxyModel model.Proxy
	err := db.Where("host = ? AND port = ?", host, port).First(&proxyModel).Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		proxyModel = model.Proxy{
			Host:       host,
			Port:       port,
			Status:     status,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		db.Create(&proxyModel)
		return
	}
	proxyModel.Status = status
	proxyModel.UpdateTime = time.Now().Unix()
	db.Save(&proxyModel)
	return
}

func (s *proxyService) DoGetProxy(getProxyService GetProxyInterface, pool *component.Pool, db *gorm.DB) {
	for _, requestUrl := range getProxyService.GetUrlList() {
		contentBody := getProxyService.GetContentHtml(requestUrl)
		if contentBody == "" {
			time.Sleep(time.Second * 5)
			continue
		}
		proxy_list := getProxyService.ParseHtml(contentBody)
		logger.Info("获取到ip:", proxy_list)
		var wg sync.WaitGroup = sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			logger.Info("wait 20s ...")
			time.Sleep(time.Second * 20)
		}(&wg)
		for _, proxyArr := range proxy_list {
			ip, port := proxyArr[0], proxyArr[1]
			pool.RunTask(func() { s.CheckProxyAndSave(ip, port, db) })
		}

		wg.Wait()
	}
}

func (s *proxyService) CheckProxyFormat(host string, port string) bool {
	ok, _ := regexp.Match(`^[\d\.]+$`, []byte(host))
	if !ok {
		return false
	}
	ok, _ = regexp.Match(`^\d+$`, []byte(port))
	if !ok {
		return false
	}
	return true
}
