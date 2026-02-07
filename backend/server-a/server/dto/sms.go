package dto

type SMSOTPSendRequest struct {
	SessionId   *string `json:"sessionId"` //nullable
	PhoneNumber string  `json:"phoneNumber"`
}
type SMSOTPVerifyRequest struct {
	SessionId      *string `json:"sessionId"` //nullable
	VerificationId string  `json:"verificationId"`
	OTP            string  `json:"otp"`
}

type SMSOTPVerifyResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	AccessToken         string `json:"accessToken"`
}
