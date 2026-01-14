<template>
  <div class="pay-container">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>正在处理支付...</p>
    </div>
    
    <div v-else-if="success" class="result success">
      <div class="icon-box">
        <i class="wechat-icon-success"></i>
      </div>
      <h3>支付成功</h3>
      <p class="amount-text">¥ {{ (amount / 100).toFixed(2) }}</p>
      <button class="btn-primary ripple" @click="close">完成</button>
    </div>
    
    <div v-else class="cashier">
      <div class="merchant-info">
        <p class="label">付款给商家</p>
        <h3 class="merchant-name">{{ merchantName }}</h3>
      </div>
      <div class="amount-display">
        <span class="currency">¥</span>
        <span class="value">{{ (amount / 100).toFixed(2) }}</span>
      </div>
      
      <div class="actions">
        <button class="btn-primary ripple" style="padding: 10px 24px;" @click="showPassword = true">立即支付</button>
      </div>
    </div>

    <!-- 模拟密码输入框 -->
    <div v-if="showPassword" class="password-modal" @click.self="showPassword = false">
      <div class="keyboard-sheet slide-up">
        <div class="sheet-header">
          <span class="close-btn" @click="showPassword = false">×</span>
          <span class="title">请输入支付密码</span>
        </div>
        <div class="pwd-display">
          <div class="pwd-input-box">
            <div v-for="i in 6" :key="i" class="dot" :class="{ filled: password.length >= i }"></div>
          </div>
        </div>
        <div class="keyboard">
          <div v-for="n in 9" :key="n" class="key ripple" @click="inputPwd(n)">{{ n }}</div>
          <div class="key empty"></div>
          <div class="key ripple" @click="inputPwd(0)">0</div>
          <div class="key delete ripple" @click="deletePwd">
            <span class="backspace-icon">⌫</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const prepayId = route.params.prepay_id
const isEmbedded = route.query.embedded === 'true'

const loading = ref(false)
const success = ref(false)
const showPassword = ref(false)
const password = ref('')
const amount = ref(100) // 默认1元
const merchantName = ref('模拟商户')

const inputPwd = (num) => {
  if (password.value.length < 6) {
    password.value += num
    if (password.value.length === 6) {
      doPay()
    }
  }
}

const deletePwd = () => {
  password.value = password.value.slice(0, -1)
}

const doPay = async () => {
  loading.value = true
  showPassword.value = false
  try {
    await axios.post('/api/internal/simulate/pay', { prepay_id: prepayId })
    setTimeout(() => {
      loading.value = false
      success.value = true
    }, 1000)
  } catch (error) {
    alert('支付失败: ' + error.message)
    loading.value = false
    password.value = ''
  }
}

const close = () => {
  if (isEmbedded) {
    window.parent.location.reload()
  } else {
    window.close()
  }
}

onMounted(async () => {
  try {
    const res = await axios.get('/api/internal/transactions', { params: { prepay_id: prepayId } })
    if (res.data && res.data.length > 0) {
      amount.value = res.data[0].amount
      merchantName.value = res.data[0].description || '模拟商户'
    }
  } catch (e) {
    console.error(e)
  }
})
</script>

<style scoped>
.pay-container {
  background-color: #ededed;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #07c160;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.cashier {
  width: 100%;
  padding-top: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.merchant-info {
  text-align: center;
  margin-bottom: 20px;
}

.label {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

.merchant-name {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.amount-display {
  margin: 30px 0 50px;
  display: flex;
  align-items: baseline;
  justify-content: center;
}

.currency {
  font-size: 32px;
  font-weight: 600;
  margin-right: 4px;
}

.value {
  font-size: 56px;
  font-weight: 700;
  line-height: 1;
}

.btn-primary {
  background-color: #07c160;
  color: white;
  border: none;
  padding: 14px 0;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  width: 80%;
  max-width: 300px;
  transition: background-color 0.2s;
  box-shadow: 0 2px 8px rgba(7, 193, 96, 0.3);
}

.btn-primary:active {
  background-color: #06ad56;
}

/* Password Modal */
.password-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.6);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  z-index: 100;
}

.keyboard-sheet {
  background: #f7f7f7;
  padding-bottom: env(safe-area-inset-bottom);
  border-radius: 16px 16px 0 0;
  overflow: hidden;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

.sheet-header {
  background: #fff;
  border-bottom: 1px solid #eee;
  padding: 16px;
  text-align: center;
  position: relative;
}

.close-btn {
  position: absolute;
  left: 16px;
  top: 12px;
  font-size: 28px;
  color: #333;
  cursor: pointer;
  line-height: 1;
}

.title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.pwd-display {
  background: #fff;
  padding: 30px 0;
  display: flex;
  justify-content: center;
}

.pwd-input-box {
  display: flex;
  border: 1px solid #ccc;
  border-radius: 6px;
  width: 90%;
  max-width: 360px;
  height: 54px;
  background: #fff;
}

.dot {
  flex: 1;
  border-right: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dot:last-child {
  border-right: none;
}

.dot.filled::after {
  content: '';
  width: 12px;
  height: 12px;
  background: #000;
  border-radius: 50%;
}

.keyboard {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: #d2d5db; /* Separator color */
  padding-top: 1px;
}

.key {
  background: #fff;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 500;
  color: #000;
  cursor: pointer;
  user-select: none;
}

.key:active {
  background-color: #e5e5e5;
}

.key.empty {
  background: #eef0f4;
  cursor: default;
}

.key.delete {
  background: #eef0f4;
  font-size: 22px;
}

.key.delete:active {
  background-color: #dcdfe5;
}

/* Success Result */
.result {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 80px;
  width: 100%;
}

.icon-box {
  margin-bottom: 24px;
}

.wechat-icon-success {
  display: inline-block;
  width: 100px;
  height: 100px;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 50 50' enable-background='new 0 0 50 50'%3E%3Cpath fill='%2307C160' d='M25,50C11.2,50,0,38.8,0,25S11.2,0,25,0s25,11.2,25,25S38.8,50,25,50z'/%3E%3Cpath fill='%23FFFFFF' d='M21.3,36.6c-0.8,0-1.6-0.3-2.2-0.9L9.5,26.1c-1.2-1.2-1.2-3.2,0-4.4c1.2-1.2,3.2-1.2,4.4,0l7.4,7.4l14.8-14.8c1.2-1.2,3.2-1.2,4.4,0c1.2,1.2,1.2,3.2,0,4.4L23.5,35.7C22.9,36.3,22.1,36.6,21.3,36.6z'/%3E%3C/svg%3E") no-repeat center;
  background-size: contain;
}

.result h3 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 16px;
  color: #333;
}

.amount-text {
  font-size: 40px;
  font-weight: 700;
  color: #333;
  margin: 0 0 40px;
}
</style>
