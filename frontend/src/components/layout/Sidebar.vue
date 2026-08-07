<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

interface NavItem {
  to: string
  label: string
  icon: string
  roles?: string[]
}

const navItems: NavItem[] = [
  { to: '/dashboard', label: 'Dashboard', icon: '📊', roles: ['admin', 'manager'] },
  { to: '/pos', label: 'POS', icon: '🛒' },
  { to: '/sales', label: 'Sales History', icon: '🧾' },
  { to: '/returns', label: 'Returns', icon: '↩️', roles: ['admin', 'manager'] },
  { to: '/products', label: 'Products', icon: '📦' },
  { to: '/categories', label: 'Categories', icon: '🏷️' },
  { to: '/inventory', label: 'Inventory', icon: '📋' },
  { to: '/purchases', label: 'Purchases', icon: '🚚', roles: ['admin', 'manager'] },
  { to: '/customers', label: 'Customers', icon: '👥' },
  { to: '/suppliers', label: 'Suppliers', icon: '🏭' },
  { to: '/reports', label: 'Reports', icon: '📈', roles: ['admin', 'manager'] },
  { to: '/users', label: 'Users', icon: '🔐', roles: ['admin'] },
  { to: '/settings', label: 'Settings', icon: '⚙️' },
]

const visibleItems = computed(() => navItems.filter((item) => !item.roles || auth.hasRole(...item.roles)))
</script>

<template>
  <aside class="flex h-full w-64 flex-col border-r border-slate-200 bg-white">
    <div class="flex items-center gap-2 px-5 py-5">
      <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-600 text-lg font-bold text-white">
        H
      </div>
      <div>
        <p class="text-sm font-semibold text-slate-900">HardwarePOS</p>
        <p class="text-xs text-slate-500">Inventory & Sales</p>
      </div>
    </div>

    <nav class="flex-1 space-y-1 overflow-y-auto px-3 pb-4">
      <RouterLink
        v-for="item in visibleItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100"
        active-class="!bg-brand-50 !text-brand-700"
      >
        <span class="text-base leading-none">{{ item.icon }}</span>
        {{ item.label }}
      </RouterLink>
    </nav>
  </aside>
</template>
