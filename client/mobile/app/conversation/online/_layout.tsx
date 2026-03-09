import { Stack } from "expo-router";

export default function OnlineLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="list"
        options={{
          title: "Online conversation list",
          headerShown: true,
        }}
      />
      <Stack.Screen
        name="create"
        options={{
          title: "Create your online conversation",
          headerShown: true,
        }}
      />
    </Stack>
  );
}
