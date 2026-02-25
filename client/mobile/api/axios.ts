import axios from "axios";
import { Platform } from "react-native";

const baseUrl = {
  android: process.env.EXPO_PUBLIC_API_URL,
  ios: process.env.EXPO_PUBLIC_API_URL,
};

const axiosInstance = axios.create({
  adapter: "fetch",
  baseURL: Platform.OS === "ios" ? baseUrl.ios : baseUrl.android,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    // "X-User-Id": "",
  },
});

export { baseUrl, axiosInstance };
