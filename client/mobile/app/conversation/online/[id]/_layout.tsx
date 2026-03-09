import { Stack } from "expo-router";

export default function DetailLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "",
          headerShown: true,
        }}
      />
      {/*<Stack.Screen*/}
      {/*  name="edit"*/}
      {/*  options={{*/}
      {/*    title: "Edit your online club",*/}
      {/*    headerShown: true,*/}
      {/*  }}*/}
      {/*/>*/}
    </Stack>
  );
}
