<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as settingsApi from '@/api/settings'
import { apiErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'

const auth = useAuthStore()
const canEdit = auth.hasRole('admin')

const form = ref({
  shop_name: '',
  address: '',
  phone: '',
  receipt_footer: '',
  tax_percentage: 0,
})

const loading = ref(true)
const saving = ref(false)
const error = ref('')
const saved = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await settingsApi.getSettings()
    const s = res.data.data
    form.value = {
      shop_name: s.shop_name,
      address: s.address,
      phone: s.phone,
      receipt_footer: s.receipt_footer,
      tax_percentage: s.tax_percentage,
    }
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load settings')
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  error.value = ''
  saved.value = false
  try {
    await settingsApi.updateSettings({
      ...form.value,
      tax_percentage: Number(form.value.tax_percentage),
    })
    saved.value = true
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to save settings')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="max-w-2xl space-y-6">
    <div>
      <h1 class="text-xl font-semibold text-slate-900">Settings</h1>
      <p class="text-sm text-slate-500">Shop details used on receipts and reports</p>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    <p v-if="saved" class="text-sm text-emerald-600">Settings saved.</p>

    <AppCard>
      <div v-if="loading" class="text-sm text-slate-400">Loading…</div>
      <form v-else class="space-y-4" @submit.prevent="save">
        <AppInput v-model="form.shop_name" label="Shop Name" required :disabled="!canEdit" />
        <AppInput v-model="form.address" label="Address" :disabled="!canEdit" />
        <AppInput v-model="form.phone" label="Phone Number" :disabled="!canEdit" />
        <AppInput v-model="form.receipt_footer" label="Receipt Footer" :disabled="!canEdit" />
        <AppInput
          v-model="form.tax_percentage"
          label="Tax Percentage"
          type="number"
          step="0.01"
          min="0"
          :disabled="!canEdit"
        />

        <div v-if="canEdit" class="flex justify-end">
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save Settings' }}</AppButton>
        </div>
        <p v-else class="text-xs text-slate-400">Only admins can edit shop settings.</p>
      </form>
    </AppCard>
  </div>
</template>
