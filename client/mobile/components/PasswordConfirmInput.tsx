import { Controller, useFormContext, useWatch } from "react-hook-form";
import InputField from "@/components/InputField";

export default function PasswordConfirmInput() {
  const { control } = useFormContext();
  const password = useWatch({ control, name: "password" });

  return (
    <Controller
      name="passwordConfirm" //use name in signup formValue
      control={control}
      rules={{
        validate: (data: string) => {
          if (data !== password) {
            return "password is not same";
          }
        },
      }}
      render={({ field: { ref, onChange, value }, fieldState: { error } }) => (
        <InputField
          ref={ref}
          variant="standard"
          label="password confirm"
          placeholder="put your password again"
          submitBehavior="blurAndSubmit"
          secureTextEntry
          textContentType="oneTimeCode"
          value={value}
          onChangeText={onChange}
          error={error?.message}
        />
      )}
    />
  );
}
