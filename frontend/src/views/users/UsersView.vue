<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import * as usersApi from '@/api/users'
import type { Role, User } from '@/types'
import { apiErrorMessage } from '@/api/client'
import AppCard from '@/components/ui/AppCard.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppBadge from '@/components/ui/AppBadge.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'

const users = ref<User[]>([])
const roles = ref<Role[]>([])
const loading = ref(true)
const error = ref('')

const showModal = ref(false)
const editing = ref<User | null>(null)
const saving = ref(false)
const formError = ref('')
const form = ref({ name: '', username: '', password: '', role_id: '', active: true })

const showConfirm = ref(false)
const pendingDeactivate = ref<User | null>(null)

const roleOptions = computed(() => roles.value.map((r) => ({ value: r.id, label: r.name })))

async function load() {
  loading.value = true
  try {
    const res = await usersApi.listUsers()
    users.value = res.data.data
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load users')
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  const res = await usersApi.listRoles()
  roles.value = res.data.data
}

function openCreate() {
  editing.value = null
  form.value = { name: '', username: '', password: '', role_id: '', active: true }
  formError.value = ''
  showModal.value = true
}

function openEdit(user: User) {
  editing.value = user
  form.value = { name: user.name, username: user.username, password: '', role_id: String(user.role_id), active: user.active }
  formError.value = ''
  showModal.value = true
}

async function save() {
  saving.value = true
  formError.value = ''
  try {
    if (editing.value) {
      await usersApi.updateUser(editing.value.id, {
        name: form.value.name,
        role_id: Number(form.value.role_id),
        active: form.value.active,
        password: form.value.password || undefined,
      })
    } else {
      await usersApi.createUser({
        name: form.value.name,
        username: form.value.username,
        password: form.value.password,
        role_id: Number(form.value.role_id),
      })
    }
    showModal.value = false
    await load()
  } catch (err) {
    formError.value = apiErrorMessage(err, 'Failed to save user')
  } finally {
    saving.value = false
  }
}

function askDeactivate(user: User) {
  pendingDeactivate.value = user
  showConfirm.value = true
}

async function confirmDeactivate() {
  if (!pendingDeactivate.value) return
  try {
    await usersApi.deactivateUser(pendingDeactivate.value.id)
    await load()
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to deactivate user')
  }
}

onMounted(async () => {
  await Promise.all([load(), loadRoles()])
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">Users</h1>
        <p class="text-sm text-slate-500">Manage staff accounts and roles</p>
      </div>
      <AppButton @click="openCreate">+ New User</AppButton>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

    <AppCard :padded="false">
      <table class="w-full text-sm">
        <thead class="border-b border-slate-100 text-left text-xs font-medium uppercase text-slate-500">
          <tr>
            <th class="px-5 py-3">Name</th>
            <th class="px-5 py-3">Username</th>
            <th class="px-5 py-3">Role</th>
            <th class="px-5 py-3">Status</th>
            <th class="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-if="loading">
            <td class="px-5 py-4 text-slate-400" colspan="5">Loading…</td>
          </tr>
          <tr v-for="u in users" :key="u.id" class="hover:bg-slate-50">
            <td class="px-5 py-3 font-medium text-slate-900">{{ u.name }}</td>
            <td class="px-5 py-3 text-slate-600">{{ u.username }}</td>
            <td class="px-5 py-3 capitalize text-slate-600">{{ u.role.name }}</td>
            <td class="px-5 py-3">
              <AppBadge :tone="u.active ? 'green' : 'slate'">{{ u.active ? 'Active' : 'Inactive' }}</AppBadge>
            </td>
            <td class="px-5 py-3 text-right">
              <button class="mr-3 text-brand-600 hover:underline" @click="openEdit(u)">Edit</button>
              <button v-if="u.active" class="text-red-600 hover:underline" @click="askDeactivate(u)">Deactivate</button>
            </td>
          </tr>
        </tbody>
      </table>
    </AppCard>

    <AppModal v-model="showModal" :title="editing ? 'Edit User' : 'New User'">
      <form class="space-y-4" @submit.prevent="save">
        <AppInput v-model="form.name" label="Full Name" required />
        <AppInput v-model="form.username" label="Username" required :disabled="!!editing" />
        <AppInput
          v-model="form.password"
          label="Password"
          type="password"
          :placeholder="editing ? 'Leave blank to keep current password' : ''"
          :required="!editing"
        />
        <AppSelect v-model="form.role_id" label="Role" :options="roleOptions" placeholder="Select role" required />
        <label v-if="editing" class="flex items-center gap-2 text-sm text-slate-600">
          <input type="checkbox" v-model="form.active" class="rounded border-slate-300" />
          Active
        </label>

        <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>

        <div class="flex justify-end gap-2">
          <AppButton type="button" variant="secondary" @click="showModal = false">Cancel</AppButton>
          <AppButton type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</AppButton>
        </div>
      </form>
    </AppModal>

    <ConfirmDialog
      v-model="showConfirm"
      title="Deactivate user"
      :message="`Deactivate '${pendingDeactivate?.name}'? They will no longer be able to log in.`"
      @confirm="confirmDeactivate"
    />
  </div>
</template>
