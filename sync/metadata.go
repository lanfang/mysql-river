package sync

import (
	"encoding/json"
	"fmt"
	"github.com/lanfang/go-lib/ha"
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/go-mysql/mysql"
	"github.com/lanfang/mysql-river/config"
	"golang.org/x/net/context"
	"sync"
	"time"
)

func NewMetaInfo(group string) (*MetaInfo, error) {
	m := &MetaInfo{}
	key := ha.EtcdClient.PositionKey(config.SERVERNAME, group)
	resp, err := ha.EtcdClient.Api().Get(context.Background(), key, nil)
	if err != nil {
		log.Info("get meta from etcd key:%v, err:%v", key, err)
		return nil, err
	}
	if resp.Node == nil {
		return nil, fmt.Errorf("node is empty")
	}
	log.Info("get meta from etcd key:%v, value:%+v", key, *(resp.Node))
	err = json.Unmarshal([]byte(resp.Node.Value), m)
	m.Group = group
	return m, nil
}

type MetaInfo struct {
	sync.RWMutex
	//binlog filename
	Name string
	//binlog position
	Pos uint32
	//server group, check server.group == meta.group
	Group string
	//save time
	LastSaveTime time.Time
	//server self role
	MyRole ha.NodeStatus
}

func (p *MetaInfo) Save(pos mysql.Position) error {
	p.Lock()
	defer p.Unlock()

	if p.MyRole != ha.Master {
		log.Error("only master can write metadata")
		return fmt.Errorf("slave server, cannot, save meta info")
	}
	p.Name, p.Pos, p.LastSaveTime = pos.Name, pos.Pos, time.Now()
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	key := ha.EtcdClient.PositionKey(config.SERVERNAME, p.Group)
	if _, err = ha.EtcdClient.Api().Set(context.Background(), key, string(d), nil); err != nil {
		log.Info("get meta from etcd key:%v, err:%v", key, err)
		return err
	}
	log.Info("save position to etcd key:%v, pos:%+v", key, pos)
	return nil
}

func (p *MetaInfo) Position() mysql.Position {
	p.Lock()
	defer p.Unlock()
	return mysql.Position{Name: p.Name, Pos: p.Pos}
}

func (p *MetaInfo) Close() {
	pos := p.Position()
	p.Save(pos)
	p.MyRole = ha.Slave
	return
}
