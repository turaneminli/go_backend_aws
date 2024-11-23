module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"], // Ensure this is correct
  theme: {
    extend: {
      colors: {
        orange: {
          light: "#ffa726",
          DEFAULT: "#fb8c00",
          dark: "#ef6c00",
        },
      },
    },
  },
  plugins: [],
};
