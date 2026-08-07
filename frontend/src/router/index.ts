import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppShell.vue'),
    children: [
      { path: '', redirect: '/dashboard' },
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { roles: ['admin', 'manager'] },
      },
      {
        path: 'products',
        name: 'products',
        component: () => import('@/views/products/ProductsView.vue'),
      },
      {
        path: 'categories',
        name: 'categories',
        component: () => import('@/views/categories/CategoriesView.vue'),
      },
      {
        path: 'inventory',
        name: 'inventory',
        component: () => import('@/views/inventory/InventoryView.vue'),
      },
      {
        path: 'purchases',
        name: 'purchases',
        component: () => import('@/views/purchases/PurchasesView.vue'),
        meta: { roles: ['admin', 'manager'] },
      },
      {
        path: 'pos',
        name: 'pos',
        component: () => import('@/views/pos/PosView.vue'),
      },
      {
        path: 'sales',
        name: 'sales',
        component: () => import('@/views/sales/SalesView.vue'),
      },
      {
        path: 'returns',
        name: 'returns',
        component: () => import('@/views/sales/ReturnsView.vue'),
        meta: { roles: ['admin', 'manager'] },
      },
      {
        path: 'customers',
        name: 'customers',
        component: () => import('@/views/customers/CustomersView.vue'),
      },
      {
        path: 'suppliers',
        name: 'suppliers',
        component: () => import('@/views/suppliers/SuppliersView.vue'),
      },
      {
        path: 'reports',
        name: 'reports',
        component: () => import('@/views/reports/ReportsView.vue'),
        meta: { roles: ['admin', 'manager'] },
      },
      {
        path: 'users',
        name: 'users',
        component: () => import('@/views/users/UsersView.vue'),
        meta: { roles: ['admin'] },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/settings/SettingsView.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.public) {
    if (to.name === 'login' && auth.isAuthenticated) return { path: '/' }
    return true
  }

  if (!auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  // Cashiers have no dashboard access — send them straight to the POS screen.
  if (to.path === '/dashboard' && !auth.hasRole('admin', 'manager')) {
    return { path: '/pos' }
  }

  const roles = to.meta.roles as string[] | undefined
  if (roles && !auth.hasRole(...roles)) {
    return { path: '/pos' }
  }

  return true
})

export default router
