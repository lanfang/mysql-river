package ha

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
)

var EtcdClient *KeysAPI
var defaultEtcdAddr []string = []string{"http://localhost:2379"}

func init() {
	EtcdClient, _ = newDefaultClient(defaultEtcdAddr...)
}

type KeysAPI struct {
	addr []string
	api  client.KeysAPI
}

func NewEtcdClient(addrlist ...string) error {
	var err error
	EtcdClient, err = newDefaultClient(addrlist...)
	return err
}

func newDefaultClient(addrlist ...string) (*KeysAPI, error) {
	if len(addrlist) > 0 {
		defaultEtcdAddr = addrlist
	}
	cfg := client.Config{
		Endpoints:               defaultEtcdAddr,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: TTL,
	}
	c, err := client.New(cfg)
	if err != nil {
		fmt.Printf("new etcd api where addr[%+v] err:%+v", defaultEtcdAddr, err)
		return nil, err
	}

	ec := &KeysAPI{
		addr: addrlist,
		api:  client.NewKeysAPI(c),
	}
	return ec, nil
}

func (c *KeysAPI) Api() client.KeysAPI {
	return c.api
}

func (c *KeysAPI) GetConfig(servername string, cfg interface{}) error {
	rsp, err := c.api.Get(context.Background(), c.ConfigKey(servername), nil)
	if err != nil {
		fmt.Printf("read config [%s] from etcd error:%v", servername, err)
		return err
	}

	if rsp.Node == nil {
		fmt.Printf("empty etcd node")
		return fmt.Errorf("empty etcd node")
	}
	return json.Unmarshal([]byte(rsp.Node.Value), cfg)

}

func (c *KeysAPI) ConfigKey(service string) string {
	return fmt.Sprintf("/config/%s", service)
}

func (c *KeysAPI) PositionKey(service, group string) string {
	return fmt.Sprintf("/position/%s/%s", service, group)
}
