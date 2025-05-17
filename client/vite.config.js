import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath } from 'url';

// Create __dirname equivalent for ESM
const __dirname = path.dirname(fileURLToPath(import.meta.url));


export default defineConfig({
  plugins: [react(), tailwindcss()],

  mode: "production",
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
    sourcemap: false,

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

  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },

  // Server configuration for development
  server: {
    proxy: {
      "/query": {
        target: "http://127.0.0.1:7654",
      },
      "/families": {
        target: "http://127.0.0.1:7654",
      },
      "/completions": {
        target: "http://127.0.0.1:5150",
      },
    },
    fs: {
      allow: [
        path.resolve(__dirname, "./"),
        path.resolve(__dirname, "../node_modules/pdfjs-dist"),
      ]
    }
  },
});
