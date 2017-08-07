package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/lanfang/mysql-river/utils"
)

//mysql 操作类型
type EventAction string

const (
	UpdateAction EventAction = "update"
	InsertAction EventAction = "insert"
	DeleteAction EventAction = "delete"
)

type TopicInfo struct {
	Topic    string
	Key      string
	KeyIndex int
}

//go:generate msgp
type EventData struct {
	Action  EventAction
	Schema  string
	Table   string
	Columns []string
	//如果是update， Rows[i]是更新前的数据, Rows[i+1]是更新后的数据
	Rows  [][]interface{}
	Owner TopicInfo `msg:"-"`
}

func (p *EventData) GetSchemal() string {
	return p.Schema
}

func (p *EventData) GetTable() string {
	return p.Table
}

func (p *EventData) GetAction() string {
	return string(p.Action)
}

func (p *EventData) Encode() ([]byte, error) {
	return p.MarshalMsg(nil)
}

func (p *EventData) Length() int {
	return p.Msgsize()
}

func (p *EventData) ToJsonArrary() (string, error) {
	if len(p.Columns) != len(p.Rows[0]) {
		return "", fmt.Errorf("column cnt(%v) != record(%v) field cnt", len(p.Columns), len(p.Rows[0]))
	}
	result := make([]interface{}, len(p.Rows))
	for i, row := range p.Rows {
		r := make(map[string]interface{})
		for i, val := range row {
			r[p.Columns[i]] = val
		}
		result[i] = r
	}
	ret, err := json.Marshal(result)
	return string(ret), err
}

func (p *EventData) ToObjArrary(container interface{}) error {
	return utils.ToObjArrary(p.Table, p.Columns, p.Rows, container)
}
