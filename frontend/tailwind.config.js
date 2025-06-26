/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{svelte,js,ts,jsx,tsx}"],
  // theme: {
  //   extend: {
  //     colors: {
  //       primary: "#4F46E5", // Indigo
  //       secondary: "#6B7280", // Gray
  //       accent: "#10B981", // Emerald
  //       neutral: "#1F2937", // Gray-800
  //       "base-100": "#FFFFFF",
  //       "base-200": "#F9FAFB",
  //       "base-300": "#F3F4F6",
  //       info: "#3ABFF8",
  //       success: "#34D399",
  //       warning: "#FBBF24",
  //       error: "#F87171",
  //     },
  //   },
  // },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light"],
  },
};
