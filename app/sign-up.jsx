import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, StyleSheet } from 'react-native';
import Icon from 'react-native-vector-icons/FontAwesome';

function SignupForm() {
    const [fullname, setFullname] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            Alert.alert('Passwords do not match!');
            return;
        }

        try {
            const response = await fetch('http://localhost:1337/api/auth/local/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    username: fullname,
                    email,
                    password,
                }),
            });
            const data = await response.json();
            if (response.ok) {
                Alert.alert('Signup successful!');
            } else {
                Alert.alert(data.message[0].messages[0].message);
            }
        } catch (error) {
            console.error('Error signing up:', error);
            Alert.alert('Error signing up');
        }
    };

    return (
        <View className="flex-1 justify-center items-center bg-gradient-to-r from-black via-blue-800 to-blue-500">
            <View className="w-4/5 max-w-md bg-white p-8 rounded-2xl shadow-xl">
                <Text className="text-4xl font-bold text-center text-blue-600 mb-6">Sign Up</Text>
                <TextInput
                    className="w-full px-5 py-3 mb-5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="Fullname"
                    value={fullname}
                    onChangeText={(text) => setFullname(text)}
                    required
                />
                <TextInput
                    className="w-full px-5 py-3 mb-5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="Email"
                    value={email}
                    onChangeText={(text) => setEmail(text)}
                    keyboardType="email-address"
                    required
                />
                <TextInput
                    className="w-full px-5 py-3 mb-5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="Password"
                    value={password}
                    onChangeText={(text) => setPassword(text)}
                    secureTextEntry
                    required
                />
                <TextInput
                    className="w-full px-5 py-3 mb-8 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="Confirm Password"
                    value={confirmPassword}
                    onChangeText={(text) => setConfirmPassword(text)}
                    secureTextEntry
                    required
                />
                <TouchableOpacity
                    className="w-full bg-blue-600 py-3 rounded-md text-white text-center mb-4"
                    onPress={handleSubmit}
                >
                    <Text className="font-semibold">Sign Up</Text>
                </TouchableOpacity>

                <View className="flex flex-row justify-center items-center mt-8 space-x-4">
                    <TouchableOpacity className="bg-red-500 py-3 px-8 rounded-md flex-row items-center">
                        <Icon name="google" size={20} color="white" />
                        <Text className="text-white ml-2">Google</Text>
                    </TouchableOpacity>
                    <TouchableOpacity className="bg-blue-700 py-3 px-8 rounded-md flex-row items-center">
                        <Icon name="facebook" size={20} color="white" />
                        <Text className="text-white ml-2">Facebook</Text>
                    </TouchableOpacity>
                </View>

                <View className="mt-6 text-center">
                    <Text className="text-sm">Already have an account? </Text>
                    <TouchableOpacity>
                        <Text className="text-blue-600 font-bold">Sign In</Text>
                    </TouchableOpacity>
                </View>
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#000000', // Black background
    },
    formContainer: {
        width: '90%',
        backgroundColor: '#ffffff',
        padding: 30,
        borderRadius: 12,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.1,
        shadowRadius: 6,
        elevation: 5,
    },
    header: {
        fontSize: 28,
        fontWeight: 'bold',
        textAlign: 'center',
        marginBottom: 20,
        color: '#1E40AF', // Blue header
    },
    input: {
        height: 50,
        borderColor: '#ddd',
        borderWidth: 1,
        borderRadius: 10,
        marginBottom: 15,
        paddingLeft: 15,
        backgroundColor: '#f1f1f1', // Light input background
    },
    button: {
        backgroundColor: '#1E40AF', // Blue button
        paddingVertical: 12,
        borderRadius: 10,
        alignItems: 'center',
        marginBottom: 20,
    },
    buttonText: {
        color: 'white',
        fontSize: 18,
        fontWeight: 'bold',
    },
    socialButtons: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginTop: 20,
    },
    socialButton: {
        flexDirection: 'row',
        paddingVertical: 10,
        paddingHorizontal: 20,
        margin: 5,
        borderRadius: 8,
        alignItems: 'center',
        justifyContent: 'center',
    },
    socialButtonText: {
        color: 'white',
        fontSize: 16,
        marginLeft: 10,
    },
    signinLink: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginTop: 20,
    },
    linkText: {
        color: '#1E40AF', // Blue link
        fontWeight: 'bold',
    },
});

export default SignupForm;
