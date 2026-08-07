<script setup lang="ts">
import AppModal from './AppModal.vue'
import AppButton from './AppButton.vue'

withDefaults(defineProps<{ modelValue: boolean; title?: string; message: string }>(), {
  title: 'Are you sure?',
})
const emit = defineEmits<{ 'update:modelValue': [value: boolean]; confirm: [] }>()

function confirm() {
  emit('confirm')
  emit('update:modelValue', false)
}
</script>

<template>
  <AppModal :model-value="modelValue" :title="title" @update:model-value="$emit('update:modelValue', $event)">
    <p class="text-sm text-slate-600">{{ message }}</p>
    <div class="mt-5 flex justify-end gap-2">
      <AppButton variant="secondary" @click="$emit('update:modelValue', false)">Cancel</AppButton>
      <AppButton variant="danger" @click="confirm">Confirm</AppButton>
    </div>
  </AppModal>
</template>
