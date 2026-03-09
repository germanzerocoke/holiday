import { Pressable, StyleSheet } from "react-native";
import { router } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { Feather } from "@expo/vector-icons";
import { colors } from "@/constants";
import OnlineConversationList from "@/components/OnlineConversationList";

export default function OnlineConversationListScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <OnlineConversationList />
      <Pressable
        style={styles.createButton}
        onPress={() => router.push("/conversation/online/create")}
      >
        <Feather name="plus" size={32} color="black" />
      </Pressable>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
  },
  createButton: {
    position: "absolute",
    bottom: 16,
    right: 16,
    backgroundColor: colors.WHITE,
    width: 64,
    height: 64,
    borderRadius: 32,
    justifyContent: "center",
    alignItems: "center",
    shadowColor: colors.BLACK,
    shadowOffset: { width: 0, height: 2 },
    shadowRadius: 3,
    shadowOpacity: 0.5,
    elevation: 2, //for android shadowing
  },
});
