/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/ui/**/*.{templ,go}"],
    theme: {
        extend: {
            colors: {
                "primary": "#11767D",
                "secondary": "#00353D",
                "accent": "#E88873",
                "error": "#ED4337",
                "neutral": "#227F92"
            },
        },
    },
    plugins: [],
}

