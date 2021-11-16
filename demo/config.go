package demo

import (
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/net4go"
	"github.com/uber/jaeger-client-go/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	etcdClient *clientv3.Client
)

func init() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.77:2379"}
	//config.Endpoints = []string{"127.0.0.1:2379"}
	var err error
	etcdClient, err = clientv3.New(config)
	if err != nil {
		log4go.Println(nil, err)
		os.Exit(-1)
	}
}

func GetETCDClient() *clientv3.Client {
	return etcdClient
}

func GetIPAddress() string {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log4go.Println(nil, err)
		return ""
	}
	listener.Close()
	return listener.Addr().String()
}

func Wait() {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
MainLoop:
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break MainLoop
		}
	}
}


type Config struct {
	config.Configuration
}

func Load(path string) (cfg *Config) {
	cfgData, _ := ioutil.ReadFile(path)

	var jCfg config.Configuration
	yaml.Unmarshal(cfgData, &jCfg)

	cfg = &Config{}
	cfg.Configuration = jCfg

	return cfg
}