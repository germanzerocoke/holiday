import { Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import AuthRoute from "@/components/AuthRoute";

export default function InitScreen() {
  return (
    <AuthRoute>
      <SafeAreaView>
        <Text>loading</Text>
      </SafeAreaView>
    </AuthRoute>
  );
}
//
// const styles = StyleSheet.create({
// });
