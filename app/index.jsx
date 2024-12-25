import { Text, View, Button } from "react-native";
import React, { useState } from "react";
import SignUp from "./sign-up"; // Adjust the path if necessary

export default function Index() {
  const [showSignUp, setShowSignUp] = useState(false);

  const handleSignUpClick = () => {
    setShowSignUp(true);
  };

  const handleBackClick = () => {
    setShowSignUp(false);
  };

  return (
    <View
      style={{
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      {!showSignUp ? (
        <>
          <Text>Edit app/index.tsx to edit this screen.</Text>
          <Button title="Sign Up" onPress={handleSignUpClick} />
        </>
      ) : (
        <SignUpScreen onBack={handleBackClick} />
      )}
    </View>
  );
}

function SignUpScreen({ onBack }) {
  return (
    <View
      style={{
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
        padding: 20,
      }}
    >
      <SignUp />
      <Button title="Back" onPress={onBack} />
    </View>
  );
}
