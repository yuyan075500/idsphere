package mail

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"ops-api/config"
	"ops-api/utils"
)

var Email EmailInfo

type EmailInfo struct {
	host     string
	form     string
	password string
	port     int
	dialer   *gomail.Dialer
	msg      *gomail.Message
}

// 获取字符串类型配置
func getConfigString(key string) (string, error) {
	if val, ok := config.Conf.Settings[key].(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("%s configuration is missing or not a string", key)
}

// 获取整数类型配置
func getConfigInt(key string) (int, error) {
	if val, ok := config.Conf.Settings[key].(int); ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s configuration is missing or not an integer", key)
}

// Setup 初始化邮件 Dialer
func (e *EmailInfo) Setup() (*gomail.Dialer, error) {

	mailHost, err := getConfigString("mailAddress")
	if err != nil {
		return nil, err
	}
	mailPort, err := getConfigInt("mailPort")
	if err != nil {
		return nil, err
	}
	mailFrom, err := getConfigString("mailForm")
	if err != nil {
		return nil, err
	}
	mailPassword, err := getConfigString("mailPassword")
	if err != nil {
		return nil, err
	}

	// 密码解密
	str, _ := utils.Decrypt(mailPassword)

	if e.dialer == nil {
		// 实例化邮件连接对象
		host := mailHost
		port := mailPort
		form := mailFrom
		password := str

		// 创建SMTP实例
		e.dialer = gomail.NewDialer(host, port, form, password)

		// 允许跳过不安全的认证
		//e.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return e.dialer, nil
}

// SendMsg 发送邮件
func (e *EmailInfo) SendMsg(to, cc, files []string, subject, body, sendType string) error {

	// 初始化连接验证
	_, err := e.Setup()
	if err != nil {
		return err
	}

	// 创建一个邮件对象
	e.msg = gomail.NewMessage()

	// 设置发件人
	mailFrom, err := getConfigString("mailForm")
	if err != nil {
		return err
	}
	e.msg.SetHeader("From", mailFrom)
	// 设置收件人
	e.msg.SetHeader("To", to...)
	// 设置抄送人
	if len(cc) > 0 {
		e.msg.SetHeader("Cc", cc...)
	}

	// 设置邮件标题与正文
	if subject == "" {
		return errors.New("邮件主题不能为空")
	}
	if body == "" {
		return errors.New("邮件内容不能为空")
	}

	// 设置邮件主题
	e.msg.SetHeader("Subject", subject)

	// 添加附件
	if len(files) > 0 {
		for _, file := range files {
			e.msg.Attach(file)
		}
	}

	// 判断发送邮件的类型并设置正文
	if sendType == "text" {
		e.msg.SetBody("text/plain", body)
	} else if sendType == "html" {
		e.msg.SetBody("text/html", body)
	}

	// 发送邮件
	if err := e.dialer.DialAndSend(e.msg); err != nil {
		return err
	}

	// 关闭连接
	e.msg.Reset()

	return nil
}
