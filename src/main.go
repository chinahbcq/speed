/**
 * @file src/main.go
 * @author chinahcbq (chinahbcq@qq.com)
 * @date 2016-04-29 15:27:13
 * @brief main.go相关操作
 *
 **/

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)
import (
	"controller"
	"rpc_framework"
	"ublog"
	"utils"
)

func loadSysConfigs(file string) error {
	configStr, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configStr, &utils.Global.Config)
	if err != nil {
		panic(err)
	}
	return nil
}
func loadErrorConfigs(file string) error {
	configStr, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configStr, &utils.ErrorMap)
	if err != nil {
		return err
	}
	return nil
}

func initLogger() {
	config := utils.Global.Config
	ubLogInfo := &ublog.UbLogInfo{
		ChannelLen:     config.LogChannelLen,
		FlushThreadNum: config.LogFlushThreadNum,
		FlushInterval:  config.LogFlushInterval,
	}
	if "" != config.LogFile {
		ublog.UbLogFdPool[3] = &ublog.UbLog{}
		ublog.UbLogFdPool[3].Init(config.LogDir,
			config.LogFile,
			config.LogLevel,
			ubLogInfo,
			false)
		utils.Global.Logger = ublog.UbLogFdPool[3]
	} else {
		utils.Global.Logger = ublog.UbLogFdPool[2]
	}
	go func() {
		for {
			utils.Global.Logger.CheckLogFile()
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func main() {
	//1 读取系统配置
	loadSysConfigs("conf/speed.conf")

	//2 初始化log
	initLogger()

	configs := utils.Global.Config
	logger := utils.Global.Logger
	logger.UbLogNotice("init ublog OK")

	//3 读取错误码配置
	err := loadErrorConfigs("conf/error_code.conf")
	if err != nil {
		panic(err)
	}
	logger.UbLogNotice("load error_code.conf OK")

	runtime.GOMAXPROCS(configs.ProcessNum)

	//4 开启TCP服务
	router := &controller.Router{}
	tcpServer := rpc_framework.TcpServer{
		ListenAddress:   configs.ListenAddressTcp,
		IdleTimeout:     configs.IdleTimeout,
		Logger:          logger,
		TerminateSignal: false,
		Router:          router,
	}
	go func() {
		tcpServer.Start()
	}()
	logger.UbLogNotice("start tcp service at %s OK", configs.ListenAddressTcp)

	//5 响应退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	tcpServer.TerminateSignal = true
	logger.UbLogNotice("service exist!")
}
