package service

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"server-a/server/constant"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) GenerateAccessToken(refreshToken string) (map[string]string, error) {
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
	exp, err := rt.Claims.GetExpirationTime()
	if err != nil {
		slog.Info("fail to get expiration time")
		return nil, errors.New("invalid token")
	}
	if exp.Unix() < time.Now().Unix() {
		slog.Info("stale token")
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
	jti, err := s.repository.FindRefreshTokenJTIById(uuid)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if rt.Claims.(jwt.MapClaims)["jti"].(string) != jti.String() {
		slog.Info("refresh token jti is not same with DB")
		return nil, errors.New("invalid token")
	}
	role, ok := rt.Claims.(jwt.MapClaims)["role"].(string)
	if !ok {
		slog.Info("fail to get role")
		return nil, errors.New("invalid token")
	}
	at, err := createToken(id, role, s.secretKeyAT, constant.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	resp := map[string]string{"accessToken": at}

	return resp, nil
}

func (s *Service) createLoginTokens(id, jti, role string) (accessToken, refreshToken string, err error) {
	at, err := createToken(id, role, s.secretKeyAT, constant.AccessTokenTTL)
	if err != nil {
		slog.Error("fail to create access token",
			"err", err,
			"id", id,
		)
		return "", "", err
	}
	rt, err := createTokenWithJTI(id, jti, role, s.secretKeyRT, constant.RefreshTokenTTL)
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

func createTokenWithJTI(id, jti, role string, secretKey []byte, ttl int64) (token string, err error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"jti":  jti,
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
