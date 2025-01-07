import { Text, View, Animated, Easing, Dimensions } from "react-native";
import { Link } from "expo-router";
import { useEffect, useRef } from "react";
import "../global.css";

export default function Index() {
  const animatedScale = useRef(new Animated.Value(1)).current;

  // Get screen dimensions
  const { width, height } = Dimensions.get("window");
  const boxWidth = width * 0.8; // 80% of the screen width
  const boxHeight = height * 0.8; // 80% of the screen height

  useEffect(() => {
    // Smooth pulse animation for the box
    const pulseAnimation = Animated.loop(
      Animated.sequence([
        Animated.timing(animatedScale, {
          toValue: 1.05,
          duration: 1000,
          easing: Easing.inOut(Easing.ease),
          useNativeDriver: true,
        }),
        Animated.timing(animatedScale, {
          toValue: 1,
          duration: 1000,
          easing: Easing.inOut(Easing.ease),
          useNativeDriver: true,
        }),
      ])
    );
    pulseAnimation.start();
  }, [animatedScale]);

  return (
    <View className="flex-1 justify-center items-center bg-gradient-to-r from-teal-400 to-blue-500">
      {/* Animated Box */}
      <Animated.View
        style={{
          width: boxWidth,
          height: boxHeight,
          transform: [{ scale: animatedScale }],
        }}
        className="rounded-xl shadow-xl"
      >
        <View className="flex-1 justify-center items-center rounded-xl bg-gradient-to-r from-blue-700 to-black p-8">
          <Text className="text-3xl font-extrabold text-gray-100 mb-4 text-center">
            Welcome to PandneyAI
          </Text>
          <Text className="text-lg text-gray-200 text-center mb-8">
            PandneyAI is an advanced AI platform designed to simplify your tasks
            and empower your creativity.
          </Text>
          <View className="flex-row space-x-4">
            {/* Sign Up Button */}
            <Link
              href="/sign-up"
              className="bg-blue-500 py-3 px-5 rounded-md shadow-sm border border-gray-300"
            >
              <Text className="text-white text-sm font-semibold">Sign Up</Text>
            </Link>
            {/* Try PandeAI Button */}
            <Link
              href="/main"
              className="bg-green-600 py-3 px-5 rounded-md shadow-sm border border-gray-300"
            >
              <Text className="text-white text-sm font-semibold">
                Try PandneyAI
              </Text>
            </Link>
          </View>
        </View>
      </Animated.View>
    </View>
  );
}
