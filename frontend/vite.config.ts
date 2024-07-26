import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/get-challenge": "http://localhost:8080/",
      "/follow-challenge": "http://localhost:8080/",
    },
  },
});
