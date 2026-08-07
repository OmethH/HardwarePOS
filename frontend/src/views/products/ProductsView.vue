<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as productsApi from '@/api/products'
import * as categoriesApi from '@/api/categories'
import type { Category, Product } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppBadge from '@/components/ui/AppBadge.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'
import { formatMoney, formatQty } from '@/lib/format'

const auth = useAuthStore()
const canManage = auth.hasRole('admin', 'manager')

const UNIT_TYPES = ['Piece', 'Box', 'Meter', 'Kg', 'Liter', 'Pack']

const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const loading = ref(true)
const error = ref('')

const search = ref('')
const categoryFilter = ref('')
const lowStockOnly = ref(false)

const showModal = ref(false)
const editing = ref<Product | null>(null)
const saving = ref(false)
const formError = ref('')

const emptyForm = (): productsApi.ProductInput => ({
  sku: '',
  barcode: null,
  name: '',
  category_id: null,
  brand: '',
  description: '',
  unit_type: 'Piece',
  purchase_price: 0,
  selling_price: 0,
  minimum_stock: 0,
  status: 'active',
})
const form = ref<productsApi.ProductInput>(emptyForm())

const showConfirm = ref(false)
const pendingDelete = ref<Product | null>(null)

const categoryOptions = computed(() => categories.value.map((c) => ({ value: c.id, label: c.name })))
const unitOptions = computed(() => UNIT_TYPES.map((u) => ({ value: u, label: u })))

function stockStatus(p: Product): { tone: 'green' | 'amber' | 'red'; label: string } {
  if (p.current_stock <= 0) return { tone: 'red', label: 'Out of stock' }
  if (p.current_stock <= p.minimum_stock) return { tone: 'amber', label: 'Low stock' }
  return { tone: 'green', label: 'In stock' }
}

async function load() {
  loading.value = true
  try {
    const res = await productsApi.listProducts({
      search: search.value || undefined,
      category_id: categoryFilter.value ? Number(categoryFilter.value) : undefined,
      low_stock: lowStockOnly.value || undefined,
    })
    products.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load products')
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  const res = await categoriesApi.listCategories()
  categories.value = res.data.data
}

function openCreate() {
  editing.value = null
  form.value = emptyForm()
  formError.value = ''
  showModal.value = true
}

function openEdit(product: Product) {
  editing.value = product
  form.value = {
    sku: product.sku,
    barcode: product.barcode,
    name: product.name,
    category_id: product.category_id,
    brand: product.brand,
    description: product.description,
    unit_type: product.unit_type,
    purchase_price: product.purchase_price,
    selling_price: product.selling_price,
    minimum_stock: product.minimum_stock,
    status: product.status,
  }
  formError.value = ''
  showModal.value = true
}

async function save() {
  saving.value = true
  formError.value = ''
  try {
    const payload = {
      ...form.value,
      category_id: form.value.category_id ? Number(form.value.category_id) : null,
      purchase_price: Number(form.value.purchase_price),
      selling_price: Number(form.value.selling_price),
      minimum_stock: Number(form.value.minimum_stock),
    }
    if (editing.value) {
      await productsApi.updateProduct(editing.value.id, payload)
    } else {
      await productsApi.createProduct(payload)
    }
    showModal.value = false
    await load()
  } catch (err) {
    formError.value = apiErrorMessage(err, 'Failed to save product')
  } finally {
    saving.value = false
  }
}

function askDelete(product: Product) {
  pendingDelete.value = product
  showConfirm.value = true
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await productsApi.deleteProduct(pendingDelete.value.id)
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to delete product')
  }
}

onMounted(async () => {
  await Promise.all([load(), loadCategories()])
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Products</h1>
        <p class="text-sm text-slate-500">Manage your product catalog</p>
      </div>
      <AppButton v-if="canManage" @click="openCreate">+ New Product</AppButton>
    </div>

    <div class="flex flex-wrap items-end gap-3">
      <AppInput v-model="search" placeholder="Search name, SKU, barcode…" @update:model-value="load" class="w-64" />
      <AppSelect v-model="categoryFilter" :options="categoryOptions" placeholder="All categories" @update:model-value="load" class="w-48" />
      <label class="flex items-center gap-2 text-sm text-slate-600">
        <input type="checkbox" v-model="lowStockOnly" @change="load" class="rounded border-slate-300" />
        Low stock only
      </label>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Product</th>
            <th class="px-5 py-3">SKU</th>
            <th class="px-5 py-3">Category</th>
            <th class="px-5 py-3 text-right">Purchase</th>
            <th class="px-5 py-3 text-right">Selling</th>
            <th class="px-5 py-3 text-right">Stock</th>
            <th class="px-5 py-3">Status</th>
            <th v-if="canManage" class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="8">Loading…</td>
          </tr>
          <tr v-else-if="products.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="8">No products found.</td>
          </tr>
          <tr v-for="p in products" :key="p.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ p.name }}</td>
            <td class="px-5 py-3 text-slate-500">{{ p.sku }}</td>
            <td class="px-5 py-3 text-slate-600">{{ p.category?.name ?? '—' }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(p.purchase_price) }}</td>
            <td class="px-5 py-3 text-right text-slate-600">{{ formatMoney(p.selling_price) }}</td>
            <td class="px-5 py-3 text-right text-slate-900">{{ formatQty(p.current_stock) }} {{ p.unit_type }}</td>
            <td class="px-5 py-3">
              <AppBadge :tone="stockStatus(p).tone">{{ stockStatus(p).label }}</AppBadge>
            </td>
            <td v-if="canManage" class="px-5 py-3 text-right">
              <button class="mr-3 text-brand-600 hover:underline" @click="openEdit(p)">Edit</button>
              <button class="text-red-600 hover:underline" @click="askDelete(p)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" :title="editing ? 'Edit Product' : 'New Product'" wide>
      <form class="space-y-4" @submit.prevent="save">
        <div class="grid grid-cols-2 gap-4">
          <AppInput v-model="form.name" label="Name" required />
          <AppInput v-model="form.sku" label="SKU" required />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <AppInput v-model="form.barcode" label="Barcode" />
          <AppSelect v-model="form.category_id" label="Category" :options="categoryOptions" placeholder="None" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <AppInput v-model="form.brand" label="Brand" />
          <AppSelect v-model="form.unit_type" label="Unit Type" :options="unitOptions" required />
        </div>
        <AppInput v-model="form.description" label="Description" />
        <div class="grid grid-cols-3 gap-4">
          <AppInput v-model="form.purchase_price" label="Purchase Price" type="number" step="0.01" min="0" required />
          <AppInput v-model="form.selling_price" label="Selling Price" type="number" step="0.01" min="0" required />
          <AppInput v-model="form.minimum_stock" label="Minimum Stock" type="number" step="0.01" min="0" required />
        </div>
        <AppSelect
          v-model="form.status"
          label="Status"
          :options="[{ value: 'active', label: 'Active' }, { value: 'inactive', label: 'Inactive' }]"
        />

        <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>

        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>

    <ConfirmDialog
      v-model="showConfirm"
      title="Delete product"
      :message="`Delete '${pendingDelete?.name}'? This cannot be undone.`"
      @confirm="confirmDelete"
    />
  </div>
</template>
