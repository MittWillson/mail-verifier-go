package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/mittwillson/mail-verifier-go/verifier"
	"net"
	"strconv"
	"strings"
)

var resultOutput = ui.NewMultilineEntry()

func main() {
	err := ui.Main(func() {
		inputEmail := ui.NewEntry()
		inputEmail.SetText("")

		inputProxy := ui.NewEntry()
		inputProxy.SetText("")

		buttonTest := ui.NewButton("测试")

		resultOutput.SetReadOnly(true)

		//

		div := ui.NewHorizontalBox()

		leftDiv := ui.NewVerticalBox()

		leftTop1 := ui.NewVerticalBox()
		leftTop1.Append(ui.NewLabel("测试邮箱:"), false)
		leftTop1.Append(inputEmail, false)
		leftTop1.Append(ui.NewLabel("测试Socks5代理 【留空则直连测试】:"), false)
		leftTop1.Append(inputProxy, false)
		leftTop1.Append(buttonTest, false)

		leftBottom1 := ui.NewVerticalBox()
		leftBottom1.Append(resultOutput, true)

		leftDiv.Append(leftTop1, true)
		leftDiv.Append(leftBottom1, true)
		leftDiv.SetPadded(false)

		div.Append(leftDiv, false)
		div.SetPadded(true)

		//
		window := ui.NewWindow("SMTP邮箱代理测试软件 made by Mitt", 150, 200, false)
		window.SetChild(div)
		window.SetMargined(true)

		buttonTest.OnClicked(func(*ui.Button) {
			var (
				proxyAddr     = ""
				proxyPort     = 0
				proxyUsername = ""
				proxyPassword = ""
				inputEmail    = strings.TrimSpace(inputEmail.Text())
				inputProxy    = strings.TrimSpace(inputProxy.Text())
			)

			if inputProxy != "" {
				proxies := strings.Split(inputProxy, ":")
				if len(proxies) == 4 {
					inputProxy = fmt.Sprintf("%s:%s", proxies[0], proxies[1])
					proxyUsername = strings.TrimSpace(proxies[2])
					proxyPassword = strings.TrimSpace(proxies[3])
				}

				var p string
				var err error
				proxyAddr, p, err = net.SplitHostPort(inputProxy)
				if err != nil {
					SetResultLabelText(fmt.Sprintf("代理地址格式有误: %s\n", err.Error()))
					return
				}

				proxyPort, _ = strconv.Atoi(p)
				if proxyPort == 0 {
					SetResultLabelText(fmt.Sprint("代理端口有误"))
					return
				}
			}

			options := verifier.VerifyOptions{
				To:            inputEmail,
				ProxyType:     "socks5",
				ProxyAddress:  proxyAddr,
				ProxyPort:     proxyPort,
				ProxyUsername: proxyUsername,
				ProxyPassword: proxyPassword,
			}

			result := verifier.Verify(options)

			if !result.CatchAll {
				printErrMsg(result.Message)
				printResult(false, options.HasProxy())
				return
			}

			if !result.ValidFormat {
				printErrMsg(result.Message)
				SetResultLabelText(fmt.Sprintf("邮箱格式有误: %s\n", inputEmail))
				return
			}

			if !result.HostExists {
				printErrMsg(result.Message)
				SetResultLabelText(fmt.Sprintf("邮箱服务器不存在\n"))
				return
			}

			if !result.Deliverable {
				printErrMsg(result.Message)
				printResult(false, options.HasProxy())
				return
			}

			printResult(true, options.HasProxy())
		})

		window.OnClosing(func(_ *ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})

	if err != nil {
		panic(err)
	}
}

func printErrMsg(msg string) {
	if msg != "" {
		SetResultLabelText(fmt.Sprintf("错误信息: %s\n", msg))
	}
}

func printResult(succeed bool, hasProxy bool) {
	if !succeed && hasProxy {
		SetResultLabelText(fmt.Sprintln("【测试结果】测试邮箱连接失败，代理无效(端口可能被封禁)"))
		return
	}

	if !succeed && !hasProxy {
		SetResultLabelText(fmt.Sprintln("【测试结果】测试邮箱连接失败，直连不通"))
		return
	}

	if !hasProxy {
		SetResultLabelText(fmt.Sprintln("【测试结果】测试邮箱连接成功, 直连有效"))
		return
	}

	SetResultLabelText(fmt.Sprintln("【测试结果】测试邮箱连接成功, 代理有效"))
}

func SetResultLabelText(msg string) {
	resultOutput.SetText(msg)
}
