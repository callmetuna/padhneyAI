import { create } from 'tailwindcss-config-native';

export default create({
  content: ['./App.{js,jsx,ts,tsx}', './**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [],
});
