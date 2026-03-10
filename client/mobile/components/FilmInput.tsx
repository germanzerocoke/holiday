import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function FilmInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="film"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Film"
          placeholder="It Happened One Night Cure The Godfather"
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
