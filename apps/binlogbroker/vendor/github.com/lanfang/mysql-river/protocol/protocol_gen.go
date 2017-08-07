package protocol

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *EventAction) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zxvk string
		zxvk, err = dc.ReadString()
		(*z) = EventAction(zxvk)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z EventAction) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z EventAction) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EventAction) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zbzg string
		zbzg, bts, err = msgp.ReadStringBytes(bts)
		(*z) = EventAction(zbzg)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z EventAction) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EventData) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zwht uint32
	zwht, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zwht > 0 {
		zwht--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Action":
			{
				var zhct string
				zhct, err = dc.ReadString()
				z.Action = EventAction(zhct)
			}
			if err != nil {
				return
			}
		case "Schema":
			z.Schema, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Table":
			z.Table, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Columns":
			var zcua uint32
			zcua, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Columns) >= int(zcua) {
				z.Columns = (z.Columns)[:zcua]
			} else {
				z.Columns = make([]string, zcua)
			}
			for zbai := range z.Columns {
				z.Columns[zbai], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "Rows":
			var zxhx uint32
			zxhx, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Rows) >= int(zxhx) {
				z.Rows = (z.Rows)[:zxhx]
			} else {
				z.Rows = make([][]interface{}, zxhx)
			}
			for zcmr := range z.Rows {
				var zlqf uint32
				zlqf, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Rows[zcmr]) >= int(zlqf) {
					z.Rows[zcmr] = (z.Rows[zcmr])[:zlqf]
				} else {
					z.Rows[zcmr] = make([]interface{}, zlqf)
				}
				for zajw := range z.Rows[zcmr] {
					z.Rows[zcmr][zajw], err = dc.ReadIntf()
					if err != nil {
						return
					}
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *EventData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Action"
	err = en.Append(0x85, 0xa6, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(string(z.Action))
	if err != nil {
		return
	}
	// write "Schema"
	err = en.Append(0xa6, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Schema)
	if err != nil {
		return
	}
	// write "Table"
	err = en.Append(0xa5, 0x54, 0x61, 0x62, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Table)
	if err != nil {
		return
	}
	// write "Columns"
	err = en.Append(0xa7, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Columns)))
	if err != nil {
		return
	}
	for zbai := range z.Columns {
		err = en.WriteString(z.Columns[zbai])
		if err != nil {
			return
		}
	}
	// write "Rows"
	err = en.Append(0xa4, 0x52, 0x6f, 0x77, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Rows)))
	if err != nil {
		return
	}
	for zcmr := range z.Rows {
		err = en.WriteArrayHeader(uint32(len(z.Rows[zcmr])))
		if err != nil {
			return
		}
		for zajw := range z.Rows[zcmr] {
			err = en.WriteIntf(z.Rows[zcmr][zajw])
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EventData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Action"
	o = append(o, 0x85, 0xa6, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, string(z.Action))
	// string "Schema"
	o = append(o, 0xa6, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61)
	o = msgp.AppendString(o, z.Schema)
	// string "Table"
	o = append(o, 0xa5, 0x54, 0x61, 0x62, 0x6c, 0x65)
	o = msgp.AppendString(o, z.Table)
	// string "Columns"
	o = append(o, 0xa7, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Columns)))
	for zbai := range z.Columns {
		o = msgp.AppendString(o, z.Columns[zbai])
	}
	// string "Rows"
	o = append(o, 0xa4, 0x52, 0x6f, 0x77, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Rows)))
	for zcmr := range z.Rows {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Rows[zcmr])))
		for zajw := range z.Rows[zcmr] {
			o, err = msgp.AppendIntf(o, z.Rows[zcmr][zajw])
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EventData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zdaf uint32
	zdaf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zdaf > 0 {
		zdaf--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Action":
			{
				var zpks string
				zpks, bts, err = msgp.ReadStringBytes(bts)
				z.Action = EventAction(zpks)
			}
			if err != nil {
				return
			}
		case "Schema":
			z.Schema, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Table":
			z.Table, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Columns":
			var zjfb uint32
			zjfb, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Columns) >= int(zjfb) {
				z.Columns = (z.Columns)[:zjfb]
			} else {
				z.Columns = make([]string, zjfb)
			}
			for zbai := range z.Columns {
				z.Columns[zbai], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "Rows":
			var zcxo uint32
			zcxo, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Rows) >= int(zcxo) {
				z.Rows = (z.Rows)[:zcxo]
			} else {
				z.Rows = make([][]interface{}, zcxo)
			}
			for zcmr := range z.Rows {
				var zeff uint32
				zeff, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Rows[zcmr]) >= int(zeff) {
					z.Rows[zcmr] = (z.Rows[zcmr])[:zeff]
				} else {
					z.Rows[zcmr] = make([]interface{}, zeff)
				}
				for zajw := range z.Rows[zcmr] {
					z.Rows[zcmr][zajw], bts, err = msgp.ReadIntfBytes(bts)
					if err != nil {
						return
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *EventData) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(string(z.Action)) + 7 + msgp.StringPrefixSize + len(z.Schema) + 6 + msgp.StringPrefixSize + len(z.Table) + 8 + msgp.ArrayHeaderSize
	for zbai := range z.Columns {
		s += msgp.StringPrefixSize + len(z.Columns[zbai])
	}
	s += 5 + msgp.ArrayHeaderSize
	for zcmr := range z.Rows {
		s += msgp.ArrayHeaderSize
		for zajw := range z.Rows[zcmr] {
			s += msgp.GuessSize(z.Rows[zcmr][zajw])
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TopicInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zrsw uint32
	zrsw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zrsw > 0 {
		zrsw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Topic":
			z.Topic, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Key":
			z.Key, err = dc.ReadString()
			if err != nil {
				return
			}
		case "KeyIndex":
			z.KeyIndex, err = dc.ReadInt()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z TopicInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Topic"
	err = en.Append(0x83, 0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Topic)
	if err != nil {
		return
	}
	// write "Key"
	err = en.Append(0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Key)
	if err != nil {
		return
	}
	// write "KeyIndex"
	err = en.Append(0xa8, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.KeyIndex)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TopicInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Topic"
	o = append(o, 0x83, 0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	o = msgp.AppendString(o, z.Topic)
	// string "Key"
	o = append(o, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendString(o, z.Key)
	// string "KeyIndex"
	o = append(o, 0xa8, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78)
	o = msgp.AppendInt(o, z.KeyIndex)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TopicInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zxpk uint32
	zxpk, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zxpk > 0 {
		zxpk--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Topic":
			z.Topic, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Key":
			z.Key, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "KeyIndex":
			z.KeyIndex, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z TopicInfo) Msgsize() (s int) {
	s = 1 + 6 + msgp.StringPrefixSize + len(z.Topic) + 4 + msgp.StringPrefixSize + len(z.Key) + 9 + msgp.IntSize
	return
}
