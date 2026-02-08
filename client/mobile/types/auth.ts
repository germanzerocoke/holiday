export interface Token {
  accessToken: string;
  refreshToken: string;
}

export interface signInWithEmailRequest {
  email: string;
  password: string;
}

export interface verifyEmailOTPRequest {
  otp: string;
  verificationId: string;
}

export interface verifyEmailOTPResponse {
  emailVerified: boolean;
  sessionId?: string;
}

export interface sendSMSOTPRequest {
  phoneNumber: string;
  sessionId: string | null;
}

export interface verifySMSOTPRequest {
  otp: string;
  verificationId: string;
  sessionId: string | null;
}

export interface verifySMSOTPResponse {
  phoneNumberVerified: boolean;
  accessToken?: string;
}

export interface loginWithEmailResponse {
  emailVerified: boolean;
  phoneNumberVerified: boolean;
  id: string;
  sessionId?: string;
  accessToken?: string;
}

export interface signInWithAppleRequest {
  user: string;
  email: string | null;
  identityToken: string;
  nonce: string;
}

export interface signInWithAppleResponse {
  phoneNumberVerified: boolean;
  sessionId?: string;
  accessToken?: string;
}
