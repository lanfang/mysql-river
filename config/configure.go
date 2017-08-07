package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/lanfang/go-lib/ha"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var g_conf_file, g_conf_etcd string
var G_Config ServerConfig = ServerConfig{}
var SERVERNAME string //= "binlogbroker"
func init() {
	if i := strings.LastIndex(os.Args[0], "/"); i >= 0 {
		i++
		SERVERNAME = os.Args[0][i:]
	}
}
func Init() {
	flag.StringVar(&G_Config.BrokerConfig.Group, "g", "", "this parameter is must, group name like mysql instance, ex: db_mall")
	flag.StringVar(&g_conf_etcd, "etcd", "", "etcd addr")
	flag.StringVar(&g_conf_file, "c", "", "json file config")
	flag.Parse()
	if G_Config.BrokerConfig.Group == "" {
		log.Fatal(" must provide -g parameter, ex:mysql instance name")
		os.Exit(1)
	}
	parseConfig()
	//用自己配置etcd地址
	if len(G_Config.BrokerConfig.EtcdAddr) != 0 {
		if err := ha.NewEtcdClient(G_Config.BrokerConfig.EtcdAddr...); err != nil {
			log.Fatal(" NewEtcdClient with addr:%+v, err:%+v", G_Config.BrokerConfig.EtcdAddr, err)
			os.Exit(1)
		}
	}
}

type TopicInfo struct {
	//table name
	Table string
	//to kafka topic name, default name is Schema name
	Topic string
	//kafka partion key name, default table name
	Key string
}

type Source struct {
	Schema string
	Tables []TopicInfo
}
type AlertConfig struct {
	Host         string
	Token        string
	DepartmentId []string
}

type DBConfig struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DSN    string
}

type SourceConfig struct {
	//db conn: root:password@tcp(localhost:3306)
	MysqlConn string
	//parse from MysqlConn
	DBConfig struct {
		User   string
		Passwd string
		Net    string
		Addr   string
		DSN    string
	} `json:"-"`
	//db and table list
	Sources []Source
	//if sources config is empty and SyncAll is true,broker sync all, else do nothing
	SyncAll bool
	BlackList []string //eg:schema.table
}

type BrokerConfig struct {
	//group name from cmd line
	Group   string `json:"-"`
	LogDir  string
	LogFile string
	//etcd url, from cmd line
	EtcdAddr []string
	//binlog push to kafka
	KafkaAddr []string
	//admin port
	RPCListen string
	// on master/slave switched, sending message enterprise WeChat
	//Alert AlertConfig
	Alert struct {
		Host         string
		Token        string
		DepartmentId []string
	}
}

type ServerConfig struct {
	//etcd: /v2/keys/config/binlogbroker
	BrokerConfig BrokerConfig
	//etcd: /v2/keys/config/binlogbroker:${group}, group name from the cmd line
	SourceConfig SourceConfig
}

func GetConfig(config_file string, config *ServerConfig) error {
	fmt.Println("config file:" + config_file)
	file, err := os.Open(config_file)
	if err != nil {
		return err
	}
	defer file.Close()

	config_str, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(config_str, config)
	if err != nil {
		fmt.Println("parase config error for:" + err.Error())
	}
	return err
}

func GetConfigFromEtcd(keyName string, cfg interface{}) error {
	fmt.Println("load %v config from etcd %v ", keyName, g_conf_etcd)
	return ha.EtcdClient.GetConfig(keyName, cfg)
}

func parseConfig() {
	var err error
	if g_conf_etcd != "" {
		if err := ha.NewEtcdClient(g_conf_etcd); err != nil {
			fmt.Fprintf(os.Stderr, "NewEtcdClient source  addr:%v\n", G_Config.BrokerConfig.EtcdAddr)
			os.Exit(1)
		}
		goto LoadEtcdConfig
	} else if g_conf_file != "" {
		goto LoadConfig
	} else {
		fmt.Fprintf(os.Stderr, "No configuration source\n")
		os.Exit(1)
	}

LoadConfig:
	//init config
	err = GetConfig(g_conf_file, &G_Config)
	if err != nil {
		log.Fatal("parse config file error: %s", err.Error())
		return
	}
	err = parseDSN(G_Config.SourceConfig.MysqlConn, &G_Config.SourceConfig)
	fmt.Printf("Config:%+v, err:%+v\n", G_Config, err)
	return
LoadEtcdConfig:
	key := SERVERNAME
	for loop := true; loop; loop = false {
		if err = GetConfigFromEtcd(SERVERNAME, &G_Config.BrokerConfig); err != nil {
			break
		}
		key = fmt.Sprintf("%v:%v", SERVERNAME, G_Config.BrokerConfig.Group)
		if err = GetConfigFromEtcd(key, &G_Config.SourceConfig); err != nil {
			break
		}

		if G_Config.SourceConfig.MysqlConn == "" {
			err = fmt.Errorf("SourceConfig.MysqlConn is empty")
			break
		}
		err = parseDSN(G_Config.SourceConfig.MysqlConn, &G_Config.SourceConfig)
	}
	if err != nil {
		log.Fatal("parse config file error: %s", err.Error())
	}
	fmt.Printf("Config:%+v\n", G_Config)
	return
}

//admin_root:bjfmg1nSsynKggb@tcp(timeline:3306) -> DSN
func parseDSN(dsn string, cfg *SourceConfig) error {
	var err error
	var left string = dsn
	for loop := true; loop; loop = false {
		var i int
		//user
		if i = strings.Index(left, ":"); i < 0 {
			err = errors.New("Invalid DSN: can not find user")
			break
		}
		cfg.DBConfig.User = left[0:i]
		i++
		left = left[i:]

		//password
		if i = strings.Index(left, "@"); i < 0 {
			err = errors.New("Invalid DSN: can not find passord")
			break
		}
		cfg.DBConfig.Passwd = left[0:i]
		i++
		left = left[i:]

		//addr
		if i = strings.Index(left, "("); i < 0 {
			err = errors.New("Invalid DSN: can not find addr")
			break
		}
		i++
		left = left[i:]

		if i = strings.Index(left, ")"); i < 0 {
			err = errors.New("Invalid DSN: can not find addr")
			break
		}
		cfg.DBConfig.Addr = left[0:i]
	}
	fmt.Printf("parseDSN %+v, %+v, err:%+v", dsn, *cfg, err)
	return err
}
