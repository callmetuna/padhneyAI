import { Text, View, Button } from "react-native";
import { React, useState } from "react";
import SignUp from "./sign-up"; // Adjust the path if necessary
import Main from "./main"; // Import Main component
import "../global.css";

export default function Index() {
  const [screen, setScreen] = useState("home"); // Track the current screen

  const handleSignUpClick = () => setScreen("signup");
  const handleTryPandeAIClick = () => setScreen("main");
  const handleBackClick = () => setScreen("home");

  return (
    <View className="flex-1 justify-center items-center bg-opacity-25 bg-teal-400">
      {screen === "home" && (
        <>
          <View className="items-center mb-5  ">
            <Text className="text-lg font-bold mt-2">Welcome to PandeAI</Text>
            <Text className="text-center text-gray-600 mt-1 px-4">
              PandeAI is an advanced AI platform designed to simplify your tasks
              and empower your creativity.
            </Text>
          </View>
          <Button title="Sign Up" onPress={handleSignUpClick} />
          <View className="mt-4">
            <Button
              title="Try PandeAI"
              onPress={handleTryPandeAIClick}
              color="#4CAF50"
            />
          </View>
        </>
      )}
      {screen === "signup" && <SignUpScreen onBack={handleBackClick} />}
      {screen === "main" && <MainScreen onBack={handleBackClick} />}
    </View>
  );
}

function SignUpScreen({ onBack }) {
  return (
    <View className="flex-1 justify-center items-center p-5 bg-gray-100">
      <SignUp />
      <View className="mt-4">
        <Button title="Back" onPress={onBack} />
      </View>
    </View>
  );
}

function MainScreen({ onBack }) {
  return (
    <View className="flex-1 justify-center items-center p-5 bg-gray-100">
      <Main />
      <View className="mt-4">
        <Button title="Back" onPress={onBack} />
      </View>
    </View>
  );
}
