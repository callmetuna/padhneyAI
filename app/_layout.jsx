import { Stack } from "expo-router";

export default function RootLayout() {
  return (
    <Stack>
      {/* Home Screen */}
      <Stack.Screen
        name="index"
        options={{ headerShown: false }} // Hide the header for the "index" screen
      />
      {/* Sign Up Screen */}
      <Stack.Screen
        name="sign-up"
        options={{ headerShown: false }} // Customize the title for the "signup" screen
      />
      {/* Main Screen */}
      <Stack.Screen
        name="main"
         // Customize the title for the "main" screen
        options={{ headerShown: false }}
      />
    </Stack>
  );
}
