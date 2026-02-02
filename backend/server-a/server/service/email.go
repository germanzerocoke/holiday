package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"server-a/server/constant"
	"server-a/server/constant/message"
	"server-a/server/dto"
	"strconv"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) IsEmailUsable(ctx context.Context, email string) (bool, error) {
	i, err := s.repository.EmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	if i {
		log.Printf("email already exist")
		return false, nil
	}
	return true, nil
}

func (s *Service) CreateMemberByEmail(ctx context.Context, email, password string) (map[string]string, error) {
	i, err := s.repository.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if i {
		log.Printf("this email already exist")
		return nil, errors.New("this email already exist")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("fail to hash password: %v", err)
		return nil, err
	}
	password = string(hashedPassword)
	id := gocql.TimeUUID()

	err = s.repository.SaveEmailMember(id, email, password)
	if err != nil {
		return nil, err
	}
	return map[string]string{"id": id.String()}, nil
}

func (s *Service) LoginWithEmail(email, password string) (*dto.EmailLoginResponse, string /*refreshToken*/, error) {
	var resp dto.EmailLoginResponse
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
		resp.EmailVerified = emailVerified
		resp.PhoneNumberVerified = phoneNumberVerified
		resp.Id = id.String()

		return &resp, "", nil
	}
	if !phoneNumberVerified {
		sid := gocql.TimeUUID()
		err = s.repository.SaveEmailBySessionId(sid, email)
		resp.EmailVerified = emailVerified
		resp.PhoneNumberVerified = phoneNumberVerified
		resp.SessionId = sid.String()
		return &resp, "", nil
	}
	at, rt, err := s.createLoginTokens(id.String(), role)
	if err != nil {
		return nil, "", err
	}
	err = s.repository.SaveRefreshTokenById(id, rt)
	if err != nil {
		return nil, "", err
	}
	resp.EmailVerified = emailVerified
	resp.PhoneNumberVerified = phoneNumberVerified
	resp.AccessToken = at
	return &resp, rt, nil
}

func (s *Service) SendEmailOTP(ctx context.Context, id string) (*dto.OTPSendResponse, error) {
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
	vid := gocql.TimeUUID()
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

	return &dto.OTPSendResponse{VerificationId: vid.String()}, nil
}

func (s *Service) VerifyEmailOTP(otp, verificationId string) (*dto.EmailOTPVerifyResponse, error) {
	vid, err := gocql.ParseUUID(verificationId)
	if err != nil {
		slog.Info("fail to parse uuid from verificationId in req", err)
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
		resp := dto.EmailOTPVerifyResponse{
			EmailVerified: false,
		}
		return &resp, nil
	}

	err = s.repository.MarkEmailVerified(email)
	if err != nil {
		return nil, err
	}

	sid := gocql.TimeUUID()
	err = s.repository.SaveEmailBySessionId(sid, email)
	if err != nil {
		return nil, err
	}

	resp := dto.EmailOTPVerifyResponse{
		EmailVerified: true,
		SessionId:     sid.String(),
	}
	return &resp, nil
}

func (s *Service) SignInWithApple(
	ctx context.Context,
	user,
	rawNonce string,
	email,
	identityToken string,
) (*dto.SignInWithAppleResponse, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		constant.AppleKeyUrl,
		nil,
	)
	if err != nil {
		slog.Error("fail to make http request",
			"err", err,
		)
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("fail to do request", "err", err)
		return nil, err
	}
	defer resp.Body.Close()
	var jwks map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		slog.Error("fail to decode jwks in time",
			"err", err,
			"resp", resp,
		)
		return nil, err
	}

	idt, err := jwt.Parse(identityToken, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodES256.Alg() {
			slog.Info("unexpected signing method")
			return nil, errors.New(message.AppleSignInFailed)
		}
		return jwks, nil
	})
	issFromClaims, err := idt.Claims.GetIssuer()
	if err != nil {
		slog.Info("fail to get issuer",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, err
	}
	if issFromClaims != constant.AppleIssuerUrl {
		slog.Info("not expected apple issuer",
			"iss", issFromClaims,
		)
		return nil, errors.New(message.AppleSignInFailed)
	}

	audsFromClaims, err := idt.Claims.GetAudience()
	if err != nil {
		slog.Info("fail to get audience",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, errors.New(message.AppleSignInFailed)
	}
	if len(audsFromClaims) != 0 || audsFromClaims[0] != s.audience {
		slog.Info("no audience or not expected audience")
		return nil, errors.New(message.AppleSignInFailed)
	}

	exp, err := idt.Claims.GetExpirationTime()
	if err != nil {
		slog.Info("fail to get expiration time")
		return nil, errors.New(message.AppleSignInFailed)
	}
	if exp.Unix() < time.Now().Unix() {
		slog.Info("stale token")
		return nil, errors.New(message.AppleSignInFailed)
	}

	nonceFromClaims, ok := idt.Claims.(jwt.MapClaims)["nonce"].(string)
	if !ok {
		slog.Info("no nonce in claims")
		return nil, errors.New(message.AppleSignInFailed)
	}
	sum := sha256.Sum256([]byte(rawNonce))
	hashedNonce := base64.RawURLEncoding.EncodeToString(sum[:])
	if nonceFromClaims != hashedNonce {
		slog.Info("not expected identityToken's nonce")
		return nil, errors.New(message.AppleSignInFailed)
	}

}
