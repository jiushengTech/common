package client

import (
	"github.com/jiushengTech/common/log"
	"testing"
)

var (
	TokenUrl               = "http://publish.sense-ecology.com/manager/api/openAPI/userProLogin"
	XSenseAppSessionHeader = "x-sense-app-session"
	XSenseAppVersionHeader = "x-sense-app-version"
	XSenseProSessionHeader = "x-sense-pro-session"
	XSenseProVersionHeader = "x-sense-pro-version"
	XSenseAppSessionValue  = "1da746d14e27483bad3312444c56a322"
	XSenseAppVersionValue  = "1.0"
	XSenseProSessionValue  = "f480d105739a40999f68637a1e30dacf"
	XSenseProVersionValue  = "1.0"
	//  body
	ProKeyBody    = "7f3882eb8cd34fe9a8e98fc17e13cb77"
	ProSecretBody = "0b5705f2c6f94b89aa857d8c6f97ba48"
)

type SessionResp struct {
	ProSession  string `json:"proSession"`
	MessageCode int    `json:"messageCode"`
	Message     string `json:"message"`
	MessageType int    `json:"messageType"`
}

func TestNewHttpClient(t *testing.T) {
	httpClient := NewHttpClient(
		WithHeader("Content-Type", "application/json"),
		WithHeader(XSenseAppSessionHeader, XSenseAppSessionValue),
		WithHeader(XSenseAppVersionHeader, XSenseAppVersionValue),
		WithHeader(XSenseProSessionHeader, XSenseProSessionValue),
		WithHeader(XSenseProVersionHeader, XSenseProVersionValue),
		WithBody("ProKey", ProKeyBody),
		WithBody("ProSecret", ProSecretBody),
	)
	res := &SessionResp{}
	err := httpClient.PostJSONWithResp(TokenUrl, res)
	log.Infof("res: %+v", *res)
	if err != nil {
		log.Error(err)
	}
}

func TestDev(t *testing.T) {
	httpClient := NewHttpClient(
		WithHeader("Content-Type", "application/json"),
	)
	data, err := httpClient.Get("http://192.168.10.119:8000/yw/system/sys_dict_type")
	if err != nil {
		log.Error("HTTP GET error:", err)
		return
	}
	// 将字节数组转换为字符串
	s := string(data)
	// 打印结果
	log.Info("Response:", s)
}
