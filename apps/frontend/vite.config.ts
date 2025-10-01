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
    ]
  }
})
