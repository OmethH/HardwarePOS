<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as categoriesApi from '@/api/categories'
import type { Category } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppModal from '@/components/ui/AppModal.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'
import { formatDate } from '@/lib/format'

const auth = useAuthStore()
const canManage = auth.hasRole('admin', 'manager')

const categories = ref<Category[]>([])
const loading = ref(true)
const error = ref('')

const showModal = ref(false)
const editing = ref<Category | null>(null)
const name = ref('')
const saving = ref(false)

const showConfirm = ref(false)
const pendingDelete = ref<Category | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await categoriesApi.listCategories()
    categories.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load categories')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  name.value = ''
  showModal.value = true
}

function openEdit(category: Category) {
  editing.value = category
  name.value = category.name
  showModal.value = true
}

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (editing.value) {
      await categoriesApi.updateCategory(editing.value.id, name.value)
    } else {
      await categoriesApi.createCategory(name.value)
    }
    showModal.value = false
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to save category')
  } finally {
    saving.value = false
  }
}

function askDelete(category: Category) {
  pendingDelete.value = category
  showConfirm.value = true
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await categoriesApi.deleteCategory(pendingDelete.value.id)
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to delete category')
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Categories</h1>
        <p class="text-sm text-slate-500">Organize your products into categories</p>
      </div>
      <AppButton v-if="canManage" @click="openCreate">+ New Category</AppButton>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Name</th>
            <th class="px-5 py-3">Created</th>
            <th v-if="canManage" class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="3">Loading…</td>
          </tr>
          <tr v-else-if="categories.length === 0">
            <td class="px-5 py-4 text-slate-400" colspan="3">No categories yet.</td>
          </tr>
          <tr v-for="cat in categories" :key="cat.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ cat.name }}</td>
            <td class="px-5 py-3 text-slate-500">{{ formatDate(cat.created_at) }}</td>
            <td v-if="canManage" class="px-5 py-3 text-right">
              <button class="mr-3 text-brand-600 hover:underline" @click="openEdit(cat)">Edit</button>
              <button class="text-red-600 hover:underline" @click="askDelete(cat)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" :title="editing ? 'Edit Category' : 'New Category'">
      <form class="space-y-4" @submit.prevent="save">
        <AppInput v-model="name" label="Name" required />
        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>

    <ConfirmDialog
      v-model="showConfirm"
      title="Delete category"
      :message="`Delete '${pendingDelete?.name}'? This cannot be undone.`"
      @confirm="confirmDelete"
    />
  </div>
</template>
