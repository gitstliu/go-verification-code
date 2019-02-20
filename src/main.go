package main

import (
	"facade/handler"

	"codegenerator"
	"conf"
	"fmt"
	"imagegenerator"
	"redisclusteradapter"
	"web/restadapter"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gitstliu/log4go"
)

func main() {

	defer panicHandler()

	log4go.LoadConfiguration("conf/log.xml")
	defer log4go.Close()

	conf.LoadConfigure("conf/config.toml")

	log4go.Info("Config Load Success!!!")

	restAdapter := restadapter.RestAdapter{initUrl(), conf.GetConfigure().ServicePort}

	log4go.Info("Config RedisCluster!!!")
	redisclusteradapter.CreadRedisCluster(conf.GetConfigure().RedisClusterAddresses, conf.GetConfigure().RedisClusterConnTimeout, conf.GetConfigure().RedisClusterReadTimeout, conf.GetConfigure().RedisClusterWriteTimeout, conf.GetConfigure().RedisClusterKeepAlive, conf.GetConfigure().RedisClusterAliveTime)
	log4go.Info("Config RedisCluster Finished!!!")
	imagegenerator.LoadFont("ttf/code.ttf")
	codegenerator.InitWords("words/words.txt")
	log4go.Info("Starting RestAdapter!!!")
	restAdapter.Start()

}

func initUrl() []*restadapter.UrlMap {
	vcodefacade := handler.VCodeFacade{}

	urls := make([]*restadapter.UrlMap, 0, 100)
	urls = append(urls, &restadapter.UrlMap{
		Url: "/vcode",
		MethodMap: map[string]rest.HandlerFunc{
			"GET":  vcodefacade.GetCode,
			"POST": vcodefacade.CheckCode}})

	return urls
}

func panicHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
		fmt.Printf("%T", r)
		panic(r)
	}
}
