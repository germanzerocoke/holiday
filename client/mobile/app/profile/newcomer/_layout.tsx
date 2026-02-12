import { Stack } from "expo-router";

export default function NewcomerLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="nickname"
        options={{
          title: "Set your nickname",
          headerShown: true,
        }}
      />
    </Stack>
  );
}
