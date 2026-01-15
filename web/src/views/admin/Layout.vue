<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="240px" class="material-aside">
        <div class="logo-area">
          <span class="logo-text">Pay Sandbox</span>
        </div>
        <el-menu
          router
          :default-active="$route.path"
          class="el-menu-vertical-demo material-menu"
          text-color="#5f6368"
          active-text-color="#1a73e8"
          :default-openeds="['wechat-pay']">
          <el-sub-menu index="wechat-pay">
            <template #title>
              <el-icon :size="20"><Wallet /></el-icon>
              <span class="menu-text">微信支付</span>
            </template>
            <el-menu-item index="/admin/merchants">
              <el-icon :size="20"><User /></el-icon>
              <span class="menu-text">商户管理</span>
            </el-menu-item>
            <el-menu-item index="/admin/transactions">
              <el-icon :size="20"><List /></el-icon>
              <span class="menu-text">交易流水</span>
            </el-menu-item>
            <el-menu-item index="/admin/refunds">
              <el-icon :size="20"><RefreshLeft /></el-icon>
              <span class="menu-text">退款流水</span>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header class="material-header">
          <div class="header-content">
            <h2 class="page-title">{{ getPageTitle($route.path) }}</h2>
          </div>
          <div class="header-actions">
            <!-- Placeholder for user profile or settings -->
            <div class="avatar-placeholder">A</div>
          </div>
        </el-header>
        <el-main class="material-main">
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { User, List, RefreshLeft, Wallet } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'

const getPageTitle = (path) => {
  if (path.includes('merchants')) return '商户管理'
  if (path.includes('transactions')) return '交易流水'
  if (path.includes('refunds')) return '退款流水'
  return '控制台'
}
</script>

<style scoped>
.common-layout {
  height: 100vh;
  background-color: #f8f9fa;
}
.el-container {
  height: 100%;
}
.material-aside {
  background-color: #fff;
  border-right: 1px solid #dadce0;
  display: flex;
  flex-direction: column;
}
.logo-area {
  height: 64px;
  display: flex;
  align-items: center;
  padding-left: 24px;
  border-bottom: 1px solid transparent;
}
.logo-text {
  font-size: 20px;
  font-weight: 500;
  color: #202124;
}
.material-menu {
  border-right: none;
  padding-top: 12px;
}
.el-menu-item, :deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
  margin: 0 12px 4px 12px;
  border-radius: 8px;
}
:deep(.el-sub-menu .el-menu-item) {
  margin-left: 24px;
  margin-right: 12px;
  width: auto;
}
.el-menu-item.is-active {
  background-color: #e8f0fe;
  color: #1a73e8;
  font-weight: 500;
}
.el-menu-item:hover {
  background-color: #f1f3f4;
}
.menu-text {
  font-size: 16px;
  margin-left: 12px;
}
.material-header {
  background-color: #fff;
  height: 64px;
  border-bottom: 1px solid #dadce0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3), 0 2px 6px 2px rgba(60,64,67,0.15);
  z-index: 10;
}
.page-title {
  font-size: 22px;
  font-weight: 400;
  color: #202124;
  margin: 0;
}
.avatar-placeholder {
  width: 36px;
  height: 36px;
  background-color: #1a73e8;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  cursor: pointer;
}
.material-main {
  padding: 32px;
  background-color: #f8f9fa;
}
</style>
