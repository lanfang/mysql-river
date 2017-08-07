package sync

import (
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/go-mysql/canal"
	"github.com/lanfang/go-mysql/mysql"
	"github.com/lanfang/go-mysql/replication"
	"github.com/lanfang/go-mysql/schema"
	"github.com/lanfang/mysql-river/protocol"
	"github.com/siddontang/go/hack"
)

type posSaver struct {
	pos   mysql.Position
	force bool
}

type eventHandler struct {
	c *SyncClient
}

func (h *eventHandler) OnRotate(e *replication.RotateEvent) error {
	pos := mysql.Position{
		string(e.NextLogName),
		uint32(e.Position),
	}

	h.c.syncCh <- posSaver{pos, true}

	return h.c.ctx.Err()
}

func (h *eventHandler) OnDDL(nextPos mysql.Position, _ *replication.QueryEvent) error {
	h.c.syncCh <- posSaver{nextPos, true}
	return h.c.ctx.Err()
}

func (h *eventHandler) OnXID(nextPos mysql.Position) error {
	h.c.syncCh <- posSaver{nextPos, false}
	return h.c.ctx.Err()
}
func (h *eventHandler) OnGTID(mysql.GTIDSet) error {
	return nil
}

func (h *eventHandler) OnRow(e *canal.RowsEvent) error {
	rule, ok := h.c.getFilterInfo(e.Table.Schema, e.Table.Name)
	if !ok {
		log.Info("%+v.%+v filtered, continue", e.Table.Schema, e.Table.Name)
		return nil
	}
	data := protocol.EventData{}
	data.Action = protocol.EventAction(e.Action)
	data.Schema = e.Table.Schema
	data.Table = e.Table.Name
	var hit_index bool
	for index, col := range e.Table.Columns {
		if !hit_index && rule.KeyIndex != defaultKeyIndex && col.Name == rule.Key {
			rule.KeyIndex = index
			hit_index = true
		}
		data.Columns = append(data.Columns, col.Name)
		if col.Type == schema.TYPE_TEXT {
			for i, row := range e.Rows {
				if index < len(row) {
					if t, ok := e.Rows[i][index].([]byte); ok {
						e.Rows[i][index] = hack.String(t)
					}
				}
			}
		}

	}
	//should not come here
	if rule.KeyIndex >= len(data.Columns) || rule.KeyIndex == defaultKeyIndex {
		rule.KeyIndex = 0
	}
	data.Owner = *rule
	data.Rows = e.Rows
	h.c.syncCh <- data
	return nil
}

func (h *eventHandler) String() string {
	return "BinlogEventHandler"
}
