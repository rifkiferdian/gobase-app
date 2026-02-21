module.exports = {
  content: ["./templates/**/*.html"],
  theme: {
    extend: {
      fontFamily: {
        display: ["Plus Jakarta Sans", "ui-sans-serif", "system-ui"],
        body: ["Manrope", "ui-sans-serif", "system-ui"],
        "display-auth": ["Space Grotesk", "ui-sans-serif", "system-ui"],
      },
      colors: {
        brand: {
          50: "#f6f1f4",
          100: "#efe2ea",
          300: "#d1a0d6",
          500: "#8c149c",
          600: "#800080",
          700: "#571a3f",
        },
        "brand-auth": {
          50: "#f6f3ff",
          100: "#e9e1ff",
          200: "#d2c1ff",
          300: "#b49bff",
          400: "#9567ff",
          500: "#7b35ff",
          600: "#6c20f2",
          700: "#5718c1",
          800: "#40128f",
          900: "#2d0d65",
        },
      },
      boxShadow: {
        card: "0 12px 26px rgba(16, 24, 40, 0.08)",
        glow: "0 30px 80px -40px rgba(90, 24, 154, 0.6)",
      },
    },
  },
  plugins: [],
};
