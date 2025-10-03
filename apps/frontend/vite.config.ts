import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0', // Bind to all interfaces for Docker
    port: 5173,
    allowedHosts: [
      'clan-organised-genealogy-powerseller.trycloudflare.com',
      '.trycloudflare.com', // Allow all cloudflare tunnel domains
      'localhost',
      '127.0.0.1'
    ],
    proxy: {
      // Proxy API requests to backend when API_URL is not available
      '/api': 'http://localhost:8080',
      '/auth': 'http://localhost:8080',
      '/user': 'http://localhost:8080',
      '/keys': 'http://localhost:8080',
      '/credits': 'http://localhost:8080',
      '/models': 'http://localhost:8080',
      '/v1': 'http://localhost:8080',
      '/analytics': 'http://localhost:8080',
      '/chathistory': 'http://localhost:8080',
      '/newchat': 'http://localhost:8080'
    }
  }
})
