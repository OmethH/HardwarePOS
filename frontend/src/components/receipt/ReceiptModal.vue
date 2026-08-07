<script setup lang="ts">
import { ref } from 'vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppButton from '@/components/ui/AppButton.vue'
import ReceiptDocument from './ReceiptDocument.vue'
import type { Sale, Settings } from '@/types'

defineProps<{ modelValue: boolean; sale: Sale; settings: Settings }>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()

const format = ref<'a4' | 'thermal'>('a4')

function print() {
  window.print()
}
</script>

<template>
  <AppModal :model-value="modelValue" title="Receipt" wide @update:model-value="$emit('update:modelValue', $event)">
    <div class="mb-4 flex items-center justify-between no-print">
      <div class="flex gap-2 text-sm">
        <button
          class="rounded-md px-3 py-1.5 font-medium"
          :class="format === 'a4' ? 'bg-brand-600 text-white' : 'bg-slate-100 text-slate-600'"
          @click="format = 'a4'"
        >
          A4 Invoice
        </button>
        <button
          class="rounded-md px-3 py-1.5 font-medium"
          :class="format === 'thermal' ? 'bg-brand-600 text-white' : 'bg-slate-100 text-slate-600'"
          @click="format = 'thermal'"
        >
          Thermal
        </button>
      </div>
      <AppButton @click="print">Print</AppButton>
    </div>

    <div class="rounded-lg border border-slate-200 bg-slate-50 py-4">
      <ReceiptDocument :sale="sale" :settings="settings" :format="format" />
    </div>
  </AppModal>
</template>
