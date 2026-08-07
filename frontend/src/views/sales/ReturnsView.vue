<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as salesApi from '@/api/sales'
import * as returnsApi from '@/api/returns'
import type { ReturnRecord, Sale } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import { formatDateTime, formatMoney, formatQty } from '@/lib/format'

const invoiceSearch = ref('')
const sale = ref<Sale | null>(null)
const searchError = ref('')
const searching = ref(false)

const returnQuantities = ref<Record<number, string>>({})
const submitting = ref(false)
const submitError = ref('')
const submitSuccess = ref('')

const history = ref<ReturnRecord[]>([])
const loadingHistory = ref(true)

const returnableItems = computed(() =>
  (sale.value?.items ?? []).map((item) => ({
    ...item,
    remaining: item.quantity - item.returned_quantity,
  })),
)

const refundTotal = computed(() =>
  returnableItems.value.reduce((sum, item) => {
    const qty = Number(returnQuantities.value[item.id] || 0)
    return sum + qty * item.selling_price
  }, 0),
)

async function findSale() {
  searchError.value = ''
  sale.value = null
  submitSuccess.value = ''
  if (!invoiceSearch.value.trim()) return
  searching.value = true
  try {
    const res = await salesApi.getSaleByInvoice(invoiceSearch.value.trim())
    sale.value = res.data.data
    returnQuantities.value = {}
  } catch (err) {
    searchError.value = apiErrorMessage(err, 'Invoice not found')
  } finally {
    searching.value = false
  }
}

async function submitReturn() {
  if (!sale.value) return
  submitError.value = ''
  const items = returnableItems.value
    .filter((item) => Number(returnQuantities.value[item.id] || 0) > 0)
    .map((item) => ({ sale_item_id: item.id, quantity: Number(returnQuantities.value[item.id]) }))

  if (items.length === 0) {
    submitError.value = 'Enter a return quantity for at least one item.'
    return
  }

  submitting.value = true
  try {
    await returnsApi.createReturn({ sale_id: sale.value.id, items })
    submitSuccess.value = `Return processed. Refund: ${formatMoney(refundTotal.value)}`
    sale.value = null
    invoiceSearch.value = ''
    await loadHistory()
  } catch (err) {
    submitError.value = apiErrorMessage(err, 'Failed to process return')
  } finally {
    submitting.value = false
  }
}

async function loadHistory() {
  loadingHistory.value = true
  try {
    const res = await returnsApi.listReturns(50)
    history.value = res.data.data
  } finally {
    loadingHistory.value = false
  }
}

onMounted(loadHistory)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-xl font-semibold text-slate-900">Returns</h1>
      <p class="text-sm text-slate-500">Look up an invoice and process item returns</p>
    </div>

    <AppCard title="Find Invoice">
      <div class="flex items-end gap-3">
        <AppInput v-model="invoiceSearch" label="Invoice number" placeholder="INV00001" class="max-w-xs" />
        <AppButton :disabled="searching" @click="findSale">{{ searching ? 'Searching…' : 'Find' }}</AppButton>
      </div>
      <p v-if="searchError" class="mt-2 text-sm text-red-600">{{ searchError }}</p>
      <p v-if="submitSuccess" class="mt-2 text-sm text-emerald-600">{{ submitSuccess }}</p>
    </AppCard>

    <AppCard v-if="sale" :title="`Invoice ${sale.invoice_number}`">
      <p class="mb-4 text-sm text-slate-500">
        {{ sale.customer?.name ?? 'Walk-in Customer' }} · {{ formatDateTime(sale.created_at) }}
      </p>

      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="py-2">Product</th>
            <th class="py-2 text-right">Sold</th>
            <th class="py-2 text-right">Already Returned</th>
            <th class="py-2 text-right">Remaining</th>
            <th class="py-2 text-right">Return Qty</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-for="item in returnableItems" :key="item.id">
            <td class="py-2 font-medium text-slate-900">{{ item.product?.name }}</td>
            <td class="py-2 text-right text-slate-600">{{ formatQty(item.quantity) }}</td>
            <td class="py-2 text-right text-slate-600">{{ formatQty(item.returned_quantity) }}</td>
            <td class="py-2 text-right text-slate-600">{{ formatQty(item.remaining) }}</td>
            <td class="py-2 text-right">
              <input
                type="number"
                min="0"
                :max="item.remaining"
                :disabled="item.remaining <= 0"
                v-model="returnQuantities[item.id]"
                class="w-20 rounded-md border border-slate-300 px-2 py-1 text-right text-sm disabled:bg-slate-100"
              />
            </td>
          </tr>
        </tbody>
      </table>

      <div class="mt-4 flex items-center justify-between border-t border-slate-100 pt-4">
        <span class="text-sm font-medium text-slate-700">Refund total: {{ formatMoney(refundTotal) }}</span>
        <AppButton :disabled="submitting" @click="submitReturn">
          {{ submitting ? 'Processing…' : 'Process Return' }}
        </AppButton>
      </div>
      <p v-if="submitError" class="mt-2 text-sm text-red-600">{{ submitError }}</p>
    </AppCard>

    <AppCard title="Recent Returns" :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Return #</th>
            <th class="px-5 py-3">Invoice</th>
            <th class="px-5 py-3 text-right">Refund</th>
            <th class="px-5 py-3">Date</th>
            <th class="px-5 py-3">By</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loadingHistory">
            <td class="px-5 py-4 text-slate-400" colspan="5">Loading…</td>
          </tr>
          <tr v-else-if="history.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="5">No returns yet.</td>
          </tr>
          <tr v-for="r in history" :key="r.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ r.return_number }}</td>
            <td class="px-5 py-3 text-slate-600">{{ r.sale?.invoice_number }}</td>
            <td class="px-5 py-3 text-right text-slate-900">{{ formatMoney(r.total_refund) }}</td>
            <td class="px-5 py-3 text-slate-500">{{ formatDateTime(r.created_at) }}</td>
            <td class="px-5 py-3 text-slate-500">{{ r.created_by_user?.name ?? '—' }}</td>
          </tr>
        </tbody>
      </table>
    </AppCard>
  </div>
</template>
