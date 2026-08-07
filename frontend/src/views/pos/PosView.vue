<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as productsApi from '@/api/products'
import * as customersApi from '@/api/customers'
import * as salesApi from '@/api/sales'
import * as settingsApi from '@/api/settings'
import type { Customer, Product, Sale, Settings } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import ReceiptModal from '@/components/receipt/ReceiptModal.vue'
import { formatMoney, formatQty } from '@/lib/format'

interface CartLine {
  product: Product
  quantity: number
}

const products = ref<Product[]>([])
const customers = ref<Customer[]>([])
const settings = ref<Settings | null>(null)
const search = ref('')
const cart = ref<CartLine[]>([])
const customerId = ref('1')
const discount = ref('0')
const paymentMethod = ref('cash')
const error = ref('')
const submitting = ref(false)

const completedSale = ref<Sale | null>(null)
const showReceipt = ref(false)

const customerOptions = computed(() => customers.value.map((c) => ({ value: c.id, label: c.name })))

const filteredProducts = computed(() => {
  const term = search.value.trim().toLowerCase()
  if (!term) return products.value.slice(0, 24)
  return products.value
    .filter(
      (p) =>
        p.name.toLowerCase().includes(term) ||
        p.sku.toLowerCase().includes(term) ||
        (p.barcode ?? '').toLowerCase().includes(term),
    )
    .slice(0, 24)
})

const subtotal = computed(() => cart.value.reduce((sum, l) => sum + l.product.selling_price * l.quantity, 0))
const total = computed(() => Math.max(0, subtotal.value - (Number(discount.value) || 0)))

function inCartQuantity(productId: number): number {
  return cart.value.find((l) => l.product.id === productId)?.quantity ?? 0
}

function addToCart(product: Product) {
  if (product.current_stock <= 0) return
  const existing = cart.value.find((l) => l.product.id === product.id)
  if (existing) {
    if (existing.quantity < product.current_stock) existing.quantity += 1
  } else {
    cart.value.push({ product, quantity: 1 })
  }
}

function setQuantity(line: CartLine, qty: number) {
  if (qty <= 0) {
    cart.value = cart.value.filter((l) => l !== line)
    return
  }
  line.quantity = Math.min(qty, line.product.current_stock)
}

function removeLine(line: CartLine) {
  cart.value = cart.value.filter((l) => l !== line)
}

function resetCart() {
  cart.value = []
  discount.value = '0'
  customerId.value = '1'
  paymentMethod.value = 'cash'
}

async function loadProducts() {
  const res = await productsApi.listProducts({ status: 'active' })
  products.value = res.data.data
}

async function loadCustomers() {
  const res = await customersApi.listCustomers()
  customers.value = res.data.data
}

async function loadSettings() {
  const res = await settingsApi.getSettings()
  settings.value = res.data.data
}

async function completeSale() {
  if (cart.value.length === 0) return
  error.value = ''
  submitting.value = true
  try {
    const res = await salesApi.checkout({
      customer_id: Number(customerId.value),
      discount: Number(discount.value) || 0,
      payment_method: paymentMethod.value,
      items: cart.value.map((l) => ({ product_id: l.product.id, quantity: l.quantity })),
    })
    completedSale.value = res.data.data
    showReceipt.value = true
    resetCart()
    await loadProducts()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Checkout failed')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadProducts(), loadCustomers(), loadSettings()])
})
</script>

<template>
  <div class="grid h-full grid-cols-1 gap-6 lg:grid-cols-3">
    <div class="lg:col-span-2 space-y-4">
      <AppInput v-model="search" placeholder="Search by name, SKU, or scan barcode…" />

      <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-4">
        <button
          v-for="p in filteredProducts"
          :key="p.id"
          class="rounded-xl border border-slate-200 bg-white p-3 text-left shadow-sm transition hover:border-brand-300 hover:shadow disabled:cursor-not-allowed disabled:opacity-40"
          :disabled="p.current_stock <= 0"
          @click="addToCart(p)"
        >
          <p class="text-sm font-medium text-slate-900 line-clamp-2">{{ p.name }}</p>
          <p class="mt-1 text-xs text-slate-500">{{ p.sku }}</p>
          <div class="mt-2 flex items-center justify-between">
            <span class="text-sm font-semibold text-brand-700">{{ formatMoney(p.selling_price) }}</span>
            <span class="text-xs text-slate-400">
              {{ p.current_stock <= 0 ? 'Out of stock' : `${formatQty(p.current_stock)} left` }}
            </span>
          </div>
          <span v-if="inCartQuantity(p.id) > 0" class="mt-1 block text-xs font-medium text-emerald-600">
            {{ inCartQuantity(p.id) }} in cart
          </span>
        </button>
        <p v-if="filteredProducts.length === 0" class="col-span-full py-8 text-center text-sm text-slate-400">
          No products match your search.
        </p>
      </div>
    </div>

    <AppCard title="Current Sale" class="flex h-fit flex-col">
      <div class="space-y-3">
        <div v-if="cart.length === 0" class="py-6 text-center text-sm text-slate-400">Cart is empty</div>
        <div v-for="line in cart" :key="line.product.id" class="flex items-center gap-2 text-sm">
          <div class="flex-1">
            <p class="font-medium text-slate-900">{{ line.product.name }}</p>
            <p class="text-xs text-slate-500">{{ formatMoney(line.product.selling_price) }} each</p>
          </div>
          <input
            type="number"
            min="0"
            :max="line.product.current_stock"
            :value="line.quantity"
            @input="setQuantity(line, Number(($event.target as HTMLInputElement).value))"
            class="w-16 rounded-md border border-slate-300 px-2 py-1 text-right text-sm"
          />
          <span class="w-20 text-right font-medium text-slate-900">
            {{ formatMoney(line.product.selling_price * line.quantity) }}
          </span>
          <button class="text-red-500 hover:text-red-700" @click="removeLine(line)">✕</button>
        </div>
      </div>

      <div class="mt-4 space-y-3 border-t border-slate-100 pt-4">
        <AppSelect v-model="customerId" label="Customer" :options="customerOptions" />
        <div class="grid grid-cols-2 gap-3">
          <AppInput v-model="discount" label="Discount" type="number" min="0" step="0.01" />
          <AppSelect
            v-model="paymentMethod"
            label="Payment"
            :options="[
              { value: 'cash', label: 'Cash' },
              { value: 'card', label: 'Card' },
              { value: 'mobile_money', label: 'Mobile Money' },
            ]"
          />
        </div>

        <div class="space-y-1 border-t border-slate-100 pt-3 text-sm">
          <div class="flex justify-between text-slate-500">
            <span>Subtotal</span>
            <span>{{ formatMoney(subtotal) }}</span>
          </div>
          <div class="flex justify-between text-slate-500">
            <span>Discount</span>
            <span>-{{ formatMoney(Number(discount) || 0) }}</span>
          </div>
          <div class="flex justify-between text-base font-semibold text-slate-900">
            <span>Total</span>
            <span>{{ formatMoney(total) }}</span>
          </div>
        </div>

        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

        <AppButton class="w-full justify-center" :disabled="cart.length === 0 || submitting" @click="completeSale">
          {{ submitting ? 'Processing…' : 'Complete Sale' }}
        </AppButton>
      </div>
    </AppCard>

    <ReceiptModal
      v-if="completedSale && settings"
      v-model="showReceipt"
      :sale="completedSale"
      :settings="settings"
    />
  </div>
</template>
