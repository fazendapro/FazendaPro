import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': '/src'
      }
    },
    define: {
      'process.env.VITE_API_URL': JSON.stringify(env.VITE_API_URL)
    }
  }
})
