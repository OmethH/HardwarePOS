<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as inventoryApi from '@/api/inventory'
import * as productsApi from '@/api/products'
import type { Product, StockMovement } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppBadge from '@/components/ui/AppBadge.vue'
import { formatDateTime, formatQty } from '@/lib/format'

const auth = useAuthStore()
const canAdjust = auth.hasRole('admin', 'manager')

const movements = ref<StockMovement[]>([])
const products = ref<Product[]>([])
const loading = ref(true)
const error = ref('')

const showModal = ref(false)
const saving = ref(false)
const formError = ref('')
const form = ref({ product_id: '', quantity: '', note: '' })

const productOptions = computed(() => products.value.map((p) => ({ value: p.id, label: `${p.name} (${p.sku})` })))

function typeTone(type: StockMovement['type']): 'green' | 'red' | 'blue' | 'amber' {
  switch (type) {
    case 'PURCHASE':
      return 'green'
    case 'SALE':
      return 'blue'
    case 'RETURN':
      return 'amber'
    default:
      return 'red'
  }
}

async function load() {
  loading.value = true
  try {
    const res = await inventoryApi.listMovements(100)
    movements.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load stock movements')
  } finally {
    loading.value = false
  }
}

async function loadProducts() {
  const res = await productsApi.listProducts()
  products.value = res.data.data
}

function openAdjust() {
  form.value = { product_id: '', quantity: '', note: '' }
  formError.value = ''
  showModal.value = true
}

async function saveAdjustment() {
  saving.value = true
  formError.value = ''
  try {
    await inventoryApi.createAdjustment({
      product_id: Number(form.value.product_id),
      quantity: Number(form.value.quantity),
      note: form.value.note,
    })
    showModal.value = false
    await load()
  } catch (err) {
    formError.value = apiErrorMessage(err, 'Failed to record adjustment')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await Promise.all([load(), loadProducts()])
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Inventory Movements</h1>
        <p class="text-sm text-slate-500">Every stock change, with its source</p>
      </div>
      <AppButton v-if="canAdjust" @click="openAdjust">+ Manual Adjustment</AppButton>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Date</th>
            <th class="px-5 py-3">Product</th>
            <th class="px-5 py-3">Type</th>
            <th class="px-5 py-3 text-right">Quantity</th>
            <th class="px-5 py-3">Note</th>
            <th class="px-5 py-3">By</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="6">Loading…</td>
          </tr>
          <tr v-else-if="movements.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="6">No stock movements yet.</td>
          </tr>
          <tr v-for="m in movements" :key="m.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 text-slate-500">{{ formatDateTime(m.created_at) }}</td>
            <td class="px-5 py-3 font-medium text-slate-900">{{ m.product?.name ?? `#${m.product_id}` }}</td>
            <td class="px-5 py-3"><AppBadge :tone="typeTone(m.type)">{{ m.type }}</AppBadge></td>
            <td class="px-5 py-3 text-right" :class="m.quantity < 0 ? 'text-red-600' : 'text-emerald-600'">
              {{ m.quantity > 0 ? '+' : '' }}{{ formatQty(m.quantity) }}
            </td>
            <td class="px-5 py-3 text-slate-500">{{ m.note }}</td>
            <td class="px-5 py-3 text-slate-500">{{ m.created_by_user?.name ?? '—' }}</td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" title="Manual Stock Adjustment">
      <form class="space-y-4" @submit.prevent="saveAdjustment">
        <AppSelect v-model="form.product_id" label="Product" :options="productOptions" placeholder="Select product" required />
        <AppInput v-model="form.quantity" label="Quantity change (use negative to reduce)" type="number" step="0.01" required />
        <AppInput v-model="form.note" label="Reason" required placeholder="e.g. Stocktake correction" />

        <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>

        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
