import { axiosInstance } from "@/api/axios";
import {
  EmailLoginResponse,
  EmailSignRequest,
  SignInWithAppleRequest,
  SignInWithAppleResponse,
  VerifyEmailOTPRequest,
  VerifyEmailOTPResponse,
  VerifySMSOTPRequest,
  VerifySMSOTPResponse,
} from "@/types/auth";

export async function postEmailSignup(body: EmailSignRequest): Promise<string> {
  console.log("post email sign up");
  const { data } = await axiosInstance.post("/auth/email/create", body);
  return data;
}

export async function postEmailLogin(
  body: EmailSignRequest,
): Promise<EmailLoginResponse> {
  const { data } = await axiosInstance.post("/auth/email/login", body);
  console.log(data);
  return data;
}

export async function requestEmailOTP(
  id: string,
): Promise<{ verificationId: string }> {
  const { data } = await axiosInstance.post("/auth/email/otp/send", { id });
  return data;
}

export async function requestSMSOTP(phoneNumber: string) {
  const { data } = await axiosInstance.post("/auth/sms/otp/send", {
    phoneNumber,
  });
  return data;
}

export async function verifyEmailOTP(
  body: VerifyEmailOTPRequest,
): Promise<VerifyEmailOTPResponse> {
  const { data } = await axiosInstance.post("/auth/email/otp/verify", body);
  return data;
}

export async function verifySMSOTP(
  body: VerifySMSOTPRequest,
): Promise<VerifySMSOTPResponse> {
  const { data } = await axiosInstance.post("/auth/sms/otp/verify", body);
  return data;
}

export async function signInWithApple(
  body: SignInWithAppleRequest,
): Promise<SignInWithAppleResponse> {
  const { data } = await axiosInstance.post("/auth/email/apple", body);
  return data;
}

export async function getMe() {
  const { data } = await axiosInstance.get("/auth/me");
  return data;
}
