import { StyleSheet, Text, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { colors } from "@/constants";
import { useJoinOnlineConversation } from "@/hooks/useConversation";

export default function OnlineConversationDetailScreen() {
  const { id } = useLocalSearchParams();
  const {
    data: conversation,
    isPending,
    isError,
  } = useJoinOnlineConversation(String(id));

  if (isPending || isError) {
    return <></>;
  }

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <Text style={styles.title}>{conversation?.title}</Text>
        <Text style={styles.description}>description</Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
  },
  title: {
    fontSize: 18,
    color: colors.BLACK,
    fontWeight: 600,
    marginVertical: 8,
  },
  description: {
    fontSize: 16,
    color: colors.BLACK,
    marginBottom: 14,
  },
});
