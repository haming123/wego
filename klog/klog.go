package klog

import (
	"sync"
	"time"
)

const tm_flush_loop = 10
type LogEngine struct {
	mu     		sync.Mutex
	out    		*FileWriter
	buf 		*LogBuffer
	close_chan	chan bool
	is_closed 	bool
}

func NewEngine(lpath string, rtype RotateType) *LogEngine {
	var eng LogEngine
	eng.out = NewFileWriter(get_file_path(lpath), rtype)
	eng.is_closed = false
	eng.close_chan = make(chan bool)
	eng.buf = NewLogBuffer()
	go eng.flush_loop()

	loglog.Debug("New Klog")
	return &eng
}

func (eng *LogEngine) flush_loop()  {
	for {
		select {
		case <-time.After(tm_flush_loop * time.Second):
			eng.Flush()
		case <-eng.close_chan:
			eng.is_closed = true
			loglog.Debug("stop flush loop")
			return
		}
	}
}

func (eng *LogEngine) Flush() error {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	return eng.out.Flush()
}

func (eng *LogEngine) Close() error {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	eng.out.Close()
	if eng.close_chan != nil {
		close(eng.close_chan)
		eng.close_chan = nil
	}

	loglog.Debug("Close Klog")
	return nil
}

func (eng *LogEngine) NewLog(class_name string) *LogRow {
	row := getLineEnt()
	row.out = eng
	row.TableName(class_name)
	return row
}

func (eng *LogEngine)Output(row *LogRow) {
	if eng.is_closed == true {
		return
	}

	eng.mu.Lock()

	row.ctime = time.Now()
	row.Encode(eng.buf)
	eng.out.Write(row.ctime, eng.buf.GetBytes())

	putLineEnt(row)
	eng.mu.Unlock()
}

func (eng *LogEngine)GetCurLogFile() string {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	fname := ""
	if eng.out.file != nil  {
		fname = eng.out.file.name
	}
	return fname
}

func (eng *LogEngine)GenCurLogFile() string {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	fw := eng.out
	fname := gen_file_name(time.Now(), fw.file_path, fw.rotate)
	return fname
}
