<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as purchasesApi from '@/api/purchases'
import * as suppliersApi from '@/api/suppliers'
import * as productsApi from '@/api/products'
import type { Product, Purchase, Supplier } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppModal from '@/components/ui/AppModal.vue'
import { formatDateTime, formatMoney } from '@/lib/format'

const purchases = ref<Purchase[]>([])
const suppliers = ref<Supplier[]>([])
const products = ref<Product[]>([])
const loading = ref(true)
const error = ref('')

const showModal = ref(false)
const saving = ref(false)
const formError = ref('')
const supplierId = ref('')
const lines = ref<{ product_id: string; quantity: string; price: string }[]>([])

const supplierOptions = computed(() => suppliers.value.map((s) => ({ value: s.id, label: s.name })))
const productOptions = computed(() => products.value.map((p) => ({ value: p.id, label: `${p.name} (${p.sku})` })))

const total = computed(() =>
  lines.value.reduce((sum, l) => sum + (Number(l.quantity) || 0) * (Number(l.price) || 0), 0),
)

async function load() {
  loading.value = true
  try {
    const res = await purchasesApi.listPurchases(50)
    purchases.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load purchases')
  } finally {
    loading.value = false
  }
}

async function loadRefs() {
  const [sRes, pRes] = await Promise.all([suppliersApi.listSuppliers(), productsApi.listProducts()])
  suppliers.value = sRes.data.data
  products.value = pRes.data.data
}

function openCreate() {
  supplierId.value = ''
  lines.value = [{ product_id: '', quantity: '', price: '' }]
  formError.value = ''
  showModal.value = true
}

function addLine() {
  lines.value.push({ product_id: '', quantity: '', price: '' })
}

function removeLine(index: number) {
  lines.value.splice(index, 1)
}

function onProductChosen(index: number) {
  const line = lines.value[index]
  const product = products.value.find((p) => p.id === Number(line.product_id))
  if (product && !line.price) {
    line.price = String(product.purchase_price)
  }
}

async function save() {
  formError.value = ''
  const items = lines.value
    .filter((l) => l.product_id && Number(l.quantity) > 0 && Number(l.price) >= 0)
    .map((l) => ({ product_id: Number(l.product_id), quantity: Number(l.quantity), price: Number(l.price) }))

  if (!supplierId.value || items.length === 0) {
    formError.value = 'Select a supplier and add at least one valid item.'
    return
  }

  saving.value = true
  try {
    await purchasesApi.createPurchase({ supplier_id: Number(supplierId.value), items })
    showModal.value = false
    await load()
  } catch (err) {
    formError.value = apiErrorMessage(err, 'Failed to create purchase')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await Promise.all([load(), loadRefs()])
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Purchases</h1>
        <p class="text-sm text-slate-500">Record incoming stock from suppliers</p>
      </div>
      <AppButton @click="openCreate">+ New Purchase</AppButton>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Reference</th>
            <th class="px-5 py-3">Supplier</th>
            <th class="px-5 py-3">Items</th>
            <th class="px-5 py-3 text-right">Total</th>
            <th class="px-5 py-3">Date</th>
            <th class="px-5 py-3">By</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="6">Loading…</td>
          </tr>
          <tr v-else-if="purchases.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="6">No purchases yet.</td>
          </tr>
          <tr v-for="p in purchases" :key="p.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ p.reference_number }}</td>
            <td class="px-5 py-3 text-slate-600">{{ p.supplier?.name }}</td>
            <td class="px-5 py-3 text-slate-500">{{ p.items.length }} item(s)</td>
            <td class="px-5 py-3 text-right text-slate-900">{{ formatMoney(p.total_amount) }}</td>
            <td class="px-5 py-3 text-slate-500">{{ formatDateTime(p.created_at) }}</td>
            <td class="px-5 py-3 text-slate-500">{{ p.created_by_user?.name ?? '—' }}</td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" title="New Purchase" wide>
      <form class="space-y-4" @submit.prevent="save">
        <AppSelect v-model="supplierId" label="Supplier" :options="supplierOptions" placeholder="Select supplier" required />

        <div class="space-y-3">
          <p class="text-sm font-medium text-slate-700">Items</p>
          <div v-for="(line, i) in lines" :key="i" class="grid grid-cols-12 items-end gap-2">
            <div class="col-span-6">
              <AppSelect
                v-model="line.product_id"
                :options="productOptions"
                placeholder="Product"
                @update:model-value="onProductChosen(i)"
              />
            </div>
            <div class="col-span-2">
              <AppInput v-model="line.quantity" type="number" step="0.01" min="0" placeholder="Qty" />
            </div>
            <div class="col-span-3">
              <AppInput v-model="line.price" type="number" step="0.01" min="0" placeholder="Price" />
            </div>
            <div class="col-span-1">
              <button type="button" class="text-red-600 hover:underline" @click="removeLine(i)">✕</button>
            </div>
          </div>
          <AppButton type="button" variant="secondary" @click="addLine">+ Add item</AppButton>
        </div>

        <div class="flex items-center justify-between border-t border-slate-100 pt-3 text-sm font-semibold text-slate-900">
          <span>Total</span>
          <span>{{ formatMoney(total) }}</span>
        </div>

        <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>

        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Confirm Purchase' }}</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
