import { useInfiniteQuery, useMutation, useQuery } from "@tanstack/react-query";
import {
  getOnlineConversations,
  joinOnlineConversation,
  createOnlineConversation,
} from "@/api/conversation";
import { queryKey } from "@/constants";

export function useGetInfiniteOnlineConversations() {
  return useInfiniteQuery({
    queryFn: ({ pageParam }) => getOnlineConversations(pageParam),
    queryKey: [queryKey.CONVERSATION, queryKey.GET_ONLINE_CONVERSATIONS],
    initialPageParam: 0,
    getNextPageParam: (lastPage, allPages) => {
      const lastPost = lastPage[lastPage.length - 1];
      return lastPost ? allPages.length + 1 : undefined;
    },
  });
}

export function useJoinOnlineConversation(id: string) {
  return useQuery({
    queryFn: () => joinOnlineConversation(id),
    queryKey: [queryKey.CONVERSATION, queryKey.JOIN_ONLINE_CONVERSATION, id],
    enabled: Boolean(id),
  });
}

export function useCreateOnlineConversation() {
  return useMutation({
    mutationFn: createOnlineConversation,
  });
}
