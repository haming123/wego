package klog

import (
	"errors"
	"sync"
	"time"
)

const tm_flush_loop = 10
type Klog struct {
	mu     		sync.Mutex
	out    		*FileWriter
	buf 		*LogBuffer
	sender 		LogSender
	close_chan	chan bool
	is_closed 	bool
}

func NewKlog(lpath string, rtype RotateType) *Klog {
	var log Klog
	log.out = NewFileWriter(get_file_path(lpath), rtype)
	log.is_closed = false
	log.close_chan = make(chan bool)
	log.buf = NewLogBuffer()
	log.sender = nil
	go log.flush_loop()

	loglog.Debug("New Klog")
	return &log
}

func (m *Klog) SetLogSender(sender LogSender) error {
	if sender == nil {
		return errors.New("sender is nil")
	}
	if m.sender != nil {
		return errors.New("klog's sender is not nil")
	}
	m.sender = sender
	m.sender.SetLogger(m)
	go m.sender.BeginSendLoop()
	return nil
}

func (m *Klog) flush_loop()  {
	for {
		select {
		case <-time.After(tm_flush_loop * time.Second):
			m.Flush()
		case <-m.close_chan:
			m.is_closed = true
			loglog.Debug("stop flush loop")
			return
		}
	}
}

func (m *Klog) Flush() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.out.Flush()
}

func (m *Klog) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.out.Close()
	if m.close_chan != nil {
		close(m.close_chan)
		m.close_chan = nil
	}

	loglog.Debug("Close Klog")
	return nil
}

func (m *Klog) NewL(class_name string) *LogLine {
	line := getLineEnt()
	line.out = m
	line.ClassName(class_name)
	return line
}

func (m *Klog)Output(line *LogLine) {
	if m.is_closed == true {
		return
	}

	m.mu.Lock()

	line.ctime = time.Now()
	line.Encode(m.buf)
	m.out.Write(line.ctime, m.buf.GetBytes())

	putLineEnt(line)
	m.mu.Unlock()
}

func (m *Klog)GetCurLogFile() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	fname := ""
	if m.out.file != nil  {
		fname = m.out.file.name
	}
	return fname
}

func (m *Klog)GenCurLogFile() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	fw := m.out
	fname := gen_file_name(time.Now(), fw.file_path, fw.rotate)
	return fname
}
