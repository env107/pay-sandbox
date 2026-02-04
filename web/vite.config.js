import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 后端端口：环境变量 BACKEND_PORT 指定，未设置时默认 8080（与 Go 服务默认一致）
const backendPort = process.env.BACKEND_PORT || '8080'
const backendTarget = `http://localhost:${backendPort}`

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: backendTarget,
        changeOrigin: true
      },
      '/v3': {
        target: backendTarget,
        changeOrigin: true
      }
    }
  }
})
