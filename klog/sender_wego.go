package klog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

//更换数据服务器的错误时长
const tm_err_change_datanode = 60
type Config4Wego struct {
	HttpAddr	string
	AccId 		string
	AppId 		string
	HostIp 		string
	HostPort 	int
}

type Sender4Wego struct {
	logger		*Klog
	log_conf	Config4Wego
	file_path 	string
	http_url	string
	node_url	string
	file_pos	FilePos
	file 		*os.File
	buff 		*Reader
	http_client *http.Client
	tm_err 		int64
}

func NewSender4Wego(log_srv string, accid string, appid string, local_ip string, port int) *Sender4Wego {
	var conf Config4Wego
	conf.HttpAddr = log_srv
	conf.AccId = accid
	conf.AppId = appid
	conf.HostIp = local_ip
	conf.HostPort = port

	var send Sender4Wego
	send.log_conf = conf
	send.http_url = conf.HttpAddr
	send.buff = NewReader()
	send.http_client = &http.Client{Timeout: 5 * time.Second}
	return &send
}

func (sender *Sender4Wego) GetConfig() *Config4Wego {
	return &sender.log_conf
}

func (sender *Sender4Wego) SetLogger(logger *Klog)  {
	sender.logger = logger
	sender.file_path = logger.out.file_path
}

//当前发送的文件是否是在写日志文件
func (sender *Sender4Wego) IsWritingFile() bool {
	fwrite := sender.logger.GetCurLogFile()
	if fwrite == "" {
		fwrite = sender.logger.GenCurLogFile()
	}
	ret := fwrite == sender.file_pos.fname
	//loglog.Debugf("IsWriting=%v :%s\n", ret, get_short_name(fwrite))
	return ret
}

//关闭当前日志文件
//需要重置read buffer
func (sender *Sender4Wego) CloseFile()  {
	if sender.file != nil {
		sender.file.Close()
		sender.file = nil
	}
	sender.file_pos.fname = ""
	sender.file_pos.offset = 0
	sender.buff.Reset()
}

//关闭当前日志文件，并获取新的日志文件名称
//清空数据服务器地址，用于重新选择一个服务器
func (sender *Sender4Wego) GetNewFile() string {
	f_next := get_next_log_file(sender.file_path, sender.file_pos.fname)
	if f_next != "" {
		sender.CloseFile()
		sender.file_pos.fname = f_next
		sender.file_pos.offset = 0
		sender.file_pos.flag_end = false
	}
	return f_next
}

//向服务器获取数据服务器地址
func (sender *Sender4Wego) GetDataNodeAddr(url_old string) (string, error) {
	host_old := ""
	if url_old != "" {
		uu, err := url.Parse(url_old)
		if err == nil {
			host_old = url.QueryEscape(uu.Host)
		}
	}

	url := fmt.Sprintf("%s/get_data_node?acc=%s&app=%s&ip=%s&port=%d&err=%s", sender.http_url,
		sender.log_conf.AccId, sender.log_conf.AppId, sender.log_conf.HostIp,  sender.log_conf.HostPort, host_old)
	str, err := HttpGet(url)
	if err != nil {
		loglog.Error(err)
		return "", err
	}

	loglog.Debug("GetDataNode: " + str)
	return str, nil
}

//获取数据服务器的日志文件信息
//服务器返回日志编码，以及文件的大小
func (sender *Sender4Wego) GetSendInfo() (FilePos, error) {
	var fpos FilePos
	url := fmt.Sprintf("%s/get_send_info?acc=%s&app=%s&ip=%s&port=%d", sender.node_url,
		sender.log_conf.AccId, sender.log_conf.AppId, sender.log_conf.HostIp,  sender.log_conf.HostPort)
	str_data, err := HttpGet(url)
	if err != nil {
		if sender.tm_err == 0 {
			sender.tm_err = time.Now().Unix()
		}
		loglog.Error(err)
		return fpos, err
	}

	var info LogInfo
	if len(str_data) > 0 {
		err := json.Unmarshal([]byte(str_data), &info)
		if err != nil {
			loglog.Debug(err)
			return fpos, err
		}
	}

	fpos.fname = info.GetFilePath(sender.file_path)
	fpos.offset = info.FSize
	fpos.flag_end = info.Status
	loglog.Debugf("GetSendInfo: %s", info.String())

	sender.tm_err = 0
	return fpos, nil
}

//发送文件结束标志
func (sender *Sender4Wego) SendFileEnd(fpos *FilePos) error {
	url := fmt.Sprintf("%s/set_send_info?acc=%s&app=%s&ip=%s&port=%d&fcode=%s", sender.node_url,
		sender.log_conf.AccId, sender.log_conf.AppId, sender.log_conf.HostIp, sender.log_conf.HostPort,
		sender.file_pos.GetFileCode())
	loglog.Debugf("SendFileEnd: %s", get_short_name(fpos.fname))
	_, err := HttpGet(url)
	if err != nil {
		if sender.tm_err == 0 {
			sender.tm_err = time.Now().Unix()
		}
		loglog.Error(err)
		return err
	}

	sender.tm_err = 0
	return nil
}

