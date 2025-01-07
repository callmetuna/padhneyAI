import { Text, View, TouchableOpacity } from "react-native";
import { Link } from "expo-router"; // Import Link component for navigation
import "../global.css"; // Ensure Tailwind is configured

export default function Index() {
  return (
    <View className="flex-1 justify-center items-center bg-gradient-to-r from-teal-400 to-blue-500">
      <View className="bg-blue-100 bg-opacity-90 rounded-xl shadow-xl w-4/5 h-4/5 p-8">
        <View className="flex-1 justify-center items-center">
          <Text className="text-3xl font-extrabold text-gray-800 mb-4 text-center">
            Welcome to PandeAI
          </Text>
          <Text className="text-lg text-gray-700 text-center mb-8">
            PandeAI is an advanced AI platform designed to simplify your tasks
            and empower your creativity.
          </Text>
          <View className="space-y-4">
            {/* Navigate to Sign Up page */}
            <Link
              href="/sign-up" // Linking to the signup page
              className="bg-blue-500 py-3 px-5 rounded-md shadow-sm"
            >
              <Text className="text-white text-sm font-semibold">Sign Up</Text>
            </Link>
            {/* Navigate to Main page */}
            <Link
              href="/main" // Linking to the main page
              className="bg-green-600 py-3 px-5 rounded-md shadow-sm"
            >
              <Text className="text-white text-sm font-semibold">
                Try PandeAI
              </Text>
            </Link>
          </View>
        </View>
      </View>
    </View>
  );
}
