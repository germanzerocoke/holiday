import { Stack } from "expo-router";

export default function ClubLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "",
          headerShown: false,
        }}
      />
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
