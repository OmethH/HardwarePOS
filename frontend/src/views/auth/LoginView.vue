<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err) {
    error.value = apiErrorMessage(err, 'Invalid username or password')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-100 px-4">
    <div class="w-full max-w-sm rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
      <div class="mb-6 flex flex-col items-center text-center">
        <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-brand-600 text-xl font-bold text-white">
          H
        </div>
        <h1 class="text-lg font-semibold text-slate-900">HardwarePOS</h1>
        <p class="text-sm text-slate-500">Sign in to manage your shop</p>
      </div>

      <form class="space-y-4" @submit.prevent="submit">
        <AppInput v-model="username" label="Username" required placeholder="admin" />
        <AppInput v-model="password" label="Password" type="password" required placeholder="••••••••" />

        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

        <AppButton type="submit" class="w-full justify-center" :disabled="loading">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </AppButton>
      </form>
    </div>
  </div>
</template>
