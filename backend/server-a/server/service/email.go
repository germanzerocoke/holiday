package service

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"math/rand"
	"net/smtp"
	"os"
	"server-a/server/dto"
	"strconv"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CheckEmailUsable(ctx context.Context, email string) (bool, error) {
	isTaken, err := s.repository.EmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	if isTaken {
		slog.Info("this email already exist",
			"email", email,
		)
		return false, nil
	}
	return true, nil
}

func (s *Service) CreateMemberByEmail(ctx context.Context, email, password string) (map[string]string, error) {
	exist, err := s.repository.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if exist {
		log.Printf("this email already exist")
		return nil, errors.New("this email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("fail to hash password: %v", err)
		return nil, err
	}

	password = string(hashedPassword)
	idv7, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for email password sign in user")
		return nil, err
	}
	id := gocql.UUID(idv7)
	err = s.repository.SaveEmailLoginInfo(id, email, password)
	if err != nil {
		return nil, err
	}

	return map[string]string{"id": id.String()}, nil
}

func (s *Service) LoginWithEmail(email, password string) (*dto.LoginWithEmailResponse, string /*refreshToken*/, error) {
	var resp dto.LoginWithEmailResponse

	emailVerified, phoneNumberVerified, id, dbPassword, role, err :=
		s.repository.FindLoginInfoByEmail(email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		slog.Info("invalid password",
			"err", err,
		)
		return nil, "", err
	}

	if !emailVerified {
		resp.EmailVerified = false
		resp.PhoneNumberVerified = false
		resp.Id = id.String()

		return &resp, "", nil
	}

	if !phoneNumberVerified {
		sid, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to make random uuid for session id")
			return nil, "", err
		}

		err = s.repository.SaveEmailBySessionId(sid, email)
		resp.EmailVerified = true
		resp.PhoneNumberVerified = false
		resp.SessionId = sid.String()
		return &resp, "", nil
	}

	jti, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to create random uuid for jti")
		return nil, "", err
	}
	at, rt, err := s.createLoginTokens(id.String(), jti.String(), role)
	if err != nil {
		return nil, "", err
	}
	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", err
	}
	resp.EmailVerified = true
	resp.PhoneNumberVerified = true
	resp.AccessToken = at
	return &resp, rt, nil
}

func (s *Service) SendEmailOTP(ctx context.Context, id string) (*dto.SendOTPResponse, error) {
	uid, err := gocql.ParseUUID(id)
	if err != nil {
		slog.Error("fail to parse id",
			"err", err,
			"id", id)
		return nil, err
	}
	email, err := s.repository.FindEmailById(uid)
	if err != nil {
		return nil, err
	}
	otp := strconv.Itoa(rand.Intn(900000) + 100000)
	vid, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for verification id")
		return nil, err
	}
	err = s.repository.SaveEmailAndOtpByVerificationId(vid, email, otp)
	if err != nil {
		return nil, err
	}
	go func() {
		from := os.Getenv("FROM_EMAIL")
		auth := smtp.PlainAuth(
			"",
			from,
			os.Getenv("FROM_EMAIL_PASSWORD"),
			os.Getenv("FROM_EMAIL_SMTP"),
		)

		headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
		message := "Subject: Verify your email\n" + headers + "\n\n" + otp + "\ncode is valid for 5 minutes"

		err = smtp.SendMail(
			os.Getenv("SMTP_ADDR"),
			auth,
			from,
			[]string{email},
			[]byte(message),
		)
		if err != nil {
			log.Printf("fail to send email: %v", err)
		}
	}()

	return &dto.SendOTPResponse{VerificationId: vid.String()}, nil
}

func (s *Service) VerifyEmailOTP(otp, verificationId string) (*dto.VerifyEmailOTPResponse, error) {
	vid, err := gocql.ParseUUID(verificationId)
	if err != nil {
		slog.Info("fail to parse uuid from verificationId in req", err)
		return nil, err
	}
	email, dbOTP, err := s.repository.FindEmailAndOTPByVerificationId(vid)
	if err != nil {
		return nil, err
	}
	if otp != dbOTP {
		log.Printf(
			"code is not same with db code- received code: %v, db code: %v",
			otp, dbOTP,
		)
		resp := dto.VerifyEmailOTPResponse{
			EmailVerified: false,
		}
		return &resp, nil
	}

	err = s.repository.MarkEmailVerified(email)
	if err != nil {
		return nil, err
	}

	sid, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for session id")
		return nil, err
	}
	err = s.repository.SaveEmailBySessionId(sid, email)
	if err != nil {
		return nil, err
	}

	resp := dto.VerifyEmailOTPResponse{
		EmailVerified: true,
		SessionId:     sid.String(),
	}
	return &resp, nil
}
