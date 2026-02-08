import { axiosInstance } from "@/api/axios";
import {
  emailLoginResponse,
  emailSignRequest,
  sendSMSOTPRequest,
  signInWithAppleRequest,
  signInWithAppleResponse,
  verifyEmailOTPRequest,
  verifyEmailOTPResponse,
  verifySMSOTPRequest,
  verifySMSOTPResponse,
} from "@/types/auth";

export async function postEmailSignup(body: emailSignRequest): Promise<string> {
  console.log("post email sign up");
  const { data } = await axiosInstance.post("/auth/email/create", body);
  return data;
}

export async function postEmailLogin(
  body: emailSignRequest,
): Promise<emailLoginResponse> {
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

export async function requestSMSOTP(
  body: sendSMSOTPRequest,
): Promise<{ verificationId: string }> {
  const { data } = await axiosInstance.post("/auth/sms/otp/send", body);
  return data;
}

export async function verifyEmailOTP(
  body: verifyEmailOTPRequest,
): Promise<verifyEmailOTPResponse> {
  const { data } = await axiosInstance.post("/auth/email/otp/verify", body);
  return data;
}

export async function verifySMSOTP(
  body: verifySMSOTPRequest,
): Promise<verifySMSOTPResponse> {
  const { data } = await axiosInstance.post("/auth/sms/otp/verify", body);
  return data;
}

export async function signInWithApple(
  body: signInWithAppleRequest,
): Promise<signInWithAppleResponse> {
  const { data } = await axiosInstance.post("/auth/email/apple", body);
  return data;
}

export async function getMe() {
  const { data } = await axiosInstance.get("/auth/me");
  return data;
}
