import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function RuleInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="rule"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Rule"
          placeholder="1. Respect each other 2. Try to speak more 3. Try to yield and listen"
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
