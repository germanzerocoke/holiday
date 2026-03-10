import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function ByInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="by"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="By"
          placeholder="William Shakespeare Nastume Soseki Ezra Pound"
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
