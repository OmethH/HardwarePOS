<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import * as reportsApi from '@/api/reports'
import type { InventoryReportRow, ProductPerformanceRow, SalesReportRow } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import { formatMoney, formatQty } from '@/lib/format'

type Tab = 'sales' | 'inventory' | 'performance'

const tabs: { value: Tab; label: string }[] = [
  { value: 'sales', label: 'Sales Report' },
  { value: 'inventory', label: 'Inventory Report' },
  { value: 'performance', label: 'Product Performance' },
]

const rangePresets = [
  { value: 'today', label: 'Today' },
  { value: 'week', label: 'This Week' },
  { value: 'month', label: 'This Month' },
  { value: 'custom', label: 'Custom Range' },
]

const activeTab = ref<Tab>('sales')
const rangePreset = ref('month')
const from = ref('')
const to = ref('')

const loading = ref(false)
const error = ref('')

const salesRows = ref<SalesReportRow[]>([])
const inventoryRows = ref<InventoryReportRow[]>([])
const performanceRows = ref<ProductPerformanceRow[]>([])

function applyPreset() {
  const now = new Date()
  const iso = (d: Date) => d.toISOString().slice(0, 10)

  if (rangePreset.value === 'today') {
    from.value = iso(now)
    to.value = iso(now)
  } else if (rangePreset.value === 'week') {
    const start = new Date(now)
    start.setDate(now.getDate() - now.getDay())
    from.value = iso(start)
    to.value = iso(now)
  } else if (rangePreset.value === 'month') {
    const start = new Date(now.getFullYear(), now.getMonth(), 1)
    from.value = iso(start)
    to.value = iso(now)
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    if (activeTab.value === 'sales') {
      const res = await reportsApi.getSalesReport(from.value, to.value)
      salesRows.value = res.data.data
    } else if (activeTab.value === 'inventory') {
      const res = await reportsApi.getInventoryReport()
      inventoryRows.value = res.data.data
    } else {
      const res = await reportsApi.getProductPerformance(from.value, to.value)
      performanceRows.value = res.data.data
    }
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load report')
  } finally {
    loading.value = false
  }
}

async function exportReport(format: 'csv' | 'pdf') {
  const report = activeTab.value === 'sales' ? 'sales' : activeTab.value === 'inventory' ? 'inventory' : 'product-performance'
  try {
    await reportsApi.downloadReport(report, format, { from: from.value, to: to.value })
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to export report')
  }
}

watch(rangePreset, () => {
  if (rangePreset.value !== 'custom') {
    applyPreset()
    load()
  }
})

watch(activeTab, load)

onMounted(() => {
  applyPreset()
  load()
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-xl font-semibold text-slate-900">Reports</h1>
      <p class="text-sm text-slate-500">Sales, inventory and product performance</p>
    </div>

    <div class="flex gap-2 border-b border-slate-200">
      <button
        v-for="t in tabs"
        :key="t.value"
        class="border-b-2 px-4 py-2 text-sm font-medium"
        :class="activeTab === t.value ? 'border-brand-600 text-brand-600' : 'border-transparent text-slate-500 hover:text-slate-700'"
        @click="activeTab = t.value"
      >
        {{ t.label }}
      </button>
    </div>

    <div class="flex flex-wrap items-end gap-3">
      <AppSelect
        v-if="activeTab !== 'inventory'"
        v-model="rangePreset"
        label="Range"
        :options="rangePresets"
        class="w-48"
      />
      <template v-if="activeTab !== 'inventory' && rangePreset === 'custom'">
        <AppInput v-model="from" label="From" type="date" @update:model-value="load" />
        <AppInput v-model="to" label="To" type="date" @update:model-value="load" />
      </template>
      <div class="ml-auto flex gap-2">
        <AppButton variant="secondary" @click="exportReport('csv')">Export CSV</AppButton>
        <AppButton variant="secondary" @click="exportReport('pdf')">Export PDF</AppButton>
      </div>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table v-if="activeTab === 'sales'" class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Date</th>
            <th class="px-5 py-3 text-right">Invoices</th>
            <th class="px-5 py-3 text-right">Revenue</th>
            <th class="px-5 py-3 text-right">Discount</th>
            <th class="px-5 py-3 text-right">Profit</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="5">Loading…</td>
          </tr>
          <tr v-else-if="salesRows.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="5">No sales in this range.</td>
          </tr>
          <tr v-for="row in salesRows" :key="row.date" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ row.date }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ row.invoices }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(row.revenue) }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(row.discount) }}</td>
            <td class="px-5 py-3 text-right font-medium text-emerald-600">{{ formatMoney(row.profit) }}</td>
          </tr>
        </tbody>
      </table>

      <table v-else-if="activeTab === 'inventory'" class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Product</th>
            <th class="px-5 py-3">SKU</th>
            <th class="px-5 py-3 text-right">Current Stock</th>
            <th class="px-5 py-3 text-right">Stock Value</th>
            <th class="px-5 py-3">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="5">Loading…</td>
          </tr>
          <tr v-else-if="inventoryRows.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="5">No products found.</td>
          </tr>
          <tr v-for="row in inventoryRows" :key="row.sku" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ row.product }}</td>
            <td class="px-5 py-3 text-slate-600">{{ row.sku }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatQty(row.current_stock) }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(row.stock_value) }}</td>
            <td class="px-5 py-3 text-slate-600">{{ row.status }}</td>
          </tr>
        </tbody>
      </table>

      <table v-else class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Product</th>
            <th class="px-5 py-3 text-right">Quantity Sold</th>
            <th class="px-5 py-3 text-right">Revenue</th>
            <th class="px-5 py-3 text-right">Profit</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="4">Loading…</td>
          </tr>
          <tr v-else-if="performanceRows.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="4">No sales in this range.</td>
          </tr>
          <tr v-for="row in performanceRows" :key="row.product" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ row.product }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatQty(row.quantity_sold) }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(row.revenue) }}</td>
            <td class="px-5 py-3 text-right font-medium text-emerald-600">{{ formatMoney(row.profit) }}</td>
          </tr>
        </tbody>
      </table>
    </AppCard>
  </div>
</template>
