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

type LogLine struct {
	out 		*Klog
	encode 		bool
	ctime		time.Time
	btime		sql.NullTime
	cname   	sql.NullString
	fname   	sql.NullString
	client		sql.NullString
	userid   	sql.NullString
	feilds 		[]LogFeild
}

func (line *LogLine)Reset() {
	line.out = nil
	line.encode  = false
	line.btime.Valid = false
	line.cname.Valid = false
	line.fname.Valid = false
	line.client.Valid = false
	line.userid.Valid = false
	line.feilds = line.feilds[:0]
}

var linePool *sync.Pool
func init() {
	linePool = &sync.Pool {
		New: func() interface{} {
			return new(LogLine)
		},
	}
}

func getLineEnt() *LogLine {
	return linePool.Get().(*LogLine)
}

func putLineEnt(line *LogLine)  {
	line.Reset()
	linePool.Put(line)
}

func (line *LogLine)Output() {
	lg := line.out
	if lg == nil {
		loglog.Error("klog is nil")
		putLineEnt(line)
		return
	}

	lg.Output(line)
}

func (line *LogLine)Begin() *LogLine {
	line.btime.Time = time.Now()
	line.btime.Valid = true
	return line
}

func (line *LogLine)BeginTime(btime time.Time) *LogLine {
	line.btime.Time = btime
	line.btime.Valid = true
	return line
}

func (line *LogLine)ClassName(val string) *LogLine {
	line.cname.String = val
	line.cname.Valid = true
	return line
}

func (line *LogLine)FuncName(val string) *LogLine {
	line.fname.String = val
	line.fname.Valid = true
	return line
}

func (line *LogLine)Client(val string) *LogLine {
	line.client.String = val
	line.client.Valid = true
	return line
}

func (line *LogLine)UserId(val string) *LogLine {
	line.userid.String = val
	line.userid.Valid = true
	return line
}

func (line *LogLine)Add(fld string, val interface{}) *LogLine {
	if line.feilds == nil {
		line.feilds = make([]LogFeild, 0, 100)
	}
	log_fld := LogFeild {fld, val}
	line.feilds = append(line.feilds, log_fld)
	return line
}

func (line *LogLine)WriteEncodeString(pbuf *LogBuffer, data string)  {
	var i, beg int = 0, 0
	nn := len (data)
	for i=0; i < nn ; i++ {
		if data[i] == '`' {
			pbuf.WriteString(data[beg:i])
			pbuf.WriteByte('\\')
			pbuf.WriteByte('`')
			line.encode = true
			beg = i+1
		} else if data[i] == '\n' {
			pbuf.WriteString(data[beg:i])
			pbuf.WriteByte('\\')
			pbuf.WriteByte('\n')
			line.encode = true
			beg = i+1
		}
	}
	if beg < nn {
		pbuf.WriteString(data[beg:nn])
		beg = nn
	}
}

func (line *LogLine)WriteLogFieldSpliter(pbuf *LogBuffer) {
	pbuf.WriteString(" `")
}

func (line *LogLine)WriteLogLineSpliter(pbuf *LogBuffer) {
	pbuf.WriteString(" `")
	if line.encode {
		pbuf.WriteByte('T')
	} else {
		pbuf.WriteByte('F')
	}
	pbuf.WriteString(" \n")
}

func (line *LogLine)Encode(pbuf *LogBuffer) {
	pbuf.Reset()
	pbuf.WriteString("ctm=")
	pbuf.WriteTimeString(line.ctime)
	line.WriteLogFieldSpliter(pbuf);pbuf.WriteString("class=")
	line.WriteEncodeString(pbuf, line.cname.String)
	if line.fname.Valid {
		line.WriteLogFieldSpliter(pbuf);pbuf.WriteString("func=")
		line.WriteEncodeString(pbuf, line.fname.String)
	}
	if line.client.Valid {
		line.WriteLogFieldSpliter(pbuf);pbuf.WriteString("client=")
		line.WriteEncodeString(pbuf, line.client.String)
	}
	if line.userid.Valid {
		line.WriteLogFieldSpliter(pbuf);pbuf.WriteString("userid=")
		line.WriteEncodeString(pbuf, line.userid.String)
	}
	if line.btime.Valid == true {
		line.WriteLogFieldSpliter(pbuf);pbuf.WriteString("etm=")
		ddd := line.ctime.Sub(line.btime.Time)
		pbuf.WriteString(fmt.Sprint(ddd.Nanoseconds()/1e6))
	}

	for i:=0; i < len(line.feilds); i++ {
		line.WriteLogFieldSpliter(pbuf);
		pbuf.WriteString(line.feilds[i].Field)
		pbuf.WriteByte('=')
		if val, ok := line.feilds[i].Value.(string); ok {
			line.WriteEncodeString(pbuf, val)
		} else {
			pbuf.WriteString(fmt.Sprint(line.feilds[i].Value))
		}
	}

	line.WriteLogLineSpliter(pbuf)
}

func GetLineEncodeFlag(data []byte) bool {
	nn := len(data)
	if nn > 3 {
		if data[nn-1] == 'T' && data[nn-2] == '`' && data[nn-3] == ' ' {
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
