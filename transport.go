package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func NewTransport(transporter *http.Transport, setting *TransportSetting) http.RoundTripper {
	if transporter.TLSClientConfig == nil {
		transporter.TLSClientConfig = &tls.Config{}
	}

	if setting.IsInsecure {
		transporter.TLSClientConfig.InsecureSkipVerify = true
	}
	if setting.Hostname != "" {
		transporter.TLSClientConfig.ServerName = setting.Hostname
	}
	if setting.IsDebug {
		log.Printf("%+v\n", transporter.TLSClientConfig)
	}

	return transporter
}

type TransportSetting struct {
	Hostname   string
	IsInsecure bool
	IsDebug    bool
}
