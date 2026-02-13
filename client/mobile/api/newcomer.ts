import { axiosInstance } from "@/api/axios";

export async function setNickname(body: { nickname: string }) {
  const { data } = await axiosInstance.post("/profile/newcomer/nickname", body);
  return data;
}

export async function setBirthYear(body: { birthYear: number }) {
  const { data } = await axiosInstance.post(
    "/profile/newcomer/birth-year",
    body,
  );
  return data;
}
