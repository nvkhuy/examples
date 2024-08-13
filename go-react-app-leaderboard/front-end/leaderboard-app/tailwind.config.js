module.exports = {
    content: ["./src/**/*.{js,jsx,ts,tsx}",], theme: {
        extend: {
            height: {
                '52': '13rem',  // Custom height example
                '60': '15rem',  // Custom height example
            }
        },
    }, plugins: [], extend: {
        colors: {
            yellow: {
                100: '#FCEFC7', // Adjust the color code to match your specific shade
            }, red: {
                100: '#F9D6D2', // Adjust the color code to match your specific shade
            }, gray: {
                100: '#F3F4F6', // Adjust if necessary
            }
        },
    }
}
