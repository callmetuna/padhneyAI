import { Text, View, Button, Image } from "react-native";
import React, { useState } from "react";
import SignUp from "./sign-up"; // Adjust the path if necessary
import "./output.css"

export default function Index() {
  const [showSignUp, setShowSignUp] = useState(false);

  const handleSignUpClick = () => setShowSignUp(true);
  const handleBackClick = () => setShowSignUp(false);

  return (
    <View className="flex-1 justify-center items-center bg-gray-100">
      {!showSignUp ? (
        <>
          <View className="items-center mb-5">
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
              onPress={() => alert("Try PandeAI clicked!")}
              color="#4CAF50"
            />
          </View>
        </>
      ) : (
        <SignUpScreen onBack={handleBackClick} />
      )}
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
