package main

import (
	"flag"
	"io/ioutil"

	"scraper-api-go/api"
	"scraper-api-go/conf"

	"k8s.io/klog"
)

var (
	privateKey        = flag.String("privateKey", "./cmd/scraper/etc/private-key.pem", "private pem")
	publicKey         = flag.String("publicKey", "./cmd/scraper/etc/public-key.pem", "public pem")
	log               = flag.String("log", "./cmd/scraper/logs/scraper.log", "log path")
	accessLog         = flag.String("accessLog", "./cmd/scraper/logs/access.log", "access log file path")
	errorLog          = flag.Int("errorLog", 2, "0 - disable error log, 1 - log only 5xx errors, 2 - log all errors")
	listen            = flag.String("listen", "0.0.0.0:8080", "address for the server to listen on")
	readTimeout       = flag.Int("readTimeout", 60, "maximum duration for reading the entire request")
	readHeaderTimeout = flag.Int("readHeaderTimeout", 2, "amount of time allowed to read request headers")
	writeTimeout      = flag.Int("writeTimeout", 60, "maximum duration before timing out writes of the response")
	idleTimeout       = flag.Int("idleTimeout", 60, "maximum amount of time to wait for the next request when keep-alives are enabled")
	maxHeaderBytes    = flag.Int("maxHeaderBytes", 4096, "maximum number of bytes the server will read parsing the request header's keys and values, including the request line")
)

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("log_file", *log)
	flag.Parse()

	prvKey, err := ioutil.ReadFile(*privateKey)
	if err != nil {
		klog.Fatalf("%v", err)
	}
	pubKey, err := ioutil.ReadFile(*publicKey)
	if err != nil {
		klog.Fatalf("%v", err)
	}

	cfg := &conf.Cfg{
		Core: conf.CoreCfg{
			PrvKey:    prvKey,
			PubKey:    pubKey,
			AccessLog: *accessLog,
		},
		Error: conf.ErrorCfg{
			Log: *errorLog,
		},
		Http: conf.HttpCfg{
			Listen:            *listen,
			ReadTimeout:       *readTimeout,
			ReadHeaderTimeout: *readHeaderTimeout,
			WriteTimeout:      *writeTimeout,
			IdleTimeout:       *idleTimeout,
			MaxHeaderBytes:    *maxHeaderBytes,
		},
	}

	a, err := api.New(cfg)
	if err != nil {
		klog.Fatalf("Unable to init scraper-api, err: %s", err.Error())
	}

	a.ListenAndServe()
}
