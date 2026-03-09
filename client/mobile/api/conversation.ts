import { axiosInstance } from "@/api/axios";
import {
  OnlineConversation,
  OnlineConversationDetail,
} from "@/types/conversation";

export async function getOnlineConversations(
  page = 0,
): Promise<OnlineConversation[]> {
  const { data } = await axiosInstance.get(
    `/conversation/online/list?page=${page}`,
  );
  return data;
}

export async function joinOnlineConversation(
  id: string,
): Promise<OnlineConversationDetail> {
  const { data } = await axiosInstance.get(`/conversation/online/${id}`);
  return data;
}
