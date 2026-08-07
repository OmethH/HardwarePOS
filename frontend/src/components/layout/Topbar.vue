<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-6">
    <div />
    <div class="flex items-center gap-4">
      <div class="text-right">
        <p class="text-sm font-medium text-slate-900">{{ auth.user?.name }}</p>
        <p class="text-xs capitalize text-slate-500">{{ auth.role }}</p>
      </div>
      <div
        class="flex h-9 w-9 items-center justify-center rounded-full bg-brand-100 text-sm font-semibold text-brand-700"
      >
        {{ auth.user?.name?.charAt(0).toUpperCase() }}
      </div>
      <button
        class="rounded-lg border border-slate-300 px-3 py-1.5 text-sm font-medium text-slate-600 hover:bg-slate-50"
        @click="logout"
      >
        Log out
      </button>
    </div>
  </header>
</template>
