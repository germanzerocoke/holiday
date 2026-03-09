import { deleteItemAsync, getItemAsync, setItemAsync } from "expo-secure-store";

async function saveSecureStore(key: string, value: string): Promise<void> {
  await setItemAsync(key, value);
}

async function getSecureStore(key: string) {
  const value = await getItemAsync(key);
  return value;
}

async function deleteSecureStore(key: string) {
  await deleteItemAsync(key);
}

export { saveSecureStore, getSecureStore, deleteSecureStore };
