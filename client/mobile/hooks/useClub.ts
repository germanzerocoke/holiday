import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { getOnlineClub, getOnlineClubs } from "@/api/club";
import { queryKey } from "@/constants";

export function useGetInfiniteOnlineClubs() {
  return useInfiniteQuery({
    queryFn: ({ pageParam }) => getOnlineClubs(pageParam),
    queryKey: [queryKey.CLUB, queryKey.GET_ONLINE_CLUBS],
    initialPageParam: 0,
    getNextPageParam: (lastPage, allPages) => {
      const lastPost = lastPage[lastPage.length - 1];
      return lastPost ? allPages.length + 1 : undefined;
    },
  });
}

export function useGetOnlineClub(id: string) {
  return useQuery({
    queryFn: () => getOnlineClub(id),
    queryKey: [queryKey.CLUB, queryKey.GET_ONLINE_CLUB, id],
    enabled: Boolean(id),
  });
}
