<template>
  <div class="material-card">
    <div class="table-header">
      <h3 class="card-title">交易流水</h3>
      <div class="filter-bar">
        <el-input v-model="filter.out_trade_no" placeholder="商户订单号" style="width: 180px" clearable />
        <el-input v-model="filter.prepay_id" placeholder="Prepay ID" style="width: 180px" clearable />
        <el-select v-model="filter.status" placeholder="状态" style="width: 120px" clearable>
          <el-option label="CREATED" value="CREATED" />
          <el-option label="SUCCESS" value="SUCCESS" />
          <el-option label="FAIL" value="FAIL" />
        </el-select>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetFilter">重置</el-button>
        <el-button type="danger" :disabled="selectedRows.length === 0" @click="handleBatchDelete">
          批量删除
        </el-button>
      </div>
    </div>

    <el-table :data="tableData" style="width: 100%" size="large" :header-cell-style="{ background: '#f8f9fa', color: '#5f6368', fontWeight: 500 }" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="created_at" label="创建时间" min-width="200">
        <template #default="scope">
          {{ new Date(scope.row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column prop="mchid" label="商户ID" min-width="160" />
      <el-table-column prop="out_trade_no" label="商户订单号" min-width="220" />
      <el-table-column prop="prepay_id" label="Prepay ID" min-width="220" show-overflow-tooltip />
      <el-table-column prop="trade_type" label="支付类型" min-width="120">
        <template #default="scope">
          {{ tradeTypeMap[scope.row.trade_type] || scope.row.trade_type }}
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额 (分)" min-width="120" align="right" />
      <el-table-column prop="status" label="状态" min-width="120" align="center">
        <template #default="scope">
          <span :class="['status-pill', getStatusClass(scope.row.status)]">{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column label="回调结果" width="120" align="center">
        <template #default="scope">
          <el-tooltip
            v-if="scope.row.callback_status === 'FAIL'"
            class="box-item"
            effect="dark"
            :content="scope.row.callback_msg || '未知错误'"
            placement="top">
            <span class="status-pill status-error">回调失败</span>
          </el-tooltip>
          <span v-else-if="scope.row.callback_status === 'SUCCESS'" class="status-pill status-success">回调成功</span>
          <span v-else class="status-pill status-gray">未回调</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button 
            v-if="scope.row.status === 'CREATED'"
            link
            type="primary" 
            @click="openPreview(scope.row.prepay_id)">
            模拟支付
          </el-button>
          <el-button 
            v-if="scope.row.status === 'SUCCESS'"
            link
            type="warning" 
            @click="openRefund(scope.row)">
            模拟退款
          </el-button>
          <el-button link type="primary" @click="viewDetails(scope.row)">回调日志</el-button>
          <el-button 
            v-if="scope.row.callback_status === 'FAIL'"
            link
            type="danger" 
            @click="retryCallback(scope.row)">
            重试
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 模拟退款弹窗 -->
    <el-dialog v-model="refundVisible" title="模拟退款" width="420px" top="15vh">
      <el-form label-position="top" size="large">
        <el-form-item label="退款金额 (元)">
          <el-input-number v-model="refundAmount" :precision="2" :step="0.01" :min="0.01" :max="currentTx.amount / 100" style="width: 100%" />
          <div class="refund-tip">原订单金额: ¥ {{ (currentTx.amount / 100).toFixed(2) }}</div>
        </el-form-item>
        <el-form-item label="退款原因">
          <el-input v-model="refundReason" placeholder="例如：用户协商退款" />
        </el-form-item>
        <el-button type="danger" size="large" style="width: 100%; margin-top: 16px;" @click="doRefund" :loading="refunding">
          确认退款
        </el-button>
      </el-form>
    </el-dialog>

    <!-- 回调日志弹窗 -->
    <el-dialog v-model="logDialogVisible" title="回调通知日志" width="900px" top="5vh">
      <el-table :data="logData" border style="width: 100%" size="default">
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="scope">
            {{ new Date(scope.row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="retry_count" label="重试次数" width="100" align="center" />
        <el-table-column prop="status_code" label="HTTP 状态" width="100" align="center">
          <template #default="scope">
            <span :class="scope.row.status_code >= 200 && scope.row.status_code < 300 ? 'text-success' : 'text-danger'">
              {{ scope.row.status_code }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="结果" width="100" align="center">
           <template #default="scope">
            <el-tag :type="scope.row.status === 'SUCCESS' ? 'success' : 'danger'" effect="dark" round>
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发送内容 (Request Body)" min-width="200">
          <template #default="scope">
            <pre class="code-block">{{ formatJson(scope.row.request_body) }}</pre>
          </template>
        </el-table-column>
        <el-table-column label="响应内容 (Response Body)" min-width="200">
          <template #default="scope">
            <div class="code-block">{{ scope.row.response_body || '(空)' }}</div>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([])
const selectedRows = ref([])
const logDialogVisible = ref(false)
const logData = ref([])
const refundVisible = ref(false)
const currentTx = ref({})
const refundAmount = ref(0)
const refundReason = ref('')
const refunding = ref(false)

const filter = ref({
  out_trade_no: '',
  prepay_id: '',
  status: ''
})

const tradeTypeMap = {
  'WX:JSAPI': 'JSAPI',
  'WX:M_JSAPI': '小程序',
  'WX:APP': 'APP支付',
}

const loadData = async () => {
  try {
    const params = {
      out_trade_no: filter.value.out_trade_no,
      prepay_id: filter.value.prepay_id,
      status: filter.value.status
    }
    const res = await axios.get('/api/internal/transactions', { params })
    tableData.value = res.data
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const resetFilter = () => {
  filter.value = { out_trade_no: '', prepay_id: '', status: '' }
  loadData()
}

const handleSelectionChange = (val) => {
  selectedRows.value = val
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(
    `确定要永久删除选中的 ${selectedRows.value.length} 条交易记录吗？`,
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const ids = selectedRows.value.map(row => row.id)
      await axios.delete('/api/internal/transactions', { data: ids })
      ElMessage.success('删除成功')
      loadData()
    } catch (e) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const getStatusClass = (status) => {
  return status === 'SUCCESS' ? 'status-success' : 'status-warning'
}

const openPreview = (prepayId) => {
  window.open(`/pay/preview/${prepayId}`, '_blank', 'width=375,height=667')
}

const openRefund = (row) => {
  currentTx.value = row
  refundAmount.value = row.amount / 100
  refundReason.value = ''
  refundVisible.value = true
}

const doRefund = async () => {
  refunding.value = true
  try {
    await axios.post('/api/internal/simulate/refund', {
      transaction_id: currentTx.value.transaction_id,
      amount: Math.round(refundAmount.value * 100),
      reason: refundReason.value
    })
    ElMessage.success('退款提交成功')
    refundVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error('退款失败: ' + (e.response?.data?.error || e.message))
  } finally {
    refunding.value = false
  }
}

const viewDetails = async (row) => {
  logData.value = []
  logDialogVisible.value = true
  try {
    const res = await axios.get(`/api/internal/transactions/${row.transaction_id}/logs`)
    logData.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const retryCallback = async (row) => {
  try {
    await axios.post(`/api/internal/transactions/${row.id}/retry-callback`)
    ElMessage.success('已加入重试队列')
    setTimeout(() => {
      loadData()
    }, 1000)
  } catch (e) {
    ElMessage.error('重试失败')
  }
}

const formatJson = (jsonStr) => {
  try {
    const obj = JSON.parse(jsonStr)
    return JSON.stringify(obj, null, 2)
  } catch (e) {
    return jsonStr
  }
}

onMounted(loadData)
</script>

<style scoped>
.table-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter-bar {
  display: flex;
  gap: 12px;
}
.card-title {
  font-size: 20px;
  font-weight: 400;
  margin: 0;
  color: #202124;
}
/* Pill Status Styles */
.status-pill {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.5px;
}
.status-success {
  background-color: #e6f4ea;
  color: #137333;
}
.status-warning {
  background-color: #fef7e0;
  color: #ea8600;
}
.status-error {
  background-color: #fce8e6;
  color: #c5221f;
  cursor: help;
}
.status-gray {
  background-color: #f1f3f4;
  color: #5f6368;
}
.text-success {
  color: #137333;
  font-weight: 500;
}
.text-danger {
  color: #d93025;
  font-weight: 500;
}
.code-block {
  font-family: 'Roboto Mono', monospace;
  font-size: 12px;
  background-color: #f1f3f4;
  padding: 4px 8px;
  border-radius: 4px;
  word-break: break-all;
  max-height: 80px;
  overflow-y: auto;
}
.refund-tip {
  font-size: 12px;
  color: #5f6368;
  margin-top: 4px;
}
</style>
