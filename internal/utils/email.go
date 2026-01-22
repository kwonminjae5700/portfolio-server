package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"portfolio-server/internal/config"
)

// GenerateVerificationCode generates a 6-digit verification code
func GenerateVerificationCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code, nil
}

// SendVerificationEmail sends a verification code to the user's email
func SendVerificationEmail(email, code string) error {
	cfg := config.LoadConfig()

	// SMTP 설정
	from := cfg.SMTP.From
	password := cfg.SMTP.Password
	smtpHost := cfg.SMTP.Host
	smtpPort := cfg.SMTP.Port

	// 이메일 내용
	subject := "이메일 인증 코드"
	body := `<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>이메일 인증</title>
</head>
<body style="margin: 0; padding: 0; background-color: #fafafa; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;">
    <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #fafafa;">
        <tr>
            <td style="padding: 40px 20px;">
                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width: 560px; margin: 0 auto; background-color: #ffffff;">
                    <tr>
                        <td style="padding: 48px 40px 32px 40px;">
                            <h1 style="margin: 0 0 8px 0; color: #434a53; font-size: 24px; font-weight: 600;">이메일 인증</h1>
                            <p style="margin: 0; color: #999; font-size: 14px;">회원가입 인증 코드입니다.</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 0 40px 40px 40px;">
                            <p style="margin: 0 0 24px 0; color: #666; font-size: 15px; line-height: 1.6;">
                                회원가입을 완료하려면 아래 인증 코드를 입력해주세요.
                            </p>
                            <table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 0 24px 0;">
                                <tr>
                                    <td style="background-color: #f8f8f8; padding: 24px; text-align: center; border: 1px solid #e0e0e0;">
                                        <div style="font-size: 36px; font-weight: 600; color: #3F35FF; letter-spacing: 6px; font-family: 'Courier New', monospace;">
                                            ` + code + `
                                        </div>
                                    </td>
                                </tr>
                            </table>
                            <p style="margin: 0 0 8px 0; color: #999; font-size: 13px; line-height: 1.5;">
                                유효 시간: 10분
                            </p>
                            <p style="margin: 0; color: #999; font-size: 13px; line-height: 1.5;">
                                본인이 요청하지 않은 경우 이 메일을 무시하세요.
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 24px 40px; background-color: #f8f8f8; border-top: 1px solid #e0e0e0;">
                            <p style="margin: 0; color: #999; font-size: 12px; text-align: center;">
                                이 메일은 자동 발송되었습니다.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

	// 메시지 구성
	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			body + "\r\n")

	// SMTP 인증
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// TLS 설정
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	// 이메일 전송 (TLS 사용)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	
	// STARTTLS를 사용한 SMTP 연결
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer client.Close()

	// STARTTLS 시작
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// 인증
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// 발신자 설정
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// 수신자 설정
	if err = client.Rcpt(email); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// 메시지 전송
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	client.Quit()

	return nil
}
