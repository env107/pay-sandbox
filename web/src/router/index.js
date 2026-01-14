import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '../views/admin/Layout.vue'
import MerchantList from '../views/admin/MerchantList.vue'
import TransactionList from '../views/admin/TransactionList.vue'
import RefundList from '../views/admin/RefundList.vue'
import PayPreview from '../views/mobile/PayPreview.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        { path: 'merchants', component: MerchantList },
        { path: 'transactions', component: TransactionList },
        { path: 'refunds', component: RefundList },
        { path: '', redirect: '/admin/merchants' }
      ]
    },
    {
      path: '/pay/preview/:prepay_id',
      component: PayPreview
    },
    {
      path: '/',
      redirect: '/admin'
    }
  ]
})

export default router
