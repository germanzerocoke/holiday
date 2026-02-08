import React from "react";
import { StyleSheet, View, type StyleProp, type ViewStyle } from "react-native";
import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

type FormValue = {
  // Kept for form compatibility with parent form
  countryCode: string;
  phoneNumber: string;
};

type PhoneNumberInputProps = {
  containerStyle?: StyleProp<ViewStyle>;
};

export default function PhoneNumberInput(props: PhoneNumberInputProps) {
  const { containerStyle } = props;
  const { control } = useFormContext<FormValue>();

  return (
    <View style={[styles.container, containerStyle]}>
      <Controller
        name="phoneNumber"
        control={control}
        defaultValue=""
        rules={{
          validate: (data: string) => {
            if (data.length <= 0) {
              return "Put your Phone Number";
            }
          },
        }}
        render={({ field: { onChange, onBlur, value }, fieldState }) => (
          <InputField
            value={value}
            onBlur={onBlur}
            onChangeText={(text) => {
              // Keep only digits, but don't format
              const digitsOnly = text.replace(/[^\d]/g, "");
              onChange(digitsOnly);
            }}
            keyboardType="phone-pad"
            placeholder="your phone number"
            variant="standard"
            error={fieldState.error?.message}
            submitBehavior="blurAndSubmit"
            returnKeyType="done"
          />
        )}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    minHeight: 44,
  },
});
