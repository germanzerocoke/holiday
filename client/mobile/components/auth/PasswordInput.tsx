import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";
import { TextInputProps } from "react-native";

interface Props {
  submitBehavior?: TextInputProps["submitBehavior"];
}

export default function PasswordInput({
  submitBehavior = "blurAndSubmit",
}: Props) {
  const { control, setFocus } = useFormContext();

  return (
    <Controller
      name="password" //use name in signup formValue
      control={control}
      rules={{
        validate: (data: string) => {
          if (data.length === 0) {
            return "put your password";
          }
          if (data.length < 8) {
            return "password should be at least more than 8 characters";
          }
        },
      }}
      render={({ field: { ref, onChange, value }, fieldState: { error } }) => (
        <InputField
          ref={ref}
          label="password"
          variant="standard"
          placeholder="please put your password"
          textContentType="oneTimeCode"
          secureTextEntry
          submitBehavior={submitBehavior}
          value={value}
          onChangeText={onChange}
          error={error?.message}
          onSubmitEditing={() => setFocus("passwordConfirm")}
        />
      )}
    />
  );
}
