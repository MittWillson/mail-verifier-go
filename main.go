package main

import (
	"fmt"
	"github.com/mittwillson/mail-verifier-go/verifier"
	"log"
)

func main() {
	//socks, err := socks5.NewClient("127.0.0.1:1086", "", "", 5, 5)
	//if err != nil {
	//	printError(err)
	//	return
	//}
	//conn, err := socks.Dial("tcp", "smtp.qq.com:587")
	//if err != nil {
	//	printError(err)
	//	return
	//}
	//client, err := smtp.NewClient(conn, "qq.com")
	//if err != nil {
	//	printError(err)
	//	return
	//}
	//err = client.StartTLS(&tls.Config{ServerName: "smtp.qq.com"})
	//printError(err)

	fmt.Printf("%+v\n", verifier.Verify(verifier.VerifyOptions{
		To:            "757549561@qq.com",
		ProxyType:     "socks5",
		ProxyAddress:  "127.0.0.1",
		ProxyPort:     1086,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
	fmt.Printf("%+v\n", verifier.Verify(verifier.VerifyOptions{
		To:            "mitt.willson@gmail.com",
		ProxyType:     "",
		ProxyAddress:  "",
		ProxyPort:     0,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
	fmt.Printf("%+v\n", verifier.Verify(verifier.VerifyOptions{
		To:            "mitt.willson@gmail.com",
		ProxyType:     "socks5",
		ProxyAddress:  "127.0.0.1",
		ProxyPort:     1086,
		ProxyUsername: "",
		ProxyPassword: "",
	}))
}

func printError(err error) {
	log.Fatalf("%+v", err)
}
