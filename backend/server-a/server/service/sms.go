package service

import (
	"fmt"
	"log/slog"
	"os"
	"server-a/server/constant"
	"server-a/server/dto"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	_ "github.com/joho/godotenv/autoload"
)

func (s *Service) SendSMSOTP(phoneNumber string) (*dto.OTPSendResponse, error) {
	serviceSid := os.Getenv("TWILIO_SERVICE_SID")

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := s.twilioClient.VerifyV2.CreateVerification(serviceSid, params)
	if err != nil {
		slog.Info("fail to send sms otp code", err)
		return nil, err
	}
	if resp.Status != nil {
		slog.Info("success to send sms otp code", "phoneNumber", resp.To, "status", *resp.Status)
	}
	vid, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for verification id")
		return nil, err
	}
	err = s.repository.SavePhoneNumberByVerificationId(vid, *resp.To)
	if err != nil {
		return nil, err
	}
	res := dto.OTPSendResponse{VerificationId: vid.String()}
	return &res, nil
}

func (s *Service) VerifySMSOTP(sessionId *string, otp, verificationId string) (*dto.SMSOTPVerifyResponse, string /*refreshToken*/, error) {
	var email string
	if sessionId != nil {
		sid, err := gocql.ParseUUID(*sessionId)
		if err != nil {
			slog.Info("fail to parse SessionId from OTPVerifyReq",
				"err", err,
				"SessionId", *sessionId,
			)
			return nil, "", err
		}

		email, err = s.repository.FindEmailBySessionId(sid)
		if err != nil {
			return nil, "", err
		}
	}

	vid, err := gocql.ParseUUID(verificationId)
	if err != nil {
		slog.Info("fail to parse verificationId from request",
			"err", err,
			"id", verificationId,
		)
	}
	phoneNumber, err := s.repository.FindPhoneNumberByVerificationId(vid)
	if err != nil {
		return nil, "", err
	}
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(otp)

	serviceSid := os.Getenv("TWILIO_SERVICE_SID")

	resp, err := s.twilioClient.VerifyV2.CreateVerificationCheck(serviceSid, params)
	if err != nil {
		slog.Error("fail to verify phone number otp", "err", err)
		return nil, "", err
	}
	if resp.Status == nil {
		slog.Error("status is nil pointer")
		return nil, "", fmt.Errorf("twilio response status is nil pointer: %v", verificationId)
	}
	if *resp.Status != "approved" {
		slog.Info("otp is not correct",
			"otp", otp,
		)
		r := dto.SMSOTPVerifyResponse{
			PhoneNumberVerified: false,
		}
		return &r, "", nil
	}

	if sessionId == nil {
		id := gocql.TimeUUID()
		err = s.repository.SavePhoneNumberLoginInfo(phoneNumber, id)
		if err != nil {
			return nil, "", err
		}
		jti, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to make random uuid for jti")
			return nil, "", err
		}
		at, rt, err := s.createLoginTokens(id.String(), jti.String(), constant.RoleUser)
		if err != nil {
			return nil, "", err
		}
		err = s.repository.SaveRefreshTokenJTIById(id, jti)
		if err != nil {
			return nil, "", err
		}
		r := dto.SMSOTPVerifyResponse{
			PhoneNumberVerified: true,
			AccessToken:         at,
		}
		return &r, rt, nil
	}

	id, role, createdTime, err := s.repository.FindMemberInfoByEmail(email)
	if err != nil {
		return nil, "", err
	}
	err = s.repository.LinkPhoneNumberToMember(id, email, phoneNumber, role, createdTime)
	if err != nil {
		return nil, "", err
	}

	jti, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for jti")
		return nil, "", err
	}
	at, rt, err := s.createLoginTokens(id.String(), jti.String(), constant.RoleUser)
	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", err
	}
	r := dto.SMSOTPVerifyResponse{
		PhoneNumberVerified: true,
		AccessToken:         at,
	}
	return &r, rt, nil
}
