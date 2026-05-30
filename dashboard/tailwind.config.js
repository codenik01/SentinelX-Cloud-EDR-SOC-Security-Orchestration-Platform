/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        cyber: {
          bg: "#0A0F1D",         // Ultra-dark background
          card: "rgba(16, 22, 42, 0.65)", // Glassmorphic card fill
          border: "rgba(40, 50, 80, 0.4)", // Cyber border
          glow: "#3B82F6",       // Main neon blue glow
          critical: "#EF4444",   // Severe red
          high: "#F97316",       // Warning orange
          medium: "#F59E0B",     // Amber SSH alert
          low: "#10B981",        // Green success/low alert
          muted: "#94A3B8"       // Gray subtitles
        }
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"]
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow-pulse': 'glow 2s ease-in-out infinite alternate',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px rgba(59, 130, 246, 0.2)' },
          '100%': { boxShadow: '0 0 15px rgba(59, 130, 246, 0.6)' }
        }
      }
    },
  },
  plugins: [],
}
