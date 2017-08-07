package main

import (
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/mysql-river/protocol"
	"github.com/lanfang/mysql-river/queue"
	"github.com/lanfang/mysql-river/utils"
	"time"
)

func init() {
	utils.RegisterModel(new(Goods))
}

// 可以通过添加 col struct tag来映射数据库字段名,也可以自动映射 GoodsId <-> goods_id
type Goods struct {
	GoodsId     string `col:"goods_id"`
	Enabled     int    `col:"enabled"`
	SpIdTag     string `col:"sp_id"`
	SpType      int
	SpName      string `col:"sp_name"`
	ExpiredTime time.Time
	ModifyTime  time.Time `col:"modify_time"`
}

//可以不提供词函数，会自动转换
func (this *Goods) TableName() string {
	return "goods"
}

var addr []string = []string{"localhost:9092"}
var servername string = "democonsumer"

func Handle(msg *protocol.EventData) error {
	log.Info("Handle get queue msg %+v", *msg)
	container := make([]Goods, 0)
	err := msg.ToObjArrary(&container)
	log.Info("ToObjArrary, cnt:%v, container:%+v, err:%+v", len(container), container, err)
	switch msg.Action {
	case protocol.UpdateAction:
		//如果是更新操作，binlog会体现更改前和更改后的记录
		for i := 0; i+1 < len(container); {
			log.Info("before %v ,update %+v", msg.Action, container[i])
			log.Info("after %v, update %+v", msg.Action, container[i+1])
			i = i + 2
		}
	case protocol.InsertAction, protocol.DeleteAction:
		for i := 0; i < len(container); i++ {
			log.Info("%v, %+v", msg.Action, container[i])
		}
	}
	jsonmsg, err2 := msg.ToJsonArrary()
	log.Info("ToJson, jsonmsg:%+v, err:%+v", jsonmsg, err2)

	return nil
}

func run() {
	log.Gen(servername, "./", servername+".log")
	cc := &queue.ConsumerCfg{
		TopicList: []string{"mall_db", "social_db", "order"}, GrouponId: servername, Addr: addr,
	}
	queue.SetConsumer(&queue.SimpleConsumer{
		Cfg: cc, Work: Handle,
	})
	queue.StartConsume()
}

func main() {
	run()
	wait := make(chan int, 1)
	<-wait
}
