/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./*.html", "./templates/**/*.html"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        "on-surface":"#191c1d","inverse-primary":"#84d6b9","primary-fixed":"#a0f3d4",
        "on-secondary-fixed":"#0d1f1b","surface-container-low":"#f2f4f5",
        "on-tertiary-fixed":"#3b0804","inverse-on-surface":"#eff1f2","outline":"#6f7a74",
        "error-container":"#ffdad6","on-primary-fixed":"#002117","on-secondary":"#ffffff",
        "surface-bright":"#f8fafb","on-primary-container":"#9aedcf","on-tertiary":"#ffffff",
        "on-secondary-container":"#566863","on-tertiary-fixed-variant":"#743329",
        "critical-red":"#EF4444","background":"#f8fafb","success-green":"#22C55E",
        "outline-variant":"#bec9c3","primary-fixed-dim":"#84d6b9","on-surface-variant":"#3f4944",
        "on-background":"#191c1d","surface-tint":"#086b53","surface":"#f8fafb",
        "on-primary-fixed-variant":"#00513e","primary":"#005440","alert-orange":"#F97316",
        "surface-container-high":"#e6e8e9","surface-container":"#eceeef","primary-container":"#0f6e56",
        "on-secondary-fixed-variant":"#394a45","tertiary-container":"#954c41","error":"#ba1a1a",
        "on-primary":"#ffffff","tertiary-fixed":"#ffdad4","surface-border":"#E2E8F0",
        "surface-dim":"#d8dadb","tertiary-fixed-dim":"#ffb4a8","secondary-fixed-dim":"#b7cbc4",
        "surface-container-lowest":"#ffffff","secondary-fixed":"#d3e7e0","surface-variant":"#e1e3e4",
        "on-tertiary-container":"#ffd3cc","text-main":"#1A2E29","on-error-container":"#93000a",
        "inverse-surface":"#2e3132","text-muted":"#64748B","secondary-container":"#d3e7e0",
        "secondary":"#50625d","on-error":"#ffffff","surface-container-highest":"#e1e3e4",
        "tertiary":"#78352b"
      },
      borderRadius: { DEFAULT: "1rem", lg: "2rem", xl: "3rem", full: "9999px" },
      fontFamily: { sans: ["Inter", "sans-serif"] }
    }
  },
  plugins: []
}