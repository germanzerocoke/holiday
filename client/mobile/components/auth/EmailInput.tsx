import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function EmailInput() {
  const { control, setFocus } = useFormContext();

  return (
    <Controller
      name="email"
      control={control}
      rules={{
        validate: (data: string) => {
          if (data.length === 0) {
            return "put your email";
          }
          if (!/^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/.test(data)) {
            return "not right email form";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => (
        <InputField
          autoFocus
          variant="standard"
          label="email"
          placeholder="please put your email"
          inputMode="email" //keyboard form numeric can also
          returnKeyType="next"
          submitBehavior="submit" //keyboard not going down
          onSubmitEditing={() => setFocus("password")}
          value={value}
          onChangeText={onChange}
          error={error?.message}
        />
      )}
    />
  );
}
