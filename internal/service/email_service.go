package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	gomail "gopkg.in/mail.v2"
)

type EmailService interface {
	GenerateCode() (string, error)
	SendEmail(inputEmail string, code string) error
	SendEmailResetPassword(inputEmail string, code string) error
}

type EmailServiceImpl struct {
}

func NewEmailService() *EmailServiceImpl {
	return &EmailServiceImpl{}
}

func (e *EmailServiceImpl) GenerateCode() (string, error) {
	max := big.NewInt(1000000)
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", num), nil
}

func (e *EmailServiceImpl) SendEmail(inputEmail string, code string) error {
	log.Println("Sending email to:", inputEmail)
	htmlCorpo := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="pt-BR">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Verificação de Email</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    color: #333;
                    margin: 0;
                    padding: 0;
                    background-color: #f4f4f4;
                }
                .container {
                    max-width: 600px;
                    margin: 0 auto;
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                }
                h1 {
                    color: #4CAF50;
                }
                p {
                    font-size: 16px;
                    line-height: 1.5;
                }
                .code {
                    font-size: 24px;
                    font-weight: bold;
                    color: #4CAF50;
                    background-color: #f0f0f0;
                    padding: 10px;
                    border-radius: 5px;
                    text-align: center;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Verificação de Email</h1>
                <p>Olá! Obrigado por se registrar! Para concluir seu cadastro, use o código de verificação abaixo:</p>
                <div class="code">%s</div>
                <p>Este código expira em 10 minutos. Se você não solicitou esta verificação, ignore este email.</p>
                <p>Atenciosamente,equipe Botany!</p>
            </div>
        </body>
        </html>
    `, code)

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_USER"))
	message.SetHeader("To", inputEmail)
	message.SetHeader("Subject", "Verificação de Email")
	message.SetBody("text/html", htmlCorpo)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("HTML Email sent successfully with a plain-text alternative!")
		return nil
	}
}

func (e *EmailServiceImpl) SendEmailResetPassword(inputEmail, code string) error {
	log.Println("Sending email  reset password to:", inputEmail)
	htmlCorpo := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="pt-BR">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Recuperação de Senha</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    color: #333;
                    margin: 0;
                    padding: 0;
                    background-color: #f4f4f4;
                }
                .container {
                    max-width: 600px;
                    margin: 0 auto;
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                }
                h1 {
                    color: #4CAF50;
                }
                p {
                    font-size: 16px;
                    line-height: 1.5;
                }
                .code {
                    font-size: 24px;
                    font-weight: bold;
                    color: #4CAF50;
                    background-color: #f0f0f0;
                    padding: 10px;
                    border-radius: 5px;
                    text-align: center;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Recuperação de Senha</h1>
                <p>Olá! Você solicitou a recuperação de senha. Use o código abaixo para redefinir sua senha:</p>
                <div class="code">%s</div>
                <p>Este código expira em 10 minutos. Se você não solicitou esta recuperação, ignore este email.</p>
                <p>Atenciosamente,equipe Botany!</p>
            </div>
        </body>
        </html>
    `, code)

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_USER"))
	message.SetHeader("To", inputEmail)
	message.SetHeader("Subject", "Recuperação de Senha")
	message.SetBody("text/html", htmlCorpo)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("HTML Email sent successfully with a plain-text alternative!")
	}
	return nil
}
