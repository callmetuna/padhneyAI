import { Text, View, Button, Alert } from "react-native";
import React, { useState } from "react";
import SignUp from "./sign-up"; // Adjust the path if necessary
import { StyleSheet } from "react-native";

export default function Index() {
  const [showSignUp, setShowSignUp] = useState(false);

  const handleSignUpClick = () => setShowSignUp(true);
  const handleBackClick = () => setShowSignUp(false);

  return (
    <View style={styles.container}>
      {!showSignUp ? (
        <>
          <View style={styles.welcomeContainer}>
            <Text style={styles.title}>Welcome to PandeAI</Text>
            <Text style={styles.description}>
              PandeAI is an advanced AI platform designed to simplify your tasks
              and empower your creativity.
            </Text>
          </View>
          <Button title="Sign Up" onPress={handleSignUpClick} />
          <View style={styles.buttonSpacing}>
            <Button
              title="Try PandeAI"
              onPress={() => Alert.alert("Try PandeAI clicked!")}
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
    <View style={styles.container}>
      <SignUp />
      <View style={styles.buttonSpacing}>
        <Button title="Back" onPress={onBack} />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#f7f7f7",
  },
  welcomeContainer: {
    alignItems: "center",
    marginBottom: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginTop: 10,
  },
  description: {
    textAlign: "center",
    color: "#6b7280",
    marginTop: 10,
    paddingHorizontal: 16,
  },
  buttonSpacing: {
    marginTop: 16,
  },
});