//发送数据
func (sender *Sender4Wego) SendData() error {
	url := fmt.Sprintf("%s/send_data?acc=%s&app=%s&ip=%s&port=%d&fcode=%s", sender.node_url,
		sender.log_conf.AccId, sender.log_conf.AppId, sender.log_conf.HostIp, sender.log_conf.HostPort,
		sender.file_pos.GetFileCode())
	loglog.Debugf("SendData: data_len=%d\n",sender.buff.Buffered())
	req, err := http.NewRequest("Post", url, bytes.NewReader(sender.buff.GetBytes()))
	if err != nil {
		if sender.tm_err == 0 {
			sender.tm_err = time.Now().Unix()
		}
		loglog.Error(err)
		return err
	}

	resp, err := sender.http_client.Do(req)
	if err != nil {
		if sender.tm_err == 0 {
			sender.tm_err = time.Now().Unix()
		}
		loglog.Error(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		loglog.Error("SendData StatusCode != 200")
		return errors.New("StatusCode != 200")
	}

	sender.tm_err = 0
	return nil
}

func (sender *Sender4Wego) BeginSendLoop() {
	tm_sleep := 0
	for {
		if sender.logger.is_closed == true {
			loglog.Debug("send loop return")
			break
		}
		if tm_sleep > 0 {
			select {
			case <-time.After(time.Duration(tm_sleep) * time.Second):
			case <-sender.logger.close_chan:
				loglog.Debug("stop send loop")
				return
			}
		}
		tm_sleep = 0
		//loglog.Debug("send_loop......")

		//向服务器获取数据服务器地址
		if sender.node_url == "" {
			node_url, err := sender.GetDataNodeAddr("")
			if err != nil {
				tm_sleep = tm_send_loop
				continue
			}
			sender.CloseFile()
			sender.node_url = node_url
		}

		//若出错时间超过规定时长，则需要重新获取服务器地址
		if sender.tm_err > 0 {
			tm_dd :=  time.Now().Unix() - sender.tm_err
			if tm_dd > tm_err_change_datanode {
				loglog.Debug("try to change data_node addr")
				node_url, err := sender.GetDataNodeAddr(sender.node_url)
				if err == nil && node_url != sender.node_url {
					sender.CloseFile()
					sender.node_url = node_url
					sender.tm_err = 0
				}
			}
		}

		//向服务器获取文件的发送位置
		if sender.file_pos.fname == "" {
			file_pos, err := sender.GetSendInfo()
			if err != nil {
				tm_sleep = tm_send_loop
				continue
			}
			sender.file_pos = file_pos
		}

		//若服务端没有提供发送位置，或当前文件已经发送完成，则获取第一个文件
		if sender.file_pos.fname == "" || sender.file_pos.flag_end == true {
			if sender.GetNewFile() == "" {
				tm_sleep = tm_send_loop
				continue
			}
		}

		//打开文件，并定位到发送位置
		//若打开文件错误（通常是被删除了），则需要打开一个新的文件
		if sender.file == nil {
			file, err := open_file_by_pos(sender.file_pos)
			if err != nil {
				if sender.GetNewFile() == "" {
					tm_sleep = tm_send_loop
					continue
				}
				tm_sleep = 0
				continue
			}
			sender.file = file
		}

		//若没有读取到数据，检查是否是当前正在写的日志文件，若是当前正在写的日志文件，则sleep
		//若不是当前正在写的日志文件，则给服务器发送结束标志，并开始读取新的日志文件
		//读取新的日志文件时，需要为新的日志文件获取日志数据服务器
		sender.buff.ReadFromFile(sender.file)
		if sender.buff.Buffered() < 1 {
			if sender.IsWritingFile() == true {
				//已经是当前正写日志文件，需要等待新写入的数据
				tm_sleep = tm_send_loop
				continue
			}
			//文件读完，给服务器发送结束标志
			err := sender.SendFileEnd(&sender.file_pos)
			if err != nil {
				tm_sleep = tm_send_loop
				continue
			}
			//开始读取新的日志文件
			if sender.GetNewFile() == "" {
				tm_sleep = tm_send_loop
				continue
			}
			//为新的日志文件获取日志数据服务器
			node_url, err := sender.GetDataNodeAddr("")
			if err == nil {
				sender.node_url = node_url
			}
			tm_sleep = 0
			continue
		}

		//发送数据给服务器
		//若发送失败，为了防止发送状态不同步，需要重新获取服务端的发送信息
		err := sender.SendData()
		if err != nil {
			sender.CloseFile()
			tm_sleep = tm_send_loop
			continue
		} else {
			sender.buff.Reset()
			tm_sleep = 0
			continue
		}
	}
}