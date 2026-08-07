<script setup lang="ts">
withDefaults(defineProps<{ label: string; value: string; tone?: 'brand' | 'green' | 'amber' | 'red' }>(), {
  tone: 'brand',
})
</script>

<template>
  <div class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <p class="text-sm font-medium text-slate-500">{{ label }}</p>
    <p
      class="mt-2 text-2xl font-semibold"
      :class="{
        'text-brand-700': tone === 'brand',
        'text-emerald-700': tone === 'green',
        'text-amber-700': tone === 'amber',
        'text-red-700': tone === 'red',
      }"
    >
      {{ value }}
    </p>
  </div>
</template>
