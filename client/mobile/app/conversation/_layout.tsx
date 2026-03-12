import { Stack } from "expo-router";

export default function ConversationLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="online"
        options={{
          title: "",
          headerShown: false,
        }}
      />
    </Stack>
  );
}
