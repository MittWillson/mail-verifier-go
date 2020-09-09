package main

import (
	"flag"
	"fmt"
	"github.com/mittwillson/mail-verifier-go/verifier"
	"net"
	"strconv"
)

var (
	toPtr            = flag.String("to", "random@qq.com", "测试邮箱地址")
	proxyAddrPtr     = flag.String("proxy-address", "", "可选 Socks5代理地址 格式: IP:PORT")
	proxyUsernamePtr = flag.String("proxy-username", "", "可选 Socks5代理用户名")
	proxyPasswordPtr = flag.String("proxy-password", "", "可选 Socks5代理密码")
)

func main() {
	flag.Parse()

	var (
		proxyAddr = ""
		proxyPort = 0
	)

	if *proxyAddrPtr != "" {
		var p string
		var err error
		proxyAddr, p, err = net.SplitHostPort(*proxyAddrPtr)
		if err != nil {
			fmt.Printf("代理地址格式有误: %s\n", err.Error())
			return
		}

		proxyPort, _ = strconv.Atoi(p)
		if proxyPort == 0 {
			fmt.Println("代理端口有误")
			return
		}
	}

	options := verifier.VerifyOptions{
		To:            *toPtr,
		ProxyType:     "socks5",
		ProxyAddress:  proxyAddr,
		ProxyPort:     proxyPort,
		ProxyUsername: *proxyUsernamePtr,
		ProxyPassword: *proxyPasswordPtr,
	}
	result := verifier.Verify(options)

	if !result.CatchAll {
		printErrMsg(result.Message)
		printResult(false, options.HasProxy())
		return
	}

	if !result.ValidFormat {
		printErrMsg(result.Message)
		fmt.Printf("邮箱格式有误: %s\n", *toPtr)
		return
	}

	if !result.HostExists {
		printErrMsg(result.Message)
		fmt.Printf("邮箱服务器不存在\n")
		return
	}

	if !result.Deliverable {
		printErrMsg(result.Message)
		printResult(false, options.HasProxy())
		return
	}

	printResult(true, options.HasProxy())
}

func printErrMsg(msg string) {
	if msg != "" {
		fmt.Printf("错误信息: %s\n", msg)
	}
}

func printResult(succeed bool, hasProxy bool) {
	if !succeed && hasProxy {
		fmt.Println("【测试结果】测试邮箱连接失败，代理无效(端口可能被封禁)")
		return
	}

	if !succeed && !hasProxy {
		fmt.Println("【测试结果】测试邮箱连接失败，直连不通")
		return
	}

	if !hasProxy {
		fmt.Println("【测试结果】测试邮箱连接成功, 直连有效")
		return
	}

	fmt.Println("【测试结果】测试邮箱连接成功, 代理有效")
}
