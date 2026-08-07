<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as salesApi from '@/api/sales'
import * as settingsApi from '@/api/settings'
import type { Sale, SaleStatus, Settings } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppBadge from '@/components/ui/AppBadge.vue'
import ReceiptModal from '@/components/receipt/ReceiptModal.vue'
import { formatDateTime, formatMoney } from '@/lib/format'

const sales = ref<Sale[]>([])
const settings = ref<Settings | null>(null)
const loading = ref(true)
const error = ref('')

const search = ref('')
const from = ref('')
const to = ref('')

const viewingSale = ref<Sale | null>(null)
const showReceipt = ref(false)

function statusTone(status: SaleStatus): 'green' | 'amber' | 'red' {
  if (status === 'COMPLETED') return 'green'
  if (status === 'PARTIALLY_RETURNED') return 'amber'
  return 'red'
}

async function load() {
  loading.value = true
  try {
    const res = await salesApi.listSales({
      search: search.value || undefined,
      from: from.value || undefined,
      to: to.value || undefined,
    })
    sales.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load sales')
  } finally {
    loading.value = false
  }
}

async function viewSale(sale: Sale) {
  const res = await salesApi.getSale(sale.id)
  viewingSale.value = res.data.data
  showReceipt.value = true
}

onMounted(async () => {
  await Promise.all([load(), settingsApi.getSettings().then((r) => (settings.value = r.data.data))])
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-xl font-semibold text-slate-900">Sales History</h1>
      <p class="text-sm text-slate-500">Browse and reprint past transactions</p>
    </div>

    <div class="flex flex-wrap items-end gap-3">
      <AppInput v-model="search" placeholder="Search invoice number…" @update:model-value="load" class="w-56" />
      <AppInput v-model="from" label="From" type="date" @update:model-value="load" />
      <AppInput v-model="to" label="To" type="date" @update:model-value="load" />
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Invoice</th>
            <th class="px-5 py-3">Customer</th>
            <th class="px-5 py-3">Date</th>
            <th class="px-5 py-3 text-right">Total</th>
            <th class="px-5 py-3">Status</th>
            <th class="px-5 py-3">Cashier</th>
            <th class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="7">Loading…</td>
          </tr>
          <tr v-else-if="sales.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="7">No sales found.</td>
          </tr>
          <tr v-for="s in sales" :key="s.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ s.invoice_number }}</td>
            <td class="px-5 py-3 text-slate-600">{{ s.customer?.name ?? 'Walk-in Customer' }}</td>
            <td class="px-5 py-3 text-slate-500">{{ formatDateTime(s.created_at) }}</td>
            <td class="px-5 py-3 text-right text-slate-900">{{ formatMoney(s.total) }}</td>
            <td class="px-5 py-3"><AppBadge :tone="statusTone(s.status)">{{ s.status.replace('_', ' ') }}</AppBadge></td>
            <td class="px-5 py-3 text-slate-500">{{ s.created_by_user?.name ?? '—' }}</td>
            <td class="px-5 py-3 text-right">
              <button class="text-brand-600 hover:underline" @click="viewSale(s)">View / Print</button>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <ReceiptModal v-if="viewingSale && settings" v-model="showReceipt" :sale="viewingSale" :settings="settings" />
  </div>
</template>
