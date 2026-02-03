package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"server-a/server/constant"
	"server-a/server/constant/message"
	"server-a/server/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) SignInWithApple(
	ctx context.Context,
	user,
	rawNonce,
	identityToken string,
	email *string,
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
	if !idt.Valid {
		slog.Info("token is not valid")
		return nil, errors.New(message.AppleSignInFailed)
	}

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

	//this check is necessary for the window for which exp does not show up at our jwt claims
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

	userFromClaims, err := idt.Claims.GetSubject()
	if err != nil {
		slog.Info("fail to get subject",
			"err", err,
		)
	}
	if userFromClaims != user {
		slog.Info("unexpected user")
		return nil, errors.New(message.AppleSignInFailed)
	}

	//sign in success

	if email != nil {
		//first time sign up user
	}

	//login user
	emailVerified, phoneNumberVerified, id, role, err := s.repository.FindLoginInfoByAppleSignInUser(user)

}
