<template>
  <div class="material-card">
    <div class="table-header">
      <h3 class="card-title">商户列表</h3>
      <div class="header-actions">
        <el-button type="danger" size="large" :disabled="selectedRows.length === 0" @click="handleBatchDelete">
          批量删除
        </el-button>
        <el-button type="primary" size="large" class="fab-style" @click="dialogVisible = true">
          <el-icon class="el-icon--left"><Plus /></el-icon>添加商户
        </el-button>
      </div>
    </div>

    <el-table :data="tableData" style="width: 100%" size="large" :header-cell-style="{ background: '#f8f9fa', color: '#5f6368', fontWeight: 500 }" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="mchid" label="商户号 (MchID)" min-width="140" />
      <el-table-column prop="appid" label="AppID" min-width="160" />
      <el-table-column prop="description" label="备注" min-width="180" />
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="scope">
          <el-button link type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button link type="primary" @click="openSimulatePay(scope.row)">模拟支付请求</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 商户配置弹窗 -->
    <el-dialog v-model="dialogVisible" title="商户配置" width="560px" top="8vh">
      <el-form :model="form" label-position="top" size="large">
        <el-form-item label="商户号 (MchID)">
          <el-input v-model="form.mchid" placeholder="请输入微信支付商户号" />
        </el-form-item>
        <el-form-item label="AppID">
          <el-input v-model="form.appid" placeholder="请输入关联的 AppID" />
        </el-form-item>
        <el-form-item label="API v3 Key">
          <el-input v-model="form.api_v3_key" placeholder="请输入 32 位 API v3 密钥" show-password />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.description" placeholder="例如：测试环境商户" />
        </el-form-item>
        <el-form-item label="支付回调地址">
          <el-input v-model="form.notify_url" placeholder="http://localhost:8080/notify" />
        </el-form-item>
        <el-form-item label="退款回调地址">
          <el-input v-model="form.refund_notify_url" placeholder="http://localhost:8080/notify/refund" />
        </el-form-item>
        <el-form-item label="回调间隔 (Duration)">
          <el-input v-model="form.interval" placeholder="例如: 1m, 5s" />
        </el-form-item>
        <el-form-item label="最大重试次数">
          <el-input-number v-model="form.max_retries" :min="0" :max="10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" size="large">取消</el-button>
          <el-button type="primary" @click="saveMerchant" size="large">保存配置</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 模拟支付弹窗 -->
    <el-dialog v-model="simulateVisible" title="模拟微信支付" width="420px" top="15vh">
      <div v-if="!currentMerchant.notify_url" class="no-notify-url">
        <el-alert title="请先配置默认回调地址" type="warning" :closable="false" show-icon description="模拟支付成功后需要回调该地址通知结果。" />
        <div style="margin-top: 24px; text-align: center;">
          <el-button type="primary" size="large" @click="simulateVisible = false; handleEdit(currentMerchant)">去配置</el-button>
        </div>
      </div>
      <div v-else class="simulate-box">
        <el-form label-position="top" size="large">
          <el-form-item label="支付类型">
            <el-select v-model="paymentType" placeholder="请选择支付类型" style="width: 100%">
              <el-option label="JSAPI" value="WX:JSAPI" />
            </el-select>
          </el-form-item>
          <el-form-item label="支付金额 (元)">
            <el-input-number v-model="payAmount" :precision="2" :step="0.01" :min="0.01" style="width: 100%" />
          </el-form-item>
          <div class="simulate-tip">
            将向 <strong>{{ currentMerchant.notify_url }}</strong> 发送回调
          </div>
          <el-button type="success" size="large" style="width: 100%; margin-top: 16px;" @click="startPay" :loading="creatingOrder">
            发起支付 (新窗口)
          </el-button>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const tableData = ref([])
const selectedRows = ref([])
const dialogVisible = ref(false)
const form = ref({
  mchid: '',
  appid: '',
  api_v3_key: '',
  description: '',
  notify_url: '',
  refund_notify_url: '',
  interval: '1m',
  max_retries: 3
})
const isEdit = ref(false)

// 模拟支付相关
const simulateVisible = ref(false)
const currentMerchant = ref({})
const payAmount = ref(0.01)
const paymentType = ref('WX:JSAPI')
const creatingOrder = ref(false)

const loadData = async () => {
  try {
    const res = await axios.get('/api/internal/merchants')
    tableData.value = res.data
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const handleSelectionChange = (val) => {
  selectedRows.value = val
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(
    `确定要永久删除选中的 ${selectedRows.value.length} 个商户吗？`,
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const ids = selectedRows.value.map(row => row.id)
      await axios.delete('/api/internal/merchants', { data: ids })
      ElMessage.success('删除成功')
      loadData()
    } catch (e) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const saveMerchant = async () => {
  try {
    const payload = {
      ...form.value,
      notify_config: JSON.stringify({
        interval: form.value.interval,
        max_retries: form.value.max_retries
      })
    }
    
    if (isEdit.value) {
      await axios.put(`/api/internal/merchants/${form.value.id}`, payload)
    } else {
      await axios.post('/api/internal/merchants', payload)
    }
    dialogVisible.value = false
    loadData()
    ElMessage.success('保存成功')
    resetForm()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const resetForm = () => {
  form.value = { mchid: '', appid: '', api_v3_key: '', description: '', notify_url: '', refund_notify_url: '', interval: '1m', max_retries: 3 }
  isEdit.value = false
}

const handleEdit = (row) => {
  let config = { interval: '1m', max_retries: 3 }
  try {
    config = JSON.parse(row.notify_config)
  } catch (e) {}
  
  form.value = { 
    ...row,
    interval: config.interval || '1m',
    max_retries: config.max_retries || 3
  }
  isEdit.value = true
  dialogVisible.value = true
}

const openSimulatePay = (row) => {
  currentMerchant.value = row
  simulateVisible.value = true
  payAmount.value = 0.01
  paymentType.value = 'WX:JSAPI'
}

const startPay = async () => {
  creatingOrder.value = true
  try {
    // 调用 mock 下单接口
    const res = await axios.post('/v3/pay/transactions/jsapi', {
      appid: currentMerchant.value.appid,
      mchid: currentMerchant.value.mchid,
      description: '模拟支付测试商品',
      out_trade_no: 'TEST_' + Date.now(),
      notify_url: currentMerchant.value.notify_url,
      trade_type: paymentType.value,
      amount: {
        total: Math.round(payAmount.value * 100), // 转为分
        currency: 'CNY'
      },
      payer: {
        openid: 'mock_openid_123'
      }
    })
    const prepayId = res.data.prepay_id
    window.open(`/pay/preview/${prepayId}`, '_blank', 'width=375,height=667')
    simulateVisible.value = false
  } catch (e) {
    ElMessage.error('下单失败: ' + (e.response?.data?.message || e.message))
  } finally {
    creatingOrder.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.card-title {
  font-size: 20px;
  font-weight: 400;
  margin: 0;
  color: #202124;
}
.header-actions {
  display: flex;
  gap: 12px;
}
.fab-style {
  border-radius: 24px;
  padding-left: 20px;
  padding-right: 24px;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3);
}
.no-notify-url {
  padding: 8px;
}
.simulate-tip {
  font-size: 14px;
  color: #5f6368;
  margin-top: 8px;
}
</style>
