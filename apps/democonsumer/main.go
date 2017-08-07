package main

import (
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/mysql-river/protocol"
	"github.com/lanfang/mysql-river/queue"
	"github.com/lanfang/mysql-river/utils"
	"time"
)

func init() {
	utils.RegisterModel(new(Goods), new(Topic), new(OrderModel))
}

// col struct tag
type Topic struct {
	Id            int64     `col:"pk"`
	TopicId       string    `col:"topic_id"`    //主题id
	TopicType     int       `col:"topic_type"`  //主题类型: 1 晒单
	SourceType    int       `col:"source_type"` //主题类型: 1 普通用户, 2 马甲用户, 3 商家
	GoodsId       string    `col:"goods_id"`    //商品ID(晒单时,存商品ID
	FromId        string    `col:"from_id"`     //发表者ID
	FromNick      string    `col:"from_nick"`   //发表者昵称
	Text          string    `col:"text"`        //文字, 图片存在pic_info
	StarsNum      int       `col:"stars_num"`   //评分,星数
	State         int       `col:"state"`       //1 正常，2 已屏蔽，3 已置顶
	ReplyNum      int       `col:"reply_num"`   //被回复的次数
	PornState     int       `col:"porn_state"`  //色情检查状态
	CreateTimeTag time.Time `col:"create_time"`
	ModifyTime    time.Time `col:"modify_time"`
}

// col struct tag
type Goods struct {
	GoodsIdTag          string    `col:"goods_id"`
	Enabled             int       `col:"enabled"`
	SpIdTag             string    `col:"sp_id"`
	SpType              int       `col:"sp_type"` //1:自营商户 2:第三方商户
	SpName              string    `col:"sp_name"`
	Title               string    `col:"title"`         //商品短名称
	GoodsDescTestTag    string    `col:"goods_desc"`    //商品描述
	GoodsPic            string    `col:"goods_pic"`     //商品图片
	SpecSpName          string    `col:"spec_sp_name"`  //商户名(暂时没有使用
	SpecSpPhone         string    `col:"spec_sp_phone"` //商户联系方式(暂时没有用
	OutUrl              string    `col:"out_url"`
	DetailDesc          string    `col:"detail_desc"` //图文详情
	GoodsType           int       `col:"goods_type"`
	PropsTitle          string    `col:"props_title"` //商品sku描述文案
	PropList            string    `col:"prop_list"`
	BuyTimesLimit       int       `col:"buy_times_limit"`
	StockState          int       `col:"stock_state"`
	State               int       `col:"state;default(1"`
	LastState           int       `col:"last_state;default(1"`
	CreateTime          time.Time `col:"create_time"`
	CreateUserId        string    `col:"create_user_id"`
	LastUpTime          time.Time `col:"last_up_time"`
	LastDownTime        time.Time `col:"last_down_time"`
	StartSellTime       time.Time `col:"start_sell_time"`
	ExpiredTime         time.Time `col:"expired_time"`
	ModifyTime          time.Time `col:"modify_time"`
	SubmitAuditTime     time.Time `col:"submit_audit_time"`
	AuthenticGuaranteed int       `col:"authentic_guaranteed"` //1:正品；2:非正品
	MarketPrice         int64     `col:"market_price"`         //市场价
	CanUseCoupon        int64     `col:"can_use_coupon"`
	ShortAdv            string    `col:"short_adv"`            //简单广告词
	DeliveryPromise     int       `col:"delivery_promise"`     //发货承诺
	ReturnGoodsPromise  int       `col:"return_goods_promise"` //退货承诺
	UnboxState          int       `col:"unbox_state"`          //晒单状态 1  启用, 2 停用
}

//Default snake column
type OrderModel struct {
	OrderId             string `orm:"pk"`
	UserId              string
	GoodsId             string
	GoodsTitle          string
	GoodsDesc           string // 目前是商品使用的名称
	SpId                string
	SpName              string
	SpType              int //1:自营商户 2:第三方商户
	Remarks             string
	TotalFee            int64
	MailFee             int64
	UserPayAmount       int64
	RealMailFee         int64
	RecvPersonName      string
	RecvProvince        int
	RecvCity            int
	RecvCounty          int
	RecvAddress         string
	RecvPhone           string
	MailType            int ////1: 需邮寄; 2: 不邮
	MailCompany         string
	MailTicket          string
	State               int
	CancelReason        string
	PayWay              int
	PayTransId          string //收款单
	TransId             string //交易单
	PayClientNotified   int    //客户端通知支付状态: 1,支付成功;
	CreateTime          time.Time
	PayTime             time.Time
	ModifyTime          time.Time
	ExpireTime          time.Time
	Callback            string
	OuterOrderId        string
	ExternalData        string
	SettledStatus       int    //结算状态1:未结算；2：以结算'
	SettleId            string //日结算单ID
	PlatformFee         int64  //咕咚平台费用
	AmountCheckedStatus int    //对账状态1:未对账；2：已对账
	ChannelFee          int64  //通道费用
	Source              int    //标示是站内还是站外支付
	ReserveInt1         int    `orm:"column(reserve_int_1)"`
	PkgId               string `orm:"column(package_id)"`
	// 注意和订单source字段区分,这里是标识订单来源是购物车、购物袋、秒杀列表等等
	SourceType int    `orm:"column(source_type)"`
	SourceId   string `orm:"column(source_id)"`
}

func (this *OrderModel) TableName() string {
	return "order"
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
