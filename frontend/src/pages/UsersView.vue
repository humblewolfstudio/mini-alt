<script setup lang="ts">

import {onMounted, ref} from "vue";
import {getLocaleDate, getLocaleDateTime} from "../utils";
import DeleteModal from "../components/modals/DeleteModal.vue";

const showDeleteModal = ref(false)

const selectedUser = ref('')

const isLoading = ref(false)
const users = ref([])
const error = ref<string | null>(null)

const fetchUsers = async () => {
  try {
    isLoading.value = true
    error.value = null

    const res = await fetch('/api/users/list')

    if (res.ok) {
      const data = await res.json()
      if (data) users.value = data
    }
  } catch (err) {
    console.error('Error fetching users:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load users'
  } finally {
    isLoading.value = false
  }
}

const promptDelete = (username: string) => {
  selectedUser.value = username
  showDeleteModal.value = true
}

const handleDelete = async () => {
  if (selectedUser.value === '') return

  const userId = users.value.find((user) => user.Username === selectedUser.value).Id

  try {
    await fetch('/api/users/delete', {
      method: 'POST',
      body: JSON.stringify({
        id: userId
      })
    })

    selectedUser.value = ''
    await fetchUsers()
  } catch (err) {
    error.value = `Failed to delete ${selectedUser.value}`
  } finally {
    showDeleteModal.value = false
  }
}

onMounted(() => {
  fetchUsers()
})

</script>

<template>
<div class="container">
  <div class="header">
    <h1>Users</h1>
    <RouterLink to="/users/create-users">Create User</RouterLink>
  </div>
  <div class="table-content">
    <div class="table-wrapper">
      <table class="table">
        <thead>
        <tr>
          <th>Username</th>
          <th>Expires</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="user in users">
          <td>{{user.Username}}</td>
          <td>{{user.ExpiresAt ? getLocaleDate(user.ExpiresAt) : 'Never'}}</td>
          <td>{{ getLocaleDateTime(user.CreatedAt)}}</td>
          <td>
            <button @click="promptDelete(user.Username)">Delete</button>
          </td>
        </tr>
        </tbody>
      </table>

      <div v-if="isLoading" class="loading-row">
        <div class="spinner-container">
          <div class="spinner"></div>
          <div>Loading users...</div>
        </div>
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <div v-if="!isLoading && users.length === 0 && !error" class="empty-message">
        No users found
      </div>
    </div>
  </div>

  <DeleteModal
    v-if="showDeleteModal"
    :content="selectedUser"
    @close="showDeleteModal = false"
    @confirm="handleDelete"
  />
</div>
</template>

<style scoped>

</style>