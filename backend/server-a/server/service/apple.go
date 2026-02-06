package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"server-a/server/constant"
	"server-a/server/constant/message"
	"server-a/server/dto"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) SignInWithApple(
	ctx context.Context,
	user,
	nonce,
	identityToken string,
	email *string,
) (
	*dto.SignInWithAppleResponse,
	string, /*refreshToken*/
	error,
) {
	keyReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		constant.AppleKeyUrl,
		nil,
	)
	if err != nil {
		slog.Error("fail to make http request",
			"err", err,
		)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	keyRes, err := client.Do(keyReq)
	if err != nil {
		slog.Error("fail to do request", "err", err)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	defer keyRes.Body.Close()

	var jwks jwt.VerificationKeySet
	err = json.NewDecoder(keyRes.Body).Decode(&jwks)
	if err != nil {
		slog.Error("fail to decode jwks in time",
			"err", err,
			"resp", keyRes,
		)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	idt, err := jwt.Parse(identityToken, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			slog.Info("unexpected signing method")
			return nil, errors.New(message.AppleSignInFailed)
		}
		return jwks.Keys, nil
	})

	issFromClaims, err := idt.Claims.GetIssuer()
	if err != nil {
		slog.Info("fail to get issuer",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if issFromClaims != constant.AppleIssuerUrl {
		slog.Info("not expected apple issuer",
			"iss", issFromClaims,
		)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	audsFromClaims, err := idt.Claims.GetAudience()
	if err != nil {
		slog.Info("fail to get audience",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if len(audsFromClaims) == 0 || audsFromClaims[0] != s.audience {
		slog.Info("no audience or not expected audience")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	exp, err := idt.Claims.GetExpirationTime()
	if err != nil {
		slog.Info("fail to get expiration time")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if exp.Unix() < time.Now().Unix() {
		slog.Info("stale token")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	nonceFromClaims, ok := idt.Claims.(jwt.MapClaims)["nonce"].(string)
	if !ok {
		slog.Info("no nonce in claims")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if nonceFromClaims != nonce {
		slog.Info("not expected identityToken's nonce")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	userFromClaims, err := idt.Claims.GetSubject()
	if err != nil {
		slog.Info("fail to get subject",
			"err", err,
		)
	}

	if userFromClaims != user {
		slog.Info("unexpected user")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if email != nil {
		id := gocql.TimeUUID()
		err = s.repository.SaveAppleSignInInfo(id, user, *email)
		if err != nil {
			return nil, "", errors.New(message.AppleSignInFailed)
		}

		sessionId, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to generate random uuid for session")
			return nil, "", errors.New(message.AppleSignInFailed)
		}

		err = s.repository.SaveEmailBySessionId(sessionId, *email)
		if err != nil {
			return nil, "", errors.New(message.AppleSignInFailed)
		}

		resp := dto.SignInWithAppleResponse{
			PhoneNumberVerified: false,
			SessionId:           sessionId.String(),
		}
		return &resp, "", nil
	}

	id, emailFromDB, role, err := s.repository.FindAppleSignInInfoByUser(user)
	if err != nil {
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	//this additional fetching can be removed to improve speed little bit
	//by adding few lines, but the advantage is also small currently and
	//make link phone number process more complicate
	phoneNumberVerified, err := s.repository.FindPhoneNumberVerifiedById(id)
	if err != nil {
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	if !phoneNumberVerified {
		sessionId, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to generate random uuid for session")
			return nil, "", errors.New(message.AppleSignInFailed)
		}

		err = s.repository.SaveEmailBySessionId(sessionId, emailFromDB)
		if err != nil {
			return nil, "", errors.New(message.AppleSignInFailed)
		}

		resp := dto.SignInWithAppleResponse{
			PhoneNumberVerified: false,
			SessionId:           sessionId.String(),
		}
		return &resp, "", nil
	}

	jti, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to create random uuid for jti")
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	at, rt, err := s.createLoginTokens(id.String(), jti.String(), role)
	if err != nil {
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", errors.New(message.AppleSignInFailed)
	}

	resp := dto.SignInWithAppleResponse{
		PhoneNumberVerified: true,
		AccessToken:         at,
	}
	return &resp, rt, nil
}
