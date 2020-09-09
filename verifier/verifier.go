package verifier

import (
	"crypto/tls"
	"fmt"
	"github.com/txthinking/socks5"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type VerifyOptions struct {
	To            string
	ProxyType     string
	ProxyAddress  string
	ProxyPort     int
	ProxyUsername string
	ProxyPassword string
	toInfo        mailAddressInfo
}

func (r *VerifyOptions) parse() bool {
	s1 := strings.Split(r.To, "@")
	if len(s1) < 2 {
		return false
	}

	r.toInfo.Address = r.To
	r.toInfo.Name = ""
	r.toInfo.Username = s1[0]
	r.toInfo.Domain = s1[1]
	return true
}

func (r VerifyOptions) HasProxy() bool {
	if r.ProxyType == "" {
		return false
	}

	if r.ProxyAddress == "" || r.ProxyPort == 0 {
		return false
	}

	return true
}

type mailAddressInfo struct {
	Name     string
	Address  string
	Username string
	Domain   string
}

type VerifyResult struct {
	MailAddress string
	ValidFormat bool
	Deliverable bool
	HostExists  bool
	CatchAll    bool
	Message     string
}

var (
	ResultEmpty         = VerifyResult{Message: "error with unknown"}
	ResultFormatInvalid = VerifyResult{Message: "email format invalid", CatchAll: true}
	ResultHostInvalid   = VerifyResult{Message: "host has no mx records", ValidFormat: true, CatchAll: true}
)

func checkMXRecord(domain string) bool {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		return false
	}

	if len(mxs) == 0 {
		return false
	}

	return true
}

type Dialer func(network, address string) (net.Conn, error)

func VerifyWithPlain(options VerifyOptions) VerifyResult {
	if !options.parse() {
		return ResultFormatInvalid
	}
	d := net.Dialer{Timeout: 5 * time.Second}
	return verifyWithDialer(d.Dial, options)
}

func verifyWithDialer(dialer Dialer, options VerifyOptions) VerifyResult {
	result := VerifyResult{}
	if !checkMXRecord(options.toInfo.Domain) {
		return ResultHostInvalid
	} else {
		result.HostExists = true
	}
	smtpDomain := fmt.Sprintf("smtp.%s", options.toInfo.Domain)
	conn, err := dialer("tcp", fmt.Sprintf("%s:%d", smtpDomain, 587))
	if err != nil {
		return VerifyResult{Message: err.Error()}
	}
	s, err := smtp.NewClient(conn, options.toInfo.Domain)
	if err != nil {
		return VerifyResult{Message: err.Error()}
	}
	err = s.StartTLS(&tls.Config{ServerName: smtpDomain})
	if err != nil {
		return VerifyResult{Message: err.Error()}
	}
	result.Deliverable = true
	result.CatchAll = true
	return result
}

func VerifyWithSocks5(options VerifyOptions) VerifyResult {
	if !options.parse() {
		return ResultFormatInvalid
	}
	socks, err := socks5.NewClient(fmt.Sprintf("%s:%d", options.ProxyAddress, options.ProxyPort), options.ProxyPassword, options.ProxyPassword, 5, 5)
	if err != nil {
		return VerifyResult{Message: err.Error()}
	}
	return verifyWithDialer(socks.Dial, options)
}

func Verify(options VerifyOptions) VerifyResult {
	if options.HasProxy() {
		return VerifyWithSocks5(options)
	}

	return VerifyWithPlain(options)
}
