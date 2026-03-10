import { StyleSheet, View } from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import { useCreateOnlineConversation } from "@/hooks/useConversation";
import { router } from "expo-router";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import NovelInput from "@/components/NovelInput";
import ShortStoryInput from "@/components/ShortStoryInput";
import PoemInput from "@/components/PoemInput";
import PlayInput from "@/components/PlayInput";
import FilmInput from "@/components/FilmInput";
import ByInput from "@/components/ByInput";
import RuleInput from "@/components/RuleInput";
import CapacityInput from "@/components/CapacityInput";

interface FormValue {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by?: string;
  rule?: string;
  capacity: number;
  year: number;
  month: number;
  day: number;
  hour: number;
  minute: number;
  length: number;
}

export default function OnlineConversationCreateScreen() {
  const createOnlineConversationMutation = useCreateOnlineConversation();
  const onlineConversationForm = useForm<FormValue>({
    defaultValues: {
      novel: "",
      shortStory: "",
      poem: "",
      play: "",
      film: "",
      by: "",
      rule: "",
      capacity: 6,
      year: new Date().getFullYear(),
      month: new Date().getMonth(),
      day: new Date().getDay(),
      hour: new Date().getHours(),
      minute: new Date().getMinutes(),
      length: 100,
    },
  });
  const onSubmit = (formValues: FormValue) => {
    const {
      novel,
      shortStory,
      poem,
      play,
      film,
      by,
      rule,
      capacity,
      year,
      month,
      day,
      hour,
      minute,
      length,
    } = formValues;
    const when = new Date(year, month, day, hour, minute).toISOString();
    createOnlineConversationMutation.mutate(
      {
        novel: novel,
        shortStory: shortStory,
        poem: poem,
        play: play,
        film: film,
        by: by,
        rule: rule,
        capacity: capacity,
        when: when,
        length: `${String(length)}m0s`,
      },
      {
        onSuccess: () => {
          router.replace("/conversation/online");
        },
      },
    );
  };

  return (
    <FormProvider {...onlineConversationForm}>
      <View style={styles.container}>
        <View style={styles.content}>
          <NovelInput />
          <ShortStoryInput />
          <PoemInput />
          <PlayInput />
          <FilmInput />
          <ByInput />
          <RuleInput />
          <CapacityInput />
          <YearInput />
          <MonthInput />
          <DayInput />
          <HourInput />
          <MinuteInput />
          <LengthInput />
        </View>
        <FixedBottomCTA
          label="Create"
          onPress={onlineConversationForm.handleSubmit(onSubmit)}
        />
      </View>
    </FormProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  content: {
    flex: 1,
    margin: 16,
    gap: 16,
    backgroundColor: colors.SAND_110,
  },
});
