package main

import (
	"context"
	"flag"
	"log"
)

const (
	DEFAULT_ENV_NAME = "<default>"
)

var (
	hostname     = flag.String("host", "", "Value of the Host header to send to the next hop server")
	port         = flag.Int("port", 8000, "Listen port number for default environment")
	envEnable    = flag.Bool("env", false, "Enabled proxy to cloudflare per environment")
	devEnable    = flag.Bool("dev", true, "Enabled proxy to cloudflare development environment")
	devPort      = flag.Int("dev-port", 8001, "Listen port number for development")
	stgEnable    = flag.Bool("stg", true, "Enabled proxy to cloudflare staging environment")
	stgPort      = flag.Int("stg-port", 8002, "Listen port number for staging")
	prdEnable    = flag.Bool("prd", true, "Enabled proxy to cloudflare production environment")
	prdPort      = flag.Int("prd-port", 8003, "Listen port number for production")
	flagDebug    = flag.Bool("debug", false, "Enable debug mode")
	flagInsecure = flag.Bool("insecure", false, "Ignore SSL certificate errors")
)

func init() {
	// Parse the flags
	flag.Parse()
	if *flagDebug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	var target = BuildCloudflareCDN(*hostname)

	var proxy = NewReverseProxy(&ReverseProxySetting{
		Enabled:     !*envEnable,
		Environment: DEFAULT_ENV_NAME,
		ListenHost:  "localhost",
		ListenPort:  *port,
		Target:      target,
		Hostname:    *hostname,
		IsDebug:     *flagDebug,
		IsInsecure:  *flagInsecure,
	})

	var devProxy = NewReverseProxy(&ReverseProxySetting{
		Enabled:     *envEnable && *devEnable,
		Environment: "development",
		ListenHost:  "localhost",
		ListenPort:  *devPort,
		Target:      target,
		Hostname:    *hostname,
		IsDebug:     *flagDebug,
		IsInsecure:  *flagInsecure,
	})

	var stgProxy = NewReverseProxy(&ReverseProxySetting{
		Enabled:     *envEnable && *stgEnable,
		Environment: "staging",
		ListenHost:  "localhost",
		ListenPort:  *stgPort,
		Target:      target,
		Hostname:    *hostname,
		IsDebug:     *flagDebug,
		IsInsecure:  *flagInsecure,
	})

	var prdProxy = NewReverseProxy(&ReverseProxySetting{
		Enabled:     *envEnable && *prdEnable,
		Environment: "production",
		ListenHost:  "localhost",
		ListenPort:  *prdPort,
		Target:      target,
		Hostname:    *hostname,
		IsDebug:     *flagDebug,
		IsInsecure:  *flagInsecure,
	})

	go proxy.Start(ctx, cancel)
	go devProxy.Start(ctx, cancel)
	go stgProxy.Start(ctx, cancel)
	go prdProxy.Start(ctx, cancel)

	if *envEnable && !*devEnable && !*stgEnable && !*prdEnable {
		cancel()
	}

	<-ctx.Done()
}
