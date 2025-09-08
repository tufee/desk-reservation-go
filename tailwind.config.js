/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/**/*.html",
    "./web/templates/**/*.{html,js}",
    "./web/static/js/**/*.js",
    "./internal/**/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("daisyui"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/forms"),
  ],
  daisyui: {
    themes: ["corporate"],
  },
};
