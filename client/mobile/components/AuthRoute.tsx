import { ReactNode } from "react";
import { router, useFocusEffect } from "expo-router";
import { useAuth } from "@/hooks/useAuth";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { auth } = useAuth();

  useFocusEffect(() => {
    !auth.id && router.replace("/auth");
    auth.id && router.replace("/home");
  });
  return <>{children}</>;
}
