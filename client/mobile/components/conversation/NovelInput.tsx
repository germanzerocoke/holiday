import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function NovelInput() {
  const { control } = useFormContext();

  return (
    <Controller
      name="novel"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Novel"
          placeholder="The Castle The Stranger Norwegian Wood"
          inputMode="text"
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={value}
          onChangeText={onChange}
        />
      )}
    />
  );
}
