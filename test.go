package main

import (
	"proxy-collect/component"
	"proxy-collect/model"
	"proxy-collect/service"
	"time"
)

func main() {
	pool := component.NewTaskPool(20)
	service.ProxyService.DoGetProxy(service.GetProxyGuoBanjia, pool, model.DB)
	time.Sleep(50 * time.Second)
}
