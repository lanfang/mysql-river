package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfang/go-lib/ha"
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/mysql-river/apps/binlogbroker/handlers"
	"github.com/lanfang/mysql-river/config"
	"github.com/lanfang/mysql-river/queue"
	"time"
)

func initSerer() error {
	config.Init()
	log.Gen(config.SERVERNAME, config.G_Config.BrokerConfig.LogDir, config.G_Config.BrokerConfig.LogFile)
	return queue.Init(config.G_Config.BrokerConfig.KafkaAddr)
}

func startServer() *ha.HAServer {
	log.Info("[%s]StartServer\n", time.Now().Format("2006-01-02 15:04:05"))
	router := gin.Default()
	router.Use(HandleTimeLoger)
	api := router.Group("/admin")
	{
		api.GET("/myrole", handlers.GetMyRole)
	}
	server, err := ha.RunAsHAServer(&handlers.BrokerServer{Name: config.SERVERNAME, Id: config.G_Config.BrokerConfig.Group})
	if err != nil {
		log.Error("RunAsHAServer Failed server:%+v, err:%+v", config.SERVERNAME, err)
	}
	go func() {
		router.Run(config.G_Config.BrokerConfig.RPCListen)
	}()
	return server
}

func HandleTimeLoger(c *gin.Context) {
	start_time := time.Now()
	defer func() {
		log.Info("Handle[%s][cost:%v]",
			c.Request.URL.Path, time.Now().Sub(start_time))
	}()
	c.Next()
}
