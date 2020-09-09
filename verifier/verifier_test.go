package verifier

import (
	"testing"
)

func TestQQ(t *testing.T) {
	t.Logf("%+v\n", Verify(VerifyOptions{
		To:            "test@qq.com",
		ProxyType:     "",
		ProxyAddress:  "",
		ProxyPort:     0,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
}

func TestQQWithLocalProxy(t *testing.T) {
	t.Logf("%+v\n", Verify(VerifyOptions{
		To:            "test@qq.com",
		ProxyType:     "socks5",
		ProxyAddress:  "127.0.0.1",
		ProxyPort:     1086,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
}

func TestGmail(t *testing.T) {
	t.Logf("%+v\n", Verify(VerifyOptions{
		To:            "test@gmail.com",
		ProxyType:     "",
		ProxyAddress:  "",
		ProxyPort:     0,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
}

func TestGmailWithLocalProxy(t *testing.T) {
	t.Logf("%+v\n", Verify(VerifyOptions{
		To:            "test@gmail.com",
		ProxyType:     "socks5",
		ProxyAddress:  "127.0.0.1",
		ProxyPort:     1086,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
}
