<script setup lang="ts">
withDefaults(defineProps<{ tone?: 'green' | 'amber' | 'red' | 'slate' | 'blue' }>(), { tone: 'slate' })
</script>

<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
    :class="{
      'bg-emerald-100 text-emerald-800': tone === 'green',
      'bg-amber-100 text-amber-800': tone === 'amber',
      'bg-red-100 text-red-800': tone === 'red',
      'bg-slate-100 text-slate-700': tone === 'slate',
      'bg-blue-100 text-blue-800': tone === 'blue',
    }"
  >
    <slot />
  </span>
</template>
