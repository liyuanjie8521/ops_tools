package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SmtpHost     string `toml:"smtp_host"`
	SmtpPort     int    `toml:"smtp_port"`
	SmtpUsername string `toml:"smtp_username"`
	SmtpPassword string `toml:"smtp_password"`
	From         string `toml:"from"`
	To           string `toml:"to"`
	Subject      string `toml:"subject"`
	Body         string `toml:"body"`
	Keep_vip     string `toml:"keep_vip"`
}

func main() {
	// 读取配置文件
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tomlFile := filepath.Join(dir, "config.toml")
	content, err := os.ReadFile(tomlFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = toml.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}

	// 2.获取VIP地址
	cmd := exec.Command("ip", "-4", "addr", "show", "label", config.Keep_vip)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	keep_vip := strings.Fields(string(output))[4]
	log.Printf("Detected VIP: %s", keep_vip)

	// 3. 检查keepalived的主备状态
	cmd = exec.Command("systemctl", "is-active", "keepalived")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	status := strings.TrimSpace(string(output))
	fmt.Printf("Keepalived status: %s", status)

	if status == "active" {
		if keep_vip != config.Keep_vip {
			sendEmail(config.SmtpHost, config.SmtpPort, config.SmtpUsername, config.SmtpPassword, config.From, config.To, config.Subject, config.Body)
			fmt.Printf("VIP changed to %s", keep_vip)
		} else { // If VIP hasn't changed, don't send an email. You may need to modify this logic based on your specific needs.
			log.Println("No VIP change detected.")
		}
	} else { // If Keepalived is inactive, don't send an email. You may need to modify this logic based on your specific needs.
		log.Println("Keepalived is inactive.")
	}
}

func sendEmail(SmtpHost string, SmtpPort int, SmtpUsername string, SmtpPassword string, From string, To string, Subject string, Body string) {
	// 配置SMTP认证信息
	auth := smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpHost)
	// 构建邮件内容
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", From, To, Subject, Body)

	// 连接SMTP服务器
	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", SmtpHost, SmtpPort))
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", SmtpHost, SmtpPort))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送认证信息
	if err = conn.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// 发送邮件内容
	w, err := conn.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
	// 发送邮件内容
	// w, err := conn.Write([]byte(message))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Sent %d bytes\n", w)
}
