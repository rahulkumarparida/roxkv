import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import process from 'process'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const port = Number(env.PORT || env.FRONTEND_PORT || 5173)
  const usePolling = env.CHOKIDAR_USEPOLLING === 'true'
  const apiBaseUrl = env.VITE_API_BASE_URL || env.VITE_BACKEND_SSE_TARGET || 'http://localhost:6971'
  const chatApiBaseUrl = env.VITE_CHAT_API_BASE_URL || env.VITE_BACKEND_CHAT_TARGET || 'http://localhost:6972'

  return {
    plugins: [react(), tailwindcss()],
    server: {
      host: '0.0.0.0',
      port,
      strictPort: true,
      watch: {
        usePolling,
      },
      proxy: {
        '/api/events': {
          target: apiBaseUrl,
          changeOrigin: true,
        },
        '/api/chat': {
          target: chatApiBaseUrl,
          changeOrigin: true,
        },
      },
    },
    preview: {
      host: '0.0.0.0',
      port,
      strictPort: true,
    },
  }
})
