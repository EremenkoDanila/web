import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        secure: false, // Добавляем эту строку для локального HTTPS
        rewrite: (path) => path.replace(/^\/api/, "/api"), // Не убираем /api
      },
    },
    watch: {
      usePolling: true,
    }, 
    host: true,
    strictPort: true,
    port: 3000,
  },
});