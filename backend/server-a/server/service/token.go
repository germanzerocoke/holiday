package service

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"server-a/server/constant"
	"server-a/server/dto"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) GenerateAccessToken(refreshToken string) (*dto.TokenRefreshResponse, error) {
	rt, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return s.secretKeyRT, nil
	})
	if err != nil {
		log.Printf("fail to parse token: %v", err)
		return nil, err
	}
	if !rt.Valid {
		log.Printf("invalid token: %v", rt)
		return nil, errors.New("invalid token")
	}
	id, err := rt.Claims.GetSubject()
	if err != nil {
		log.Printf("fail to get subject from claim: %v", err)
		return nil, err
	}
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		log.Printf("fail to parse gocql uuid from id: %v", err)
	}
	rtInDB, err := s.repository.FindRefreshTokenById(uuid)
	if err != nil {
		log.Printf("fail to find token: %v", err)
		return nil, err
	}
	if rt.Raw != rtInDB {
		log.Printf(
			"token is not same with DB- received token: %v, db token: %v",
			rt.Raw, rtInDB,
		)
		return nil, fmt.Errorf("invalid token: %v", rt.Raw)
	}
	role, ok := rt.Claims.(jwt.MapClaims)["role"].(string)
	if !ok {
		log.Printf("fail to get role - token: %v", rt.Raw)
		return nil, fmt.Errorf("fail to get role - token: %v", rt.Raw)
	}
	at, err := createToken(id, role, s.secretKeyAT, constant.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	resp := dto.TokenRefreshResponse{AccessToken: at}

	return &resp, nil
}

func (s *Service) createLoginTokens(id, role string) (accessToken, refreshToken string, err error) {
	at, err := createToken(id, role, s.secretKeyAT, constant.AccessTokenTTL)
	if err != nil {
		slog.Error("fail to create access token",
			"err", err,
			"id", id,
		)
		return "", "", err
	}
	rt, err := createToken(id, role, s.secretKeyRT, constant.RefreshTokenTTL)
	if err != nil {
		slog.Error("fail to create refresh token",
			"err", err,
			"id", id,
		)
		return "", "", err
	}
	return at, rt, nil
}

func createToken(id, role string, secretKey []byte, ttl int64) (token string, err error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		log.Printf("fail to make token")
		return "", err
	}
	return token, nil
}
