package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func NewReverseProxy(setting *ReverseProxySetting) *ReverseProxy {
	var cookies []*http.Cookie

	if setting.Environment != DEFAULT_ENV_NAME {
		cookies = []*http.Cookie{
			{
				Name:  "environment",
				Value: setting.Environment,
			}, {
				Name:  setting.Environment,
				Value: "true",
			},
		}
	}

	return &ReverseProxy{
		setting: setting,
		proxy: &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(setting.Target)
				pr.Out.Host = setting.Hostname
				for _, cookie := range cookies {
					pr.Out.AddCookie(cookie)
				}

				if setting.IsDebug {
					log.Printf("[%s] Incoming request: %+v\n", setting.Environment, pr.In)
					log.Printf("[%s] Proxy request: %+v\n", setting.Environment, pr.Out)
				} else {
					log.Printf("[%s] Forwarding '%s' to '%s'\n",
						setting.Environment,
						pr.In.URL.String(),
						pr.Out.URL.String(),
					)
				}
			},
			Transport: NewTransport(http.DefaultTransport.(*http.Transport), &TransportSetting{
				Hostname:   setting.Hostname,
				IsInsecure: setting.IsInsecure,
				IsDebug:    setting.IsDebug,
			}),
		},
	}
}

type ReverseProxySetting struct {
	Enabled     bool
	Environment string
	ListenHost  string
	ListenPort  int
	Target      *url.URL
	Hostname    string
	IsDebug     bool
	IsInsecure  bool
}

type ReverseProxy struct {
	proxy   *httputil.ReverseProxy
	setting *ReverseProxySetting
}

func (rp *ReverseProxy) Start(ctx context.Context, cancel context.CancelFunc) {
	if rp.setting.Enabled {
		var addr = net.JoinHostPort(rp.setting.ListenHost, strconv.Itoa(rp.setting.ListenPort))
		var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rp.proxy.ServeHTTP(w, r)
		})

		server := &http.Server{Addr: addr, Handler: handler, BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, rp.setting, l.Addr().String())
			return ctx
		}}

		// Start the server
		fmt.Printf("Starting reverse proxy server on http://%s, forwarding to %s (%s)\n\n",
			addr,
			rp.setting.Target.Host,
			rp.setting.Environment,
		)
		var err = server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Reverse proxy %s: closed\n", rp.setting.Environment)
		} else if err != nil {
			fmt.Printf("Reverse proxy %s: error (%v)\n", rp.setting.Environment, err)
		}

		cancel()
	}
}
