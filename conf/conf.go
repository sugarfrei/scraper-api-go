package conf

type CoreCfg struct {
	PrvKey    []byte
	PubKey    []byte
	AccessLog string
}

type HttpCfg struct {
	Listen            string
	ReadTimeout       int
	ReadHeaderTimeout int
	WriteTimeout      int
	IdleTimeout       int
	MaxHeaderBytes    int
}

type ErrorCfg struct {
	Log int
}

type Cfg struct {
	Core  CoreCfg
	Http  HttpCfg
	Error ErrorCfg
}
