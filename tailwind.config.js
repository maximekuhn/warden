/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/ui/**/*.{templ,go}"],
    theme: {
        extend: {
            colors: {
                "primary": "#11767D",
                "secondary": "#A37774",
                "accent": "#E88873",
                "neutral": "#484A47"
            },
        },
    },
    plugins: [],
}

