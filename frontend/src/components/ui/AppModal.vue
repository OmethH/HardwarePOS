<script setup lang="ts">
withDefaults(defineProps<{ modelValue: boolean; title?: string; wide?: boolean }>(), { wide: false })
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>

<template>
  <Teleport to="body">
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4">
      <div
        class="w-full rounded-xl bg-white shadow-xl"
        :class="wide ? 'max-w-2xl' : 'max-w-md'"
        @click.stop
      >
        <div class="flex items-center justify-between border-b border-slate-100 px-5 py-4">
          <h3 class="text-base font-semibold text-slate-900">{{ title }}</h3>
          <button
            class="rounded-md p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
            @click="$emit('update:modelValue', false)"
          >
            ✕
          </button>
        </div>
        <div class="max-h-[75vh] overflow-y-auto p-5">
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>
