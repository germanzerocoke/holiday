import { Stack } from "expo-router";

export default function OnlineLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="list"
        options={{
          title: "Online club list",
          headerShown: true,
        }}
      />
      <Stack.Screen
        name="create"
        options={{
          title: "Create your online club",
          headerShown: true,
        }}
      />
    </Stack>
  );
}
