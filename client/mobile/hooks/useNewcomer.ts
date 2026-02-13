import { useMutation } from "@tanstack/react-query";
import { setBirthYear, setNickname } from "@/api/newcomer";
import Toast from "react-native-toast-message";

function useSetNickname() {
  return useMutation({
    mutationFn: setNickname,
    onSuccess: async () => {},
    onError: (error) => {
      Toast.show({
        type: "error",
        text1: error.message,
      });
    },
  });
}

function useSetBirthYear() {
  return useMutation({
    mutationFn: setBirthYear,
    onSuccess: async () => {},
    onError: (error) => {
      Toast.show({
        type: "error",
        text1: error.message,
      });
    },
  });
}

export function useNewcomer() {
  const setNicknameMutation = useSetNickname();
  const setBirthYearMutation = useSetBirthYear();

  return {
    setNicknameMutation,
    setBirthYearMutation,
  };
}
