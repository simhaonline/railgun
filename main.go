package main

import (
	"fmt"
	"github.com/MisakaSystem/LastOrder/common"
	"github.com/MisakaSystem/LastOrder/discovery"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//初始化配置
	cli, _ := discovery.NewClientDis([]string{"localhost:2379"})
	ctx, _ := cli.InitServices("/service")

	//创建反向代理
	proxy := common.NewReverseProxy(ctx)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		var http01 = http.NewServeMux()
		http01.Handle("/", proxy)
		//http01.HandleFunc("/", )
		errc <- http.ListenAndServe(":9090", http01)

	}()

	<-errc
}