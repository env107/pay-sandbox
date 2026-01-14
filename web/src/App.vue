<template>
  <router-view></router-view>
</template>

<script setup>
import { onMounted } from 'vue'
import { ElNotification } from 'element-plus'

onMounted(() => {
  const eventSource = new EventSource('/api/internal/events')
  
  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.type === 'callback') {
        ElNotification({
          title: '回调通知',
          message: data.payload.message,
          type: 'success',
          duration: 5000,
          position: 'top-right',
        })
      }
    } catch (e) {
      console.error('SSE Parse Error:', e)
    }
  }

  eventSource.onerror = (error) => {
    console.error('SSE Error:', error)
    // eventSource.close() // 可选：出错时关闭连接或尝试重连
  }
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap');

:root {
  /* Google Blue Theme */
  --el-color-primary: #1a73e8;
  --el-color-primary-light-3: #4285f4;
  --el-color-primary-light-5: #8ab4f8;
  --el-color-primary-light-7: #aecbfa;
  --el-color-primary-light-9: #e8f0fe;
  --el-color-primary-dark-2: #1557b0;
  
  /* Increased Base Sizes */
  --el-font-size-base: 16px;
  --el-font-size-small: 14px;
  --el-font-size-large: 18px;
  --el-border-radius-base: 8px;
  --el-border-radius-small: 4px;
  --el-border-radius-round: 20px;
}

body {
  margin: 0;
  font-family: 'Roboto', -apple-system, BlinkMacSystemFont, "Segoe UI", "Helvetica Neue", Arial, sans-serif;
  background-color: #f8f9fa; /* Material Light Gray Background */
  color: #202124;
  -webkit-font-smoothing: antialiased;
}

/* Global Component Overrides for "Bigger" feel */
.el-button {
  height: 40px; /* Default larger height */
  font-weight: 500;
  padding: 10px 24px;
}

.el-button--small {
  height: 32px;
  padding: 8px 16px;
  font-size: 14px;
}

.el-button--large {
  height: 48px;
  padding: 12px 28px;
  font-size: 18px;
}

.el-input__wrapper {
  padding: 4px 12px;
  box-shadow: 0 0 0 1px #dadce0 inset;
}

.el-input__wrapper.is-focus {
  box-shadow: 0 0 0 2px var(--el-color-primary) inset !important;
}

.el-dialog {
  border-radius: 12px;
  box-shadow: 0 24px 38px 3px rgba(0,0,0,0.14), 0 9px 46px 8px rgba(0,0,0,0.12), 0 11px 15px -7px rgba(0,0,0,0.2);
}

.el-dialog__header {
  padding: 24px 24px 10px;
  margin-right: 0;
  border-bottom: none;
}

.el-dialog__title {
  font-size: 22px;
  font-weight: 500;
}

.el-dialog__body {
  padding: 10px 24px 30px;
}

.el-dialog__footer {
  padding: 10px 24px 24px;
}

/* Card Style Container */
.material-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3), 0 1px 3px 1px rgba(60,64,67,0.15);
  padding: 24px;
  margin-bottom: 24px;
}
</style>
