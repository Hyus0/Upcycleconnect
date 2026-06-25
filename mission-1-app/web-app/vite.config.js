import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  server: {
    host: "0.0.0.0",
    port: 5174,
    proxy: {
      "/go": {
        target: "http://localhost:8081",
        rewrite: (path) => path.replace(/^\/go/, ""),
        changeOrigin: true,
      },
      "/annonces": "http://localhost:8081"
    }
  }
});
