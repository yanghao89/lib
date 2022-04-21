package nacos

import (
	"bytes"

	"lib/config"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	configClient "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	configRwMutex sync.RWMutex
	nacosConfig   map[string]*viper.Viper
)

func newClient(nc config.NacosConfig) (configClient.IConfigClient, error) {
	if nacConfig == nil {
		nacConfigMutex.Lock()
		defer nacConfigMutex.Unlock()
		if nacConfig == nil {
			cc, sc := getConfig(nc)
			//服务注册, 服务连接
			cli, err := clients.NewConfigClient(vo.NacosClientParam{
				ClientConfig:  cc,
				ServerConfigs: sc,
			})
			if err != nil {
				return nil, err
			}
			nacConfig = cli
		}
	}
	return nacConfig, nil
}

func NewConfigClient(nc config.NacosConfig, cp []ConfigParam) error {
	ncl, err := newClient(nc)
	if err != nil {
		return err
	}
	for k := range cp {
		defaultConfig := viper.New()
		defaultConfig.SetConfigType("json")
		param := cp[k]
		content, err := ncl.GetConfig(vo.ConfigParam{
			DataId: param.DataId,
			Group:  param.Group,
		})
		if err != nil {
			return err
		}
		if err := defaultConfig.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
			return err
		}
		nacosConfig[param.Name] = defaultConfig
	}
	go func() {
		for k := range cp {
			defaultConfig := viper.New()
			param := cp[k]
			defaultConfig.SetConfigType("json")
			_ = ncl.ListenConfig(vo.ConfigParam{
				DataId: param.DataId,
				Group:  param.Group,
				OnChange: func(namespace, group, dataId, data string) {
					if err := defaultConfig.ReadConfig(bytes.NewReader([]byte(data))); err == nil {
						configRwMutex.Lock()
						nacosConfig[param.Name] = defaultConfig
						configRwMutex.Unlock()
					}
				},
			})
		}
	}()
	return nil
}
