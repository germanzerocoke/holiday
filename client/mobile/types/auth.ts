export interface Token {
  accessToken: string;
  refreshToken: string;
}

export interface EmailSignRequest {
  email: string;
  password: string;
}

export interface VerifyEmailOTPRequest {
  otp: string;
  verificationId: string;
}

export interface VerifyEmailOTPResponse {
  emailVerified: boolean;
  sessionId?: string;
}

export interface VerifySMSOTPRequest {
  otp: string;
  verificationId: string;
  sessionId?: string;
}

export interface VerifySMSOTPResponse {
  phoneNumberVerified: boolean;
  accessToken?: string;
}

export interface EmailLoginResponse {
  emailVerified: boolean;
  phoneNumberVerified: boolean;
  id: string;
  sessionId?: string;
  accessToken?: string;
}

export interface SignInWithAppleRequest {
  user: string;
  email: string | null;
  identityToken: string | null;
  nonce: string,
}

export interface SignInWithAppleResponse {
  phoneNumberVerified: boolean;
  sessionId?: string;
  accessToken?: string;
}
