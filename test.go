package main

import (
	"proxy-collect/component"
	"proxy-collect/component/logger"
	"proxy-collect/service"
)

var num int64 = 0

func add() {
	num = num + 1
}

func main() {

	//time.Sleep(50 * time.Second)
	//scheduler.UpdateIpInfo{}.Run()
	//res := service.ProxyService.CheckIpStatus("129.28.173.182", "8388")
	//logger.Info(res)
	//go func() {
	//	fmt.Println("pprof start...")
	//	fmt.Println(http.ListenAndServe(":9876", nil))
	//}()
	//go test()
	//for true {
	//	logger.Info(runtime.NumGoroutine())
	//	time.Sleep(time.Second)
	//}
	//scheduler.UpdateIpInfo{}.Run()
	ipInfoDto := service.ProxyService.GetIpInfo("35.196.118.22", "80")
	logger.Info("get ip info:", ipInfoDto)
}

func test() {
	pool := component.NewTaskPool(10000)
	defer pool.Close()
	for i := 0; i < 100000; i++ {
		pool.RunTask(func() {
			logger.Info("hello")
		})
	}
}
