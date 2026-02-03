package dto

type EmailRequest struct {
	Email string `json:"email"`
}
type EmailMemberSaveRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailMemberLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailLoginResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	EmailVerified       bool   `json:"emailVerified"`
	Id                  string `json:"id"`
	SessionId           string `json:"sessionId"`
	AccessToken         string `json:"accessToken"`
}

type EmailOTPSendRequest struct {
	Id string `json:"id"`
}

type EmailOTPVerifyRequest struct {
	VerificationId string `json:"verificationId"`
	OTP            string `json:"otp"`
}
type EmailOTPVerifyResponse struct {
	EmailVerified bool   `json:"emailVerified"`
	SessionId     string `json:"session_id"`
}

type SignInWithAppleRequest struct {
	User          string  `json:"user"`
	Email         *string `json:"email"` //nullable
	IdentityToken string  `json:"identityToken"`
	Nonce         string  `json:"nonce"`
}

type SignInWithAppleResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	SessionId           string `json:"sessionId"`
	AccessToken         string `json:"accessToken"`
}
