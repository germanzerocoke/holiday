import { StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import { FormProvider, useForm } from "react-hook-form";
import hairlineWidth = StyleSheet.hairlineWidth;
import NicknameInput from "@/components/NicknameInput";

interface FormValue {
  nickname: string;
}

export default function NicknameScreen() {
  const nicknameForm = useForm<FormValue>({
    defaultValues: {
      nickname: "",
    },
  });

  const onSubmit = async (formValue: FormValue) => {
    const { nickname } = formValue;
  };

  return (
    <FormProvider {...nicknameForm}>
      <SafeAreaView style={styles.container}>
        <Text style={styles.guideLine}>
          We strongly recommend to set a nickname which is convenient to be
          called by others
        </Text>
        <View style={styles.content}>
          <NicknameInput />
        </View>
        <FixedBottomCTA
          label="Set nickname"
          onPress={nicknameForm.handleSubmit(onSubmit)}
        />
      </SafeAreaView>
    </FormProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
    borderTopWidth: hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  content: {
    flex: 1,
    width: "100%",
    paddingHorizontal: 100,
    paddingTop: 100,
    gap: 50,
  },
  guideLine: {
    position: "absolute",
    top: 40,
    paddingHorizontal: 30,
    fontSize: 15,
  },
});
