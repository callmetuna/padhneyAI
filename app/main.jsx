import React, { useState } from "react";
import { View, Text, TouchableOpacity, ScrollView, TextInput, FlatList, Dimensions } from "react-native";
import { useRouter } from "expo-router";
import * as DocumentPicker from "expo-document-picker";
import RNPickerSelect from "react-native-picker-select";
import Ionicons from "react-native-vector-icons/Ionicons"; // Import the icons

const SidebarLayout = () => {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("chat");
  const [searchQuery, setSearchQuery] = useState("");
  const [messages, setMessages] = useState([
    { id: "1", text: "Hello, how can I assist you today?" },
    { id: "2", text: "I am a virtual assistant! Upload your PDF or file and ask me any question about it." },
  ]);
  const [selectedModel, setSelectedModel] = useState("Model A");
  const [isDropdownVisible, setIsDropdownVisible] = useState(false);

  const handleSendMessage = (message) => {
    if (message.trim()) {
      setMessages([...messages, { id: String(messages.length + 1), text: message }]);
    }
  };

  const handleAttachFile = async () => {
    try {
      const result = await DocumentPicker.getDocumentAsync({});
      if (result.type === "success") {
        setMessages([...messages, { id: String(messages.length + 1), text: `File attached: ${result.name}` }]);
      }
    } catch (error) {
      console.error("Error picking file:", error);
    }
  };

  const handleNewChat = () => {
    setMessages([]);
  };

  const windowWidth = Dimensions.get("window").width;
  const windowHeight = Dimensions.get("window").height;

  return (
    <View style={{ width: windowWidth, height: windowHeight }} className="relative bg-gray-900">
      {/* Log Button at Top Right */}
      <TouchableOpacity
        onPress={() => console.log("Log button clicked")}
        className="absolute top-4 right-4 bg-blue-600 py-1 px-3 rounded-md shadow-md"
      >
        <Text className="text-white font-semibold text-xs">Log</Text>
      </TouchableOpacity>

      <View className="flex-row h-full">
        {/* Sidebar */}
        <View className="w-1/4 bg-gray-800 text-white h-full flex flex-col justify-between">
          <ScrollView contentContainerStyle={{ flexGrow: 1 }}>
            {/* Login and Signup Buttons */}
            <View className="flex-row justify-start py-4 px-4">
              <TouchableOpacity
                onPress={() => console.log("Login clicked")}
                className="bg-blue-600 py-1 px-3 rounded-md mr-2"
              >
                <Text className="text-white font-semibold text-xs">Login</Text>
              </TouchableOpacity>

              <TouchableOpacity
                onPress={() => console.log("Signup clicked")}
                className="bg-green-600 py-1 px-3 rounded-md"
              >
                <Text className="text-white font-semibold text-xs">Signup</Text>
              </TouchableOpacity>
            </View>

            {/* ChatGPT Title */}
            <View className="py-6 px-4 border-b border-gray-700 flex-row justify-between items-center">
              <Text className="text-xl font-bold text-gray-100">PandneyAI</Text>
              {/* New Chat Button (with Icon) */}
              <TouchableOpacity
                onPress={handleNewChat}
                className="bg-blue-600 py-1 px-3 rounded-md shadow-md flex-row items-center"
              >
                <Ionicons name="chatbox-ellipses" size={16} color="white" />
                <Text className="text-white font-semibold text-xs ml-1">New Chat</Text>
              </TouchableOpacity>
            </View>

            {/* Search Bar in Sidebar */}
            <View className="py-4 px-4">
              <TextInput
                value={searchQuery}
                onChangeText={setSearchQuery}
                placeholder="Search..."
                placeholderTextColor="#bbb"
                className="w-full p-2 bg-gray-700 text-white rounded-md"
              />
            </View>

            {/* Tabs for Chat */}
            <TouchableOpacity
              className={`py-4 px-4 ${
                activeTab === "chat" ? "bg-gray-700" : "bg-transparent"
              }`}
              onPress={() => setActiveTab("chat")}
            >
              <Text
                className={`text-gray-200 ${
                  activeTab === "chat" ? "font-bold" : "font-normal"
                }`}
              >
                Chat
              </Text>
            </TouchableOpacity>

            {/* Chat History */}
            <View className="py-4 px-4 border-t border-gray-700">
              <Text className="text-gray-200 text-sm font-semibold">Chat History</Text>
              <FlatList
                data={messages}
                renderItem={({ item }) => (
                  <TouchableOpacity
                    onPress={() => {
                      setMessages([item]);
                      setActiveTab("chat");
                    }}
                    className="py-2 px-4 mt-2 bg-gray-700 rounded-md"
                  >
                    <Text className="text-gray-200 text-xs">{item.text}</Text>
                  </TouchableOpacity>
                )}
                keyExtractor={(item) => item.id}
                style={{ maxHeight: 150 }}
              />
            </View>
          </ScrollView>

          {/* Settings at the bottom */}
          <TouchableOpacity
            className={`py-4 px-4 ${
              activeTab === "settings" ? "bg-gray-700" : "bg-transparent"
            }`}
            onPress={() => setActiveTab("settings")}
          >
            <Text
              className={`text-gray-200 ${
                activeTab === "settings" ? "font-bold" : "font-normal"
              }`}
            >
              Settings
            </Text>
          </TouchableOpacity>
        </View>

        {/* Main Content Area */}
        <View className="flex-1 bg-gray-100 h-full">
          {/* Chat Interface */}
          {activeTab === "chat" && (
            <View className="flex-1 p-4 bg-gray-100">
              {/* Select Model Button (with Icon) */}
              <TouchableOpacity
                onPress={() => setIsDropdownVisible(!isDropdownVisible)}
                className="w-32 py-1 px-2 bg-gray-700 text-white rounded-md text-xs flex-row items-center"
              >
                <Ionicons name="settings-sharp" size={16} color="white" />
              </TouchableOpacity>

              {/* Conditionally render the dropdown */}
              {isDropdownVisible && (
                <View className="mt-2 bg-gray-700 rounded-md">
                  <RNPickerSelect
                    onValueChange={(value) => {
                      setSelectedModel(value);
                      setIsDropdownVisible(false);
                    }}
                    value={selectedModel}
                    items={[
                      { label: "Model A", value: "Model A" },
                      { label: "Model B", value: "Model B" },
                      { label: "Model C", value: "Model C" },
                    ]}
                    style={{
                      inputIOS: {
                        backgroundColor: "#333",
                        borderRadius: 5,
                        padding: 10,
                        color: "#fff",
                      },
                      inputAndroid: {
                        backgroundColor: "#333",
                        borderRadius: 5,
                        padding: 10,
                        color: "#fff",
                      },
                    }}
                  />
                </View>
              )}

              {/* Chat Box */}
              <View className="flex-1 p-3 bg-gray-200 rounded-lg">
                <FlatList
                  data={messages}
                  renderItem={({ item }) => (
                    <View
                      className={`mb-2 p-3 rounded-lg ${
                        item.id % 2 === 0 ? "bg-gray-300" : "bg-blue-300"
                      }`}
                    >
                      <Text className="text-gray-800">{item.text}</Text>
                    </View>
                  )}
                  keyExtractor={(item) => item.id}
                  style={{ flexGrow: 0 }}
                />
              </View>

              {/* Input Box for Chat */}
              <View className="flex-row items-center mt-4 border-t border-gray-300 pt-2">
                {/* Attach File Button (with Icon) */}
                <TouchableOpacity
                  onPress={handleAttachFile}
                  className="bg-gray-700 p-3 rounded-md flex-row items-center"
                >
                  <Ionicons name="attach" size={20} color="white" />
                </TouchableOpacity>

                <TextInput
                  placeholder="Type a message..."
                  className="flex-1 p-2 ml-2 border border-gray-300 rounded-md"
                  onSubmitEditing={(e) => handleSendMessage(e.nativeEvent.text)}
                />
                {/* Send Button (with Icon) */}
                <TouchableOpacity
                  onPress={() => handleSendMessage(searchQuery)}
                  className="ml-2 bg-blue-500 p-2 rounded-md flex-row items-center"
                >
                  <Ionicons name="send" size={20} color="white" />
                </TouchableOpacity>
              </View>
            </View>
          )}

          {/* Settings Content */}
          {activeTab === "settings" && (
            <View className="flex-1 p-6 bg-gray-100">
              <Text className="text-lg font-semibold text-gray-800">Settings</Text>
              <Text className="text-gray-600 mt-2">Configure your preferences here.</Text>
            </View>
          )}
        </View>
      </View>
    </View>
  );
};

export default SidebarLayout;
