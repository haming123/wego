package klog

import (
	"fmt"
)

type LogInfo struct {
	AccId 		string	`json:"acc_id"`
	AppId 		string	`json:"app_id"`
	CHost 		string	`json:"host_ip"`
	CPort 		int		`json:"host_port"`
	FCode 		string	`json:"file_code"`
	FSize 		int64	`json:"file_size"`
	Status 		bool	`json:"status"`
	Data		[]byte	`json:"data"`
}

func (info *LogInfo) GetFileName() string {
	if info.FCode == "" {
		return ""
	}
	return fmt.Sprintf("klog_%s.log", info.FCode)
}

func (info *LogInfo) GetFilePath(file_path string) string {
	if info.FCode == "" {
		return ""
	}
	return file_path + info.GetFileName()
}

func (info *LogInfo) String() string {
	if info.FCode == "" {
		return ""
	}
	return fmt.Sprintf("file=%s_%s_%s_%d_%s;len=%d;over=%v",
		info.AccId, info.AppId, info.CHost, info.CPort, info.FCode, info.FSize, info.Status)
}