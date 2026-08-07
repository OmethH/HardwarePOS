<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
    type?: 'button' | 'submit'
    disabled?: boolean
  }>(),
  { variant: 'primary', type: 'button', disabled: false },
)
</script>

<template>
  <button
    :type="type"
    :disabled="disabled"
    class="inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50"
    :class="{
      'bg-brand-600 text-white hover:bg-brand-700': variant === 'primary',
      'bg-white text-slate-700 border border-slate-300 hover:bg-slate-50': variant === 'secondary',
      'bg-red-600 text-white hover:bg-red-700': variant === 'danger',
      'text-slate-600 hover:bg-slate-100': variant === 'ghost',
    }"
  >
    <slot />
  </button>
</template>
