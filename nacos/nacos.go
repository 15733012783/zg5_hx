package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

type T struct {
	Mysql struct {
		Username  string `json:"Username"`
		Password  string `json:"Password"`
		Host      string `json:"Host"`
		Port      string `json:"Port"`
		Mysqlbase string `json:"Mysqlbase"`
	} `json:"mysql"`
	Grpc struct {
		Address string `json:"Address"`
		Host    string `json:"Host"`
		Port    int    `json:"port"`
	} `json:"grpc"`
	Consul struct {
		Name string `json:"name"`
		Host string `json:"Host"`
	} `json:"consul"`
}

var GoodsT T
var success bool

func createClientConfig() (constant.ClientConfig, []constant.ServerConfig) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         "",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "10.2.171.70",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}
	return clientConfig, serverConfigs
}

func NaCosConfig(Group, DataId string, Port int) {
	clientConfig, serverConfigs := createClientConfig()
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return
	}
	config, err3 := client.GetConfig(vo.ConfigParam{
		DataId: DataId,
		Group:  Group,
	})
	if err3 != nil {
		return
	}
	json.Unmarshal([]byte(config), &GoodsT)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	success, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "10.2.171.70",
		Port:        8881,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
}

func NaocsServiceDiscovery(Group, DataId string) {
	clientConfig, serverConfigs := createClientConfig()
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})
	if err != nil {
		log.Println(err)
		return
	}
	for i, v := range instances {
		fmt.Println(i, v, "nacos,********************************")
	}
}

func ListenConfig(client config_client.IConfigClient) {
	//Listen config change,key=dataId+group+namespaceId.
	err := client.ListenConfig(vo.ConfigParam{
		DataId: "test-data",
		Group:  "test-group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", NaCosT.Username,
			//	NaCosT.Password, NaCosT.Host, NaCosT.Port, NaCosT.Mysqlbase)
			//updateDbConnection(dsn)
		},
	})
	if err != nil {
		log.Println(err)
	}
}
