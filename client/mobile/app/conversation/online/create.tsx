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
import YearInput from "@/components/YearInput";
import MonthDayInput from "@/components/MonthDayInput";
import HourInput from "@/components/HourInput";
import MinuteInput from "@/components/MinuteInput";
import LengthInput from "@/components/LengthInput";

interface FormValue {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by?: string;
  rule?: string;
  capacity: number;
  year: string;
  monthDay: string;
  hour: string;
  minute: string;
  length: string;
}

export default function OnlineConversationCreateScreen() {
  const now = new Date();
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
      year: String(now.getFullYear()),
      monthDay: `${now.getMonth() + 1}.${now.getDate()}`,
      hour: String(now.getHours()),
      minute: String(now.getMinutes()),
      length: "100",
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
      monthDay,
      hour,
      minute,
      length,
    } = formValues;
    const monthDayParts = monthDay.split(".");
    const month = monthDayParts[0] ?? String(now.getMonth() + 1);
    const day = monthDayParts[1] ?? String(now.getDate());
    const when = new Date(
      Number(year),
      Number(month) - 1,
      Number(day),
      Number(hour),
      Number(minute),
    ).toISOString();
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
        length: `${length}m0s`,
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
          <View style={styles.dateTimeRow}>
            <View style={styles.yearBox}>
              <YearInput />
            </View>
            <View style={styles.monthDayBox}>
              <MonthDayInput />
            </View>
            <View style={styles.timeBox}>
              <HourInput />
            </View>
            <View style={styles.timeBox}>
              <MinuteInput />
            </View>
          </View>
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
  dateTimeRow: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: 8,
  },
  yearBox: {
    flex: 1.2,
  },
  monthDayBox: {
    flex: 2,
  },
  timeBox: {
    flex: 1,
  },
});
