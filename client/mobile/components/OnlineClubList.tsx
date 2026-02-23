import { useRef, useState } from "react";
import { FlatList, StyleSheet } from "react-native";
import { colors } from "@/constants";
import { useScrollToTop } from "@react-navigation/native";
import OnlineClubItem from "@/components/OnlineClubItem";
import { useGetInfiniteOnlineClubs } from "@/hooks/useClub";

export default function OnlineClubList() {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, refetch } =
    useGetInfiniteOnlineClubs();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const ref = useRef<FlatList | null>(null);
  useScrollToTop(ref);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await refetch();
    setIsRefreshing(false);
  };

  const handleEndReached = () => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  return (
    <FlatList
      ref={ref}
      data={data?.pages.flat()}
      renderItem={({ item }) => <OnlineClubItem club={item} />}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.contentContainer}
      onEndReached={handleEndReached}
      onEndReachedThreshold={0.5}
      refreshing={isRefreshing}
      onRefresh={handleRefresh}
    />
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
});
