package gows

var accept_options AcceptOptions

func init() {
	accept_options.init()
}

func DefaultAcceptOptions() *AcceptOptions {
	return &accept_options
}

func SetOriginCheckFunc(fn OriginCheckFunc) {
	accept_options.checkOrigin = fn
}

func UseFlate(val ...CompressAlloter) {
	accept_options.UseFlate(val...)
}

func SetMinCompressSize(size int) {
	accept_options.SetMinCompressSize(size)
}

func SetFrameReadBuffSize(size int) {
	accept_options.SetFrameReadBuffSize(size)
}

func SetMessageBufferSize(size int) {
	accept_options.SetMessageBufferSize(size)
}
