import { axiosInstance } from "@/api/axios";
import {
  CreateOnlineConversationRequest,
  OnlineConversation,
  OnlineConversationDetail,
} from "@/types/conversation";
import { data } from "browserslist";

export async function getOnlineConversations(
  page = 0,
): Promise<OnlineConversation[]> {
  const { data } = await axiosInstance.get(
    `/online/conversation/list?page=${page}`,
  );
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
