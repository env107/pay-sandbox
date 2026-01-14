<template>
  <div class="material-card">
    <div class="table-header">
      <h3 class="card-title">退款流水记录</h3>
      <div class="filter-bar">
        <el-input v-model="filter.out_refund_no" placeholder="商户退款单号" style="width: 180px" clearable />
        <el-input v-model="filter.transaction_id" placeholder="支付订单号" style="width: 180px" clearable />
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetFilter">重置</el-button>
        <el-button 
          type="danger" 
          @click="handleBatchDelete" 
          :disabled="multipleSelection.length === 0">
          批量删除
        </el-button>
      </div>
    </div>

    <el-table 
      :data="tableData" 
      style="width: 100%" 
      size="large" 
      :header-cell-style="{ background: '#f8f9fa', color: '#5f6368', fontWeight: 500 }"
      @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="created_at" label="退款时间" min-width="180">
        <template #default="scope">
          {{ new Date(scope.row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column prop="out_refund_no" label="商户退款单号" min-width="200" />
      <el-table-column prop="transaction_id" label="支付订单号" min-width="220" />
      <el-table-column prop="amount" label="退款金额 (分)" min-width="120" align="right" />
      <el-table-column prop="reason" label="原因" min-width="150" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" min-width="120" align="center">
        <template #default="scope">
          <span class="status-pill status-success">{{ scope.row.status }}</span>
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
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="scope">
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

    <!-- 回调日志弹窗 -->
    <el-dialog v-model="logDialogVisible" title="退款回调日志" width="900px" top="5vh">
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
const logDialogVisible = ref(false)
const logData = ref([])
const multipleSelection = ref([])
const filter = ref({
  out_refund_no: '',
  transaction_id: ''
})

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const handleBatchDelete = () => {
  if (multipleSelection.value.length === 0) return

  ElMessageBox.confirm(
    `确定要永久删除选中的 ${multipleSelection.value.length} 条记录吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const ids = multipleSelection.value.map(item => item.id)
      await axios.delete('/api/internal/refunds', { data: ids })
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const loadData = async () => {
  try {
    const params = {
      out_refund_no: filter.value.out_refund_no,
      transaction_id: filter.value.transaction_id
    }
    const res = await axios.get('/api/internal/refunds', { params })
    tableData.value = res.data
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const resetFilter = () => {
  filter.value = { out_refund_no: '', transaction_id: '' }
  loadData()
}

const viewDetails = async (row) => {
  logData.value = []
  logDialogVisible.value = true
  try {
    // 使用 transaction_id 字段查询日志，因为在 TriggerRefundCallback 中我们将 RefundID 存入了 TransactionID 字段
    const res = await axios.get(`/api/internal/refunds/${row.refund_id}/logs`)
    logData.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const retryCallback = async (row) => {
  try {
    await axios.post(`/api/internal/refunds/${row.refund_id}/retry-callback`)
    ElMessage.success('已加入重试队列')
    loadData()
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
</style>