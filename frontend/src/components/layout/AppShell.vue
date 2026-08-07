<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import Topbar from './Topbar.vue'
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-slate-50">
    <Sidebar class="no-print" />
    <div class="flex flex-1 flex-col overflow-hidden">
      <Topbar class="no-print" />
      <main class="flex-1 overflow-y-auto p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
