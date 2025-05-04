import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],

  // Set the base directory for the build
  base: "./",

  // Configure the build output
  build: {
    // Output directory for the build (relative to vite.config.js)
    outDir: "../cmd/dashboard/web",
    emptyOutDir: true,
    // Make sure assets are placed in the same directory
    assetsDir: "",

    // For smaller builds, disable source maps in production
    sourcemap: process.env.NODE_ENV !== "production",

    // Configure the output format
    rollupOptions: {
      input: {
        main: path.resolve(__dirname, "./index.html"),
      },
      output: {
        entryFileNames: "assets/[name].[hash].js",
        chunkFileNames: "assets/[name].[hash].js",
        assetFileNames: "assets/[name].[hash].[ext]",
      },
    },
  },

  // Server configuration for development
  server: {
    port: 3000,
    open: true,
  },
});
