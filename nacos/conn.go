package nacos

import (
	"context"
	"fmt"
	"lib/config"
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	tGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
)

func getConfig(cc config.NacosConfig) (*constant.ClientConfig, []constant.ServerConfig) {
	return &constant.ClientConfig{
			NamespaceId:         cc.NamespaceID,
			TimeoutMs:           cc.TimeoutMs,
			NotLoadCacheAtStart: cc.NotLoadCacheAtStart,
			LogDir:              cc.LogDir,
			CacheDir:            cc.CacheDir,
			LogLevel:            cc.LogLevel,
			RotateTime:          cc.RotateTime,
			MaxAge:              cc.MaxAge,
		}, []constant.ServerConfig{
			*constant.NewServerConfig(cc.NacosServer.IP, cc.NacosServer.Port),
		}
}

var (
	nacClient      *nacos.Registry
	nacMutex       sync.Mutex
	nacConfig      config_client.IConfigClient
	nacConfigMutex sync.Mutex
)

type ConfigParam struct {
	DataId string `yaml:"dataId"` //required
	Group  string `yaml:"group"`  //required
	Name   string `yaml:"name"`
}

type NacosListenConfig struct {
	DataId   string
	Group    string
	OnChange func() func(namespace, group, dataId, data string)
}

// registry 服务注册与发现
func registry(nc config.NacosConfig) (*nacos.Registry, error) {
	if nacClient == nil {
		nacMutex.Lock()
		defer nacMutex.Unlock()
		if nacClient == nil {
			cc, sc := getConfig(nc)
			//服务注册, 服务连接
			cli, err := clients.NewNamingClient(
				vo.NacosClientParam{
					ClientConfig:  cc,
					ServerConfigs: sc,
				},
			)
			if err != nil {
				return nil, err
			}
			nacClient = nacos.New(cli)
		}
	}
	return nacClient, nil
}

type Conn interface {
	Close()
}

var (
	_         Conn = (*conn)(nil)
	clientMap      = map[string]Conn{}
	connMap        = map[string]*grpc.ClientConn{}
)

func init() {
	clientMap = make(map[string]Conn, 0)
	connMap = make(map[string]*grpc.ClientConn, 0)
}

type conn struct {
	name string
}

func newGRpcConn(name string, gc *grpc.ClientConn) (Conn, error) {
	c := new(conn)
	c.setConn(name, gc)
	return c, nil
}

//GetConn  获取链接
func GetConn(name string) *grpc.ClientConn {
	return connMap[name]
}

//设置链接
func (c *conn) setConn(name string, conn *grpc.ClientConn) {
	connMap[name] = conn
	c.name = name
}

func (c *conn) Close() {
	if c.name != "" {
		connMap[c.name].Close()
	}
}

func clientConn(ctx context.Context, endpoint string, nc config.NacosConfig) (Conn, error) {
	var (
		nr  *nacos.Registry
		err error
		cn  *grpc.ClientConn
	)
	if nr, err = registry(nc); err != nil {
		return nil, err
	}
	if cn, err = tGrpc.DialInsecure(
		ctx,
		//endpoint
		tGrpc.WithEndpoint(fmt.Sprintf("discovery:///%s.grpc", endpoint)),
		//注册服务
		tGrpc.WithDiscovery(nr),
		tGrpc.WithMiddleware(
			mmd.Client(),
			mmd.Server(),
			tracing.Client(),
			tracing.Server(),
		),
		tGrpc.WithTimeout(0),
	); err != nil {
		return nil, err
	}
	return newGRpcConn(endpoint, cn)
}

//InitClient  初始化服务
func InitClient(nc config.NacosConfig, name ...string) (func(), error) {
	ctx := context.TODO()
	for _, v := range name {
		var (
			gc  Conn
			err error
		)
		if gc, err = clientConn(ctx, v, nc); err == nil {
			clientMap[v] = gc
		} else {
			continue
		}
	}
	return func() {
		for k := range clientMap {
			clientMap[k].Close()
		}
	}, nil
}
