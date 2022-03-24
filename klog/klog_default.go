package klog

var defaultEngine *LogEngine = nil
func InitEngine(lpath string, rtype RotateType) *LogEngine {
	Close()
	defaultEngine = NewEngine(lpath, rtype)
	return defaultEngine
}

func GetEngine() *LogEngine {
	return defaultEngine
}

func Close() {
	if defaultEngine != nil {
		defaultEngine.Close()
	}
}

func NewLog(class_name string) *LogRow {
	row := getLineEnt()
	row.out = defaultEngine
	row.TableName(class_name)
	return row
}
