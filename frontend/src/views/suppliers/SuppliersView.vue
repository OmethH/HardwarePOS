<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as suppliersApi from '@/api/suppliers'
import type { Supplier } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppModal from '@/components/ui/AppModal.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'

const auth = useAuthStore()
const canManage = auth.hasRole('admin', 'manager')

const suppliers = ref<Supplier[]>([])
const loading = ref(true)
const error = ref('')
const search = ref('')

const showModal = ref(false)
const editing = ref<Supplier | null>(null)
const form = ref({ name: '', phone: '', email: '', address: '', company: '' })
const saving = ref(false)

const showConfirm = ref(false)
const pendingDelete = ref<Supplier | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await suppliersApi.listSuppliers(search.value)
    suppliers.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load suppliers')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.value = { name: '', phone: '', email: '', address: '', company: '' }
  showModal.value = true
}

function openEdit(supplier: Supplier) {
  editing.value = supplier
  form.value = {
    name: supplier.name,
    phone: supplier.phone,
    email: supplier.email,
    address: supplier.address,
    company: supplier.company,
  }
  showModal.value = true
}

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (editing.value) {
      await suppliersApi.updateSupplier(editing.value.id, form.value)
    } else {
      await suppliersApi.createSupplier(form.value)
    }
    showModal.value = false
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to save supplier')
  } finally {
    saving.value = false
  }
}

function askDelete(supplier: Supplier) {
  pendingDelete.value = supplier
  showConfirm.value = true
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await suppliersApi.deleteSupplier(pendingDelete.value.id)
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to delete supplier')
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Suppliers</h1>
        <p class="text-sm text-slate-500">Vendors you purchase inventory from</p>
      </div>
      <AppButton v-if="canManage" @click="openCreate">+ New Supplier</AppButton>
    </div>

    <div class="flex gap-2">
      <AppInput v-model="search" placeholder="Search suppliers…" @update:model-value="load" class="max-w-xs" />
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Name</th>
            <th class="px-5 py-3">Company</th>
            <th class="px-5 py-3">Phone</th>
            <th class="px-5 py-3">Email</th>
            <th v-if="canManage" class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="5">Loading…</td>
          </tr>
          <tr v-else-if="suppliers.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="5">No suppliers yet.</td>
          </tr>
          <tr v-for="s in suppliers" :key="s.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ s.name }}</td>
            <td class="px-5 py-3 text-slate-600">{{ s.company }}</td>
            <td class="px-5 py-3 text-slate-600">{{ s.phone }}</td>
            <td class="px-5 py-3 text-slate-600">{{ s.email }}</td>
            <td v-if="canManage" class="px-5 py-3 text-right">
              <button class="mr-3 text-brand-600 hover:underline" @click="openEdit(s)">Edit</button>
              <button class="text-red-600 hover:underline" @click="askDelete(s)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" :title="editing ? 'Edit Supplier' : 'New Supplier'">
      <form class="space-y-4" @submit.prevent="save">
        <AppInput v-model="form.name" label="Name" required />
        <AppInput v-model="form.company" label="Company" />
        <div class="grid grid-cols-2 gap-4">
          <AppInput v-model="form.phone" label="Phone" />
          <AppInput v-model="form.email" label="Email" type="email" />
        </div>
        <AppInput v-model="form.address" label="Address" />
        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>

    <ConfirmDialog
      v-model="showConfirm"
      title="Delete supplier"
      :message="`Delete '${pendingDelete?.name}'? This cannot be undone.`"
      @confirm="confirmDelete"
    />
  </div>
</template>
