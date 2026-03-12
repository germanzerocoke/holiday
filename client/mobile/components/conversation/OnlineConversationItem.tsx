import { Pressable, StyleSheet, Text, View } from "react-native";
import { colors } from "@/constants";
import { router } from "expo-router";
import { ConversationFeedResponse } from "@/types/conversation";

interface OnlineConversationItemProps {
  conversation: ConversationFeedResponse;
}

export default function OnlineConversationItem({
  conversation,
}: OnlineConversationItemProps) {
  return (
    <Pressable
      style={styles.container}
      onPress={() => router.push(`/conversation/online/${conversation.id}`)}
    >
      <View style={styles.content}>
        <Text style={styles.title}>{conversation.when}</Text>
        <Text style={styles.description}>description</Text>
      </View>
    </Pressable>
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
