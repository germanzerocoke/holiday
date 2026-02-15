export interface SignInWithEmailRequest {
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

export interface SendSMSOTPRequest {
  phoneNumber: string;
  sessionId: string | null;
}

export interface VerifySMSOTPRequest {
  otp: string;
  verificationId: string;
  sessionId: string | null;
}

export interface VerifySMSOTPResponse {
  phoneNumberVerified: boolean;
  accessToken?: string;
}

export interface LoginWithEmailResponse {
  emailVerified: boolean;
  phoneNumberVerified: boolean;
  verificationId?: string;
  sessionId?: string;
  accessToken?: string;
}

export interface SignInWithAppleRequest {
  user: string;
  email: string | null;
  identityToken: string;
  nonce: string;
}

export interface SignInWithAppleResponse {
  phoneNumberVerified: boolean;
  sessionId?: string;
  accessToken?: string;
}
