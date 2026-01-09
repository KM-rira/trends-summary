import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      // GoバックエンドのAPIへプロキシ
      '/github-trending': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ai-repository-summary': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/rss': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/golang-repository-trending': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/tiobe-graph': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ai-article-summary': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/golang-weekly-content': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/google-cloud-content': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/aws-content': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/azure-content': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/google-cloud-content-ja': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/aws-content-ja': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/azure-content-ja': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/rss-ja': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ai-trends-summary': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
