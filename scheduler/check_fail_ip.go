package scheduler

import (
	"proxy-collect/component"
	"proxy-collect/component/logger"
	"proxy-collect/model"
	"proxy-collect/service"
)

type CheckFailIp struct {
}

func (s CheckFailIp) Run() {
	logger.Success("check fail ip start run")
	var proxies []model.Proxy
	model.DB.Where("status<>?", 1).Where("check_count>0").Find(&proxies)
	logger.Info("count:%d, cap: %d\n", len(proxies), cap(proxies))
	pool := component.NewTaskPool(40)
	defer pool.Close()
	for _, proxy := range proxies {
		var proxyTmp model.Proxy = proxy
		pool.RunTask(func() { service.ProxyService.CheckProxyAndSave(proxyTmp.Host, proxyTmp.Port, "") })
	}
}
