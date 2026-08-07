<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as customersApi from '@/api/customers'
import type { Customer } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppModal from '@/components/ui/AppModal.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'

const auth = useAuthStore()
const canManage = auth.hasRole('admin', 'manager')

const customers = ref<Customer[]>([])
const loading = ref(true)
const error = ref('')
const search = ref('')

const showModal = ref(false)
const editing = ref<Customer | null>(null)
const form = ref({ name: '', phone: '', address: '' })
const saving = ref(false)

const showConfirm = ref(false)
const pendingDelete = ref<Customer | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await customersApi.listCustomers(search.value)
    customers.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load customers')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.value = { name: '', phone: '', address: '' }
  showModal.value = true
}

function openEdit(customer: Customer) {
  editing.value = customer
  form.value = { name: customer.name, phone: customer.phone, address: customer.address }
  showModal.value = true
}

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (editing.value) {
      await customersApi.updateCustomer(editing.value.id, form.value)
    } else {
      await customersApi.createCustomer(form.value)
    }
    showModal.value = false
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to save customer')
  } finally {
    saving.value = false
  }
}

function askDelete(customer: Customer) {
  pendingDelete.value = customer
  showConfirm.value = true
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await customersApi.deleteCustomer(pendingDelete.value.id)
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to delete customer')
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Customers</h1>
        <p class="text-sm text-slate-500">Optional customer records for sales</p>
      </div>
      <AppButton v-if="canManage" @click="openCreate">+ New Customer</AppButton>
    </div>

    <AppInput v-model="search" placeholder="Search customers…" @update:model-value="load" class="max-w-xs" />

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Name</th>
            <th class="px-5 py-3">Phone</th>
            <th class="px-5 py-3">Address</th>
            <th v-if="canManage" class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="4">Loading…</td>
          </tr>
          <tr v-for="c in customers" :key="c.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ c.name }}</td>
            <td class="px-5 py-3 text-slate-600">{{ c.phone }}</td>
            <td class="px-5 py-3 text-slate-600">{{ c.address }}</td>
            <td v-if="canManage" class="px-5 py-3 text-right">
              <template v-if="c.id !== 1">
                <button class="mr-3 text-brand-600 hover:underline" @click="openEdit(c)">Edit</button>
                <button class="text-red-600 hover:underline" @click="askDelete(c)">Delete</button>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" :title="editing ? 'Edit Customer' : 'New Customer'">
      <form class="space-y-4" @submit.prevent="save">
        <AppInput v-model="form.name" label="Name" required />
        <AppInput v-model="form.phone" label="Phone" />
        <AppInput v-model="form.address" label="Address" />
        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>

    <ConfirmDialog
      v-model="showConfirm"
      title="Delete customer"
      :message="`Delete '${pendingDelete?.name}'?`"
      @confirm="confirmDelete"
    />
  </div>
</template>
