import { axiosInstance } from "@/api/axios";
import { OnlineClub, OnlineClubDetail } from "@/types/club";

export async function getOnlineClubs(page = 1): Promise<OnlineClub[]> {
  const { data } = await axiosInstance.get(`/club/online/list?page=${page}`);
  return data;
}

export async function getOnlineClub(id: string): Promise<OnlineClubDetail> {
  const { data } = await axiosInstance.get(`/club/online/${id}`);
  return data;
}
