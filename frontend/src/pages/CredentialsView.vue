<script setup lang="ts">

import {onMounted, ref} from "vue";
import DeleteModal from "../components/modals/DeleteModal.vue";
import {getLocaleDate} from "../utils";
import EditCredentialsModal from "../components/credentials/EditCredentialsModal.vue";

const showDeleteModal = ref(false)
const showEditModal = ref(false)
const data = ref<any | null>(null)

const selectedAccessKey = ref('')

const isLoading = ref(false)
const credentials = ref([])
const error = ref<string | null>(null)

const fetchCredentials = async () => {
  try {
    isLoading.value = true
    error.value = null

    const res = await fetch('/api/credentials')

    if (res.ok) {
      const data = await res.json()
      if (data) credentials.value = data
    }
  } catch (err) {
    console.error("Error fetching credentials:", err)
    error.value = err instanceof Error ? err.message : "Failed to load credentials"
  } finally {
    isLoading.value = false
  }
}

const promptDelete = (accessKey: string) => {
  selectedAccessKey.value = accessKey
  showDeleteModal.value = true
}

const promptEdit = (accessKey: string, editData: any) => {
  selectedAccessKey.value = accessKey
  showEditModal.value = true
  data.value = editData
}

const handleDelete = async () => {
  if (selectedAccessKey.value === '') return

  try {
    await fetch('/api/credentials/delete', {
      method: 'POST',
      body: JSON.stringify({
        accessKey: selectedAccessKey.value
      })
    })

    selectedAccessKey.value = ''
    await fetchCredentials()
  } catch (err) {
    error.value = `Failed to delete ${selectedAccessKey.value}`
  } finally {
    showDeleteModal.value = false
  }
}

const handleEdit = async (data: any) => {
  if (selectedAccessKey.value === '') return

  try {
    await fetch('/api/credentials/edit', {
      method: 'POST',
      body: JSON.stringify({
        accessKey: selectedAccessKey.value,
        name: data.name,
        description: data.description,
        status: data.status,
        expiresAt: data.expiresAt
      })
    })

    selectedAccessKey.value = ''
    await fetchCredentials()
  } catch (err) {
    error.value = `Failed to edit ${selectedAccessKey.value}`
  } finally {
    showEditModal.value = false
  }
}

onMounted(() => {
  fetchCredentials()
})

</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Credentials</h1>
      <RouterLink to="/credentials/create-credentials">Create Credentials</RouterLink>
    </div>
    <div class="table-content">
      <div class="table-wrapper">
        <table class="table">
          <thead>
          <tr>
            <th>Access Key</th>
            <th>Expires</th>
            <th>Status</th>
            <th>Name</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="cred in credentials" :key="cred.Id">
            <td>{{ cred.AccessKey }}</td>
            <td>{{ cred.ExpiresAt ? getLocaleDate(cred.ExpiresAt) : 'Never' }}</td>
            <td>{{ cred.Status ? 'Enabled' : 'Disabled' }}</td>
            <td>{{ cred.Name }}</td>
            <td>{{ cred.Description }}</td>
            <td>
              <button @click="promptDelete(cred.AccessKey)">Delete</button>
              <button @click="promptEdit(cred.AccessKey, cred)">Edit</button>
            </td>
          </tr>
          </tbody>
        </table>

        <div v-if="isLoading" class="loading-row">
          <div class="spinner-container">
            <div class="spinner"></div>
            <div>Loading credentials...</div>
          </div>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-if="!isLoading && credentials.length === 0 && !error" class="empty-message">
          No credentials found
        </div>
      </div>
    </div>

    <DeleteModal
        v-if="showDeleteModal"
        :content="selectedAccessKey"
        @close="showDeleteModal = false"
        @confirm="handleDelete"
      />

    <EditCredentialsModal
      v-if="showEditModal"
      :accessKey="selectedAccessKey"
      :data="data"
      @close="showEditModal = false"
      @continue="handleEdit"
      />
  </div>
</template>

<style scoped>

</style>