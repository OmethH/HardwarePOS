<script setup lang="ts">
import { computed } from 'vue'
import type { Sale, Settings } from '@/types'
import { formatDateTime, formatMoney, formatQty } from '@/lib/format'

const props = withDefaults(
  defineProps<{ sale: Sale; settings: Settings; format?: 'a4' | 'thermal' }>(),
  { format: 'a4' },
)

const isThermal = computed(() => props.format === 'thermal')
</script>

<template>
  <div
    id="receipt-print-area"
    class="mx-auto bg-white text-slate-900"
    :class="isThermal ? 'w-[80mm] p-2 text-xs' : 'w-full max-w-2xl p-8 text-sm'"
  >
    <div class="text-center">
      <h1 :class="isThermal ? 'text-sm font-bold' : 'text-xl font-bold'">{{ settings.shop_name }}</h1>
      <p v-if="settings.address" class="text-slate-500">{{ settings.address }}</p>
      <p v-if="settings.phone" class="text-slate-500">{{ settings.phone }}</p>
    </div>

    <div class="my-3 border-t border-dashed border-slate-400" />

    <div class="flex justify-between">
      <span>Invoice:</span>
      <span class="font-medium">{{ sale.invoice_number }}</span>
    </div>
    <div class="flex justify-between">
      <span>Date:</span>
      <span>{{ formatDateTime(sale.created_at) }}</span>
    </div>
    <div class="flex justify-between">
      <span>Customer:</span>
      <span>{{ sale.customer?.name ?? 'Walk-in Customer' }}</span>
    </div>

    <div class="my-3 border-t border-dashed border-slate-400" />

    <table class="w-full">
      <thead v-if="!isThermal">
        <tr class="text-left text-slate-500">
          <th class="pb-1">Item</th>
          <th class="pb-1 text-right">Qty</th>
          <th class="pb-1 text-right">Price</th>
          <th class="pb-1 text-right">Total</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in sale.items" :key="item.id" class="align-top">
          <td class="py-0.5" :colspan="isThermal ? 1 : 1">{{ item.product?.name }}</td>
          <td v-if="isThermal" class="py-0.5 text-right" colspan="3">
            {{ formatQty(item.quantity) }} x {{ formatMoney(item.selling_price) }} =
            {{ formatMoney(item.quantity * item.selling_price) }}
          </td>
          <template v-else>
            <td class="py-0.5 text-right">{{ formatQty(item.quantity) }}</td>
            <td class="py-0.5 text-right">{{ formatMoney(item.selling_price) }}</td>
            <td class="py-0.5 text-right">{{ formatMoney(item.quantity * item.selling_price) }}</td>
          </template>
        </tr>
      </tbody>
    </table>

    <div class="my-3 border-t border-dashed border-slate-400" />

    <div class="flex justify-between">
      <span>Subtotal:</span>
      <span>{{ formatMoney(sale.subtotal) }}</span>
    </div>
    <div class="flex justify-between">
      <span>Discount:</span>
      <span>{{ formatMoney(sale.discount) }}</span>
    </div>
    <div class="flex justify-between font-bold" :class="isThermal ? 'text-sm' : 'text-base'">
      <span>TOTAL:</span>
      <span>{{ formatMoney(sale.total) }}</span>
    </div>
    <div class="mt-1 flex justify-between text-slate-500">
      <span>Payment:</span>
      <span class="capitalize">{{ sale.payment_method }}</span>
    </div>

    <div class="my-3 border-t border-dashed border-slate-400" />

    <p class="text-center text-slate-600">{{ settings.receipt_footer }}</p>
  </div>
</template>
