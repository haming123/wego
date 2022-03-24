package klog

import (
	"bytes"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type LogFeild struct {
	Field   	string
	Value		interface{}
}

type LogRow struct {
	out 		*LogEngine
	encode 		bool
	ctime		time.Time
	btime		sql.NullTime
	cname   	sql.NullString
	fname   	sql.NullString
	client		sql.NullString
	userid   	sql.NullString
	feilds 		[]LogFeild
}

func (row *LogRow)Reset() {
	row.out = nil
	row.encode  = false
	row.btime.Valid = false
	row.cname.Valid = false
	row.fname.Valid = false
	row.client.Valid = false
	row.userid.Valid = false
	row.feilds = row.feilds[:0]
}

var linePool *sync.Pool
func init() {
	linePool = &sync.Pool {
		New: func() interface{} {
			return new(LogRow)
		},
	}
}

func getLineEnt() *LogRow {
	return linePool.Get().(*LogRow)
}

func putLineEnt(row *LogRow)  {
	row.Reset()
	linePool.Put(row)
}

func (row *LogRow)Output() {
	eng := row.out
	if eng == nil {
		putLineEnt(row)
		return
	}
	eng.Output(row)
}

func (row *LogRow)Begin() *LogRow {
	row.btime.Time = time.Now()
	row.btime.Valid = true
	return row
}

func (row *LogRow)BeginTime(btime time.Time) *LogRow {
	row.btime.Time = btime
	row.btime.Valid = true
	return row
}

func (row *LogRow)TableName(val string) *LogRow {
	row.cname.String = val
	row.cname.Valid = true
	return row
}

func (row *LogRow)FuncName(val string) *LogRow {
	row.fname.String = val
	row.fname.Valid = true
	return row
}

func (row *LogRow)ClientIP(val string) *LogRow {
	row.client.String = val
	row.client.Valid = true
	return row
}

func (row *LogRow)UserId(val string) *LogRow {
	row.userid.String = val
	row.userid.Valid = true
	return row
}

func (row *LogRow)Add(fld string, val interface{}) *LogRow {
	if row.out == nil {
		return row
	}
	if row.feilds == nil {
		row.feilds = make([]LogFeild, 0, 100)
	}
	log_fld := LogFeild {fld, val}
	row.feilds = append(row.feilds, log_fld)
	return row
}

func (row *LogRow)WriteEncodeString(pbuf *LogBuffer, data string)  {
	var i, beg int = 0, 0
	nn := len (data)
	for i=0; i < nn ; i++ {
		if data[i] == '`' {
			pbuf.WriteString(data[beg:i])
			pbuf.WriteByte('\\')
			pbuf.WriteByte('`')
			row.encode = true
			beg = i+1
		} else if data[i] == '\n' {
			pbuf.WriteString(data[beg:i])
			pbuf.WriteByte('\\')
			pbuf.WriteByte('\n')
			row.encode = true
			beg = i+1
		}
	}
	if beg < nn {
		pbuf.WriteString(data[beg:nn])
		beg = nn
	}
}

func (row *LogRow)WriteLogFieldSpliter(pbuf *LogBuffer) {
	pbuf.WriteString(" `")
}

func (row *LogRow)WriteLogLineSpliter(pbuf *LogBuffer) {
	pbuf.WriteString(" `")
	if row.encode {
		pbuf.WriteByte('.')
	} else {
		pbuf.WriteByte(' ')
	}
	pbuf.WriteString(" \n")
}

func (row *LogRow)Encode(pbuf *LogBuffer) {
	pbuf.Reset()
	pbuf.WriteString("_ctm=")
	pbuf.WriteTimeString(row.ctime)
	row.WriteLogFieldSpliter(pbuf);pbuf.WriteString("_tname=")
	row.WriteEncodeString(pbuf, row.cname.String)
	if row.fname.Valid {
		row.WriteLogFieldSpliter(pbuf);pbuf.WriteString("_fname=")
		row.WriteEncodeString(pbuf, row.fname.String)
	}
	if row.client.Valid {
		row.WriteLogFieldSpliter(pbuf);pbuf.WriteString("_cip=")
		row.WriteEncodeString(pbuf, row.client.String)
	}
	if row.userid.Valid {
		row.WriteLogFieldSpliter(pbuf);pbuf.WriteString("_uid=")
		row.WriteEncodeString(pbuf, row.userid.String)
	}
	if row.btime.Valid == true {
		row.WriteLogFieldSpliter(pbuf);pbuf.WriteString("_etm=")
		ddd := row.ctime.Sub(row.btime.Time)
		pbuf.WriteString(fmt.Sprint(ddd.Nanoseconds()/1e6))
	}

	for i:=0; i < len(row.feilds); i++ {
		row.WriteLogFieldSpliter(pbuf);
		pbuf.WriteString(row.feilds[i].Field)
		pbuf.WriteByte('=')
		if val, ok := row.feilds[i].Value.(string); ok {
			row.WriteEncodeString(pbuf, val)
		} else {
			pbuf.WriteString(fmt.Sprint(row.feilds[i].Value))
		}
	}

	row.WriteLogLineSpliter(pbuf)
}

func GetLineEncodeFlag(data []byte) bool {
	nn := len(data)
	if nn > 3 {
		if data[nn-1] == '.' && data[nn-2] == '`' && data[nn-3] == ' ' {
			return true
		}
	}
	return false
}

func LineFieldDecode(data []byte) []byte {
	nn := len(data)
	back := 0
	for i:=0; i < nn; i++ {
		if i > 0 && data[i-1] == '\\' && (data[i] == '`' || data[i] == '\n') {
			back += 1
		}
		if back > 0 {
			data[i-back] = data[i]
		}
	}
	data = data[0:nn-back]
	return data
}

func LogLineDecode(data []byte) map[string]string {
	ret := make(map[string]string)

	flag := GetLineEncodeFlag(data)
	parts := bytes.Split(data, []byte{' ', '`'})
	for _, item := range parts {
		ind := bytes.Index(item, []byte{'='})
		if ind < 1 {
			continue
		}
		value := item[ind+1:]
		if len(value) < 1 {
			continue
		}

		name := string(item[0:ind])
		if flag {
			ret[name] = string(LineFieldDecode(value))
		} else {
			ret[name] = string(value)
		}
	}

	return ret
}
