import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function ShortStoryInput() {
  const { control } = useFormContext();

  return (
    <Controller
      name="shortStory"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="email"
          placeholder="The Horse Dealer’s Daughter Cathedral Barn Burning"
          inputMode="text"
          value={value}
          onChangeText={onChange}
        />
      )}
    />
  );
}
