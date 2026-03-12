import { axiosInstance } from "@/api/axios";
import {
  ConversationFeedResponse,
  CreateOnlineConversationRequest,
  OnlineConversationDetail,
} from "@/types/conversation";

export async function getOnlineConversations(
  page = 1,
): Promise<ConversationFeedResponse[]> {
  const { data } = await axiosInstance.get(
    `/online/conversation/list?page=${page}`,
  );
  console.log(data);
  return data;
}

export async function joinOnlineConversation(
  id: string,
): Promise<OnlineConversationDetail> {
  const { data } = await axiosInstance.get(`/online/conversation?id=${id}`);
  return data;
}

export async function createOnlineConversation(
  body: CreateOnlineConversationRequest,
): Promise<{ id: string }> {
  const { data } = await axiosInstance.post(
    "/online/conversation/create",
    body,
  );
  return data;
}
