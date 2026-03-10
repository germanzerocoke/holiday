import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function PoemInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="poem"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="poem"
          placeholder="I'm Nobody! Who are you? The Love Song of J. Alfred Prufrock Sonnet 18"
          inputMode="text"
          returnKeyType="next"
          submitBehavior="submit"
          value={value}
          onChangeText={onChange}
        />
      )}
    />
  );
}
