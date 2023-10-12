/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
		"./web/views/*.{html,js}",
		"./components/*.templ"
	],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

