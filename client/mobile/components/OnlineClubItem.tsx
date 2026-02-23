import { Pressable, StyleSheet, Text, View } from "react-native";
import { colors } from "@/constants";
import { OnlineClub } from "@/types/club";
import { router } from "expo-router";

interface OnlineClubItemProps {
  club: OnlineClub;
}

export default function OnlineClubItem({ club }: OnlineClubItemProps) {
  return (
    <Pressable
      style={styles.container}
      onPress={() => router.push(`/club/online/${club.id}`)}
    >
      <View style={styles.content}>
        <Text style={styles.title}>{club.title}</Text>
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
