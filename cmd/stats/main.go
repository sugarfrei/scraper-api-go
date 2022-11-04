package main

import (
	"flag"

	"github/sugarfrei/stats-api-go/api"

	"k8s.io/klog"
)

func main() {
	var confpath = flag.String("conf", "./etc/conf.yaml", "Path to a config file")

	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("log_file", "./logs/stats.log")
	flag.Parse()

	a, err := api.New(*confpath)
	if err != nil {
		klog.Fatalf("Unable to init stats-api, err: %s", err.Error())
	}

	a.ListenAndServe()
}
