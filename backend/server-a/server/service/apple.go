package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"server-a/server/constant"
	"server-a/server/dto"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		return nil, "", errAppleSignInFailed
	}

	client := &http.Client{Timeout: 10 * time.Second}
	keyRes, err := client.Do(keyReq)
	if err != nil {
		slog.Error("fail to do request", "err", err)
		return nil, "", errAppleSignInFailed
	}

	defer keyRes.Body.Close()

	var jwks jwt.VerificationKeySet
	err = json.NewDecoder(keyRes.Body).Decode(&jwks)
	if err != nil {
		slog.Error("fail to decode jwks in time",
			"err", err,
			"resp", keyRes,
		)
		return nil, "", errAppleSignInFailed
	}

	idt, err := jwt.Parse(identityToken, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			slog.Info("unexpected signing method")
			return nil, errAppleSignInFailed
		}
		return jwks.Keys, nil
	})

	issFromClaims, err := idt.Claims.GetIssuer()
	if err != nil {
		slog.Info("fail to get issuer",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", errAppleSignInFailed
	}

	if issFromClaims != constant.AppleIssuerUrl {
		slog.Info("not expected apple issuer",
			"iss", issFromClaims,
		)
		return nil, "", errAppleSignInFailed
	}

	audsFromClaims, err := idt.Claims.GetAudience()
	if err != nil {
		slog.Info("fail to get audience",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", errAppleSignInFailed
	}

	if len(audsFromClaims) == 0 || audsFromClaims[0] != s.audience {
		slog.Info("no audience or not expected audience")
		return nil, "", errAppleSignInFailed
	}

	exp, err := idt.Claims.GetExpirationTime()
	if err != nil {
		slog.Info("fail to get expiration time")
		return nil, "", errAppleSignInFailed
	}

	if exp.Unix() < time.Now().Unix() {
		slog.Info("stale token")
		return nil, "", errAppleSignInFailed
	}

	nonceFromClaims, ok := idt.Claims.(jwt.MapClaims)["nonce"].(string)
	if !ok {
		slog.Info("no nonce in claims")
		return nil, "", errAppleSignInFailed
	}

	if nonceFromClaims != nonce {
		slog.Info("not expected identityToken's nonce")
		return nil, "", errAppleSignInFailed
	}

	userFromClaims, err := idt.Claims.GetSubject()
	if err != nil {
		slog.Info("fail to get subject",
			"err", err,
		)
	}

	if userFromClaims != user {
		slog.Info("unexpected user")
		return nil, "", errAppleSignInFailed
	}

	if email != nil {
		_, phoneNumberVerified, id, _, role, err := s.repository.FindLoginInfoByEmail(*email)
		if errors.Is(err, gocql.ErrNotFound) {
			err = nil
			idv7, err := uuid.NewV7()
			if err != nil {
				slog.Error("fail to create uuid v7 for apple sign in user")
				return nil, "", errAppleSignInFailed
			}
			id = gocql.UUID(idv7)
			err = s.repository.SaveAppleSignInInfo(id, user, *email, false)
			if err != nil {
				return nil, "", errAppleSignInFailed
			}

			sessionId, err := gocql.RandomUUID()
			if err != nil {
				slog.Error("fail to generate random uuid for session")
				return nil, "", errAppleSignInFailed
			}

			err = s.repository.SaveEmailBySessionId(sessionId, *email)
			if err != nil {
				return nil, "", errAppleSignInFailed
			}

			resp := dto.SignInWithAppleResponse{
				PhoneNumberVerified: false,
				SessionId:           sessionId.String(),
			}
			return &resp, "", nil
		}
		if err != nil {
			return nil, "", errAppleSignInFailed
		}

		err = s.repository.SaveAppleSignInInfo(id, user, *email, phoneNumberVerified)
		if err != nil {
			return nil, "", errAppleSignInFailed
		}

		if !phoneNumberVerified {
			sessionId, err := gocql.RandomUUID()
			if err != nil {
				slog.Error("fail to generate random uuid for session")
				return nil, "", errAppleSignInFailed
			}

			err = s.repository.SaveEmailBySessionId(sessionId, *email)
			if err != nil {
				return nil, "", errAppleSignInFailed
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
			return nil, "", errAppleSignInFailed
		}

		at, rt, err := s.createLoginTokens(id.String(), jti.String(), role)
		if err != nil {
			return nil, "", errAppleSignInFailed
		}

		err = s.repository.SaveRefreshTokenJTIById(id, jti)
		if err != nil {
			return nil, "", errAppleSignInFailed
		}

		resp := dto.SignInWithAppleResponse{
			PhoneNumberVerified: true,
			AccessToken:         at,
		}
		return &resp, rt, nil

	}

	id, emailFromDB, role, err := s.repository.FindAppleSignInInfoByUser(user)
	if err != nil {
		return nil, "", errAppleSignInFailed
	}

	//this additional fetching can be removed to improve speed little bit
	//by adding few lines, but the advantage is also small currently and
	//make link phone number process more complicate
	phoneNumberVerified, err := s.repository.FindPhoneNumberVerifiedById(id)
	if err != nil {
		return nil, "", errAppleSignInFailed
	}

	if !phoneNumberVerified {
		sessionId, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to generate random uuid for session")
			return nil, "", errAppleSignInFailed
		}

		err = s.repository.SaveEmailBySessionId(sessionId, emailFromDB)
		if err != nil {
			return nil, "", errAppleSignInFailed
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
		return nil, "", errAppleSignInFailed
	}

	at, rt, err := s.createLoginTokens(id.String(), jti.String(), role)
	if err != nil {
		return nil, "", errAppleSignInFailed
	}

	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", errAppleSignInFailed
	}

	resp := dto.SignInWithAppleResponse{
		PhoneNumberVerified: true,
		AccessToken:         at,
	}
	return &resp, rt, nil
}
