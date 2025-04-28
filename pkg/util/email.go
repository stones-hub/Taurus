package util

import (
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailInfo struct {
	Username       string        // 61647649@qq.com
	Password       string        // 邮箱授权码
	ConnectTimeout time.Duration // 链接超时时间
	SendTimeout    time.Duration // 发送邮件超时时间
	Host           string        // smtp.example.com ， 邮件服务器地址
	Port           int           // 端口
	KeepAlive      bool          // 是否长链
}

// Content 邮件内容
type EmailContent struct {
	From    string   // 来源 61647649 <61647649@qq.com>  需要跟发送邮件的用户名一致
	Subject string   // 标题
	Body    string   // 内容，目前只支持html解析
	File    []string // 内容带附件
}

func send(content *EmailContent, to string, smtpClient *mail.SMTPClient) error {
	//Create the email message
	email := mail.NewMSG()

	email.SetFrom(content.From).AddTo(to).SetSubject(content.Subject)

	//Get from each mail
	email.GetFrom()
	email.SetBody(mail.TextHTML, content.Body)

	//Send with high priority
	email.SetPriority(mail.PriorityHigh)

	// 判断是否有File来发送附件
	if len(content.File) > 0 {
		for _, file := range content.File {
			filename := file
			parts := strings.Split(file, "/")
			if len(parts) > 0 {
				filename = parts[len(parts)-1]
			}
			email.AddAttachment(file, filename)
		}
	}

	// always check error after send
	if email.Error != nil {
		return email.Error
	}

	//Pass the client to the email message to send it
	return email.Send(smtpClient)
}

func SendMail(content *EmailContent, to string, info EmailInfo) (err error) {
	var (
		emailClient = mail.NewSMTPClient()
		smtpClient  *mail.SMTPClient
	)

	emailClient.Host = info.Host
	emailClient.Port = info.Port
	emailClient.Username = info.Username
	emailClient.Password = info.Password
	emailClient.Encryption = mail.EncryptionNone
	emailClient.ConnectTimeout = info.ConnectTimeout
	emailClient.SendTimeout = info.SendTimeout
	emailClient.KeepAlive = info.KeepAlive

	smtpClient, err = emailClient.Connect()
	if err != nil {
		return
	}

	//NOOP command, optional, usedfor avoid timeout when KeepAlive is true and you aren't sending mails.
	//Execute this command each 30 seconds is ideal for persistent connection
	err = smtpClient.Noop()
	if err != nil {
		return
	}

	err = send(content, to, smtpClient)
	return
}

func SendMultipleEmails(content *EmailContent, toList []string, info EmailInfo) (err error) {

	var (
		emailClient = mail.NewSMTPClient()
		smtpClient  *mail.SMTPClient
	)
	emailClient.Host = info.Host
	emailClient.Port = info.Port
	emailClient.Username = info.Username
	emailClient.Password = info.Password
	emailClient.Encryption = mail.EncryptionNone
	emailClient.ConnectTimeout = info.ConnectTimeout
	emailClient.SendTimeout = info.SendTimeout

	//KeepAlive true because the connection need to be open for multiple emails
	//For avoid inactivity timeout, every 30 second you can send a NO OPERATION command to smtp client
	//use smtpClient.Client.Noop() after 30 second of inactivity in this example
	emailClient.KeepAlive = true

	//For authentication you can use AuthPlain, AuthLogin or AuthCRAMMD5
	emailClient.Authentication = mail.AuthPlain

	smtpClient, err = emailClient.Connect()

	if err != nil {
		return
	}

	// toList := [3]string{"to1@example1.com", "to3@example2.com", "to4@example3.com"}
	for _, to := range toList {
		err = send(content, to, smtpClient)
		if err != nil {
			return
		}
	}
	return
}

func SendMailWithTSL(content *EmailContent, to string, info EmailInfo) (err error) {
	var (
		emailClient = mail.NewSMTPClient()
		smtpClient  *mail.SMTPClient
	)
	emailClient.Host = info.Host
	emailClient.Port = info.Port
	emailClient.Username = info.Username
	emailClient.Password = info.Password
	emailClient.Encryption = mail.EncryptionSTARTTLS
	emailClient.ConnectTimeout = info.ConnectTimeout
	emailClient.SendTimeout = info.SendTimeout
	emailClient.KeepAlive = info.KeepAlive

	smtpClient, err = emailClient.Connect()
	if err != nil {
		return
	}

	err = send(content, to, smtpClient)
	return
}

func SendMailWithSSL(content *EmailContent, to string, info EmailInfo) (err error) {
	var (
		emailClient = mail.NewSMTPClient()
		smtpClient  *mail.SMTPClient
	)
	emailClient.Host = info.Host
	emailClient.Port = info.Port
	emailClient.Username = info.Username
	emailClient.Password = info.Password
	emailClient.Encryption = mail.EncryptionSSLTLS
	emailClient.ConnectTimeout = info.ConnectTimeout
	emailClient.SendTimeout = info.SendTimeout
	emailClient.KeepAlive = info.KeepAlive

	smtpClient, err = emailClient.Connect()
	if err != nil {
		return
	}

	err = send(content, to, smtpClient)
	return
}
