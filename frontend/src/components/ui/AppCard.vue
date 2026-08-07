<script setup lang="ts">
withDefaults(defineProps<{ title?: string; padded?: boolean }>(), { padded: true })
</script>

<template>
  <div class="rounded-xl border border-slate-200 bg-white shadow-sm">
    <div v-if="title || $slots.actions" class="flex items-center justify-between border-b border-slate-100 px-5 py-4">
      <h2 v-if="title" class="text-base font-semibold text-slate-900">{{ title }}</h2>
      <div v-if="$slots.actions" class="flex items-center gap-2">
        <slot name="actions" />
      </div>
    </div>
    <div :class="padded ? 'p-5' : ''">
      <slot />
    </div>
  </div>
</template>
