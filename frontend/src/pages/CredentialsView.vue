<script setup lang="ts">

import {onMounted, ref} from "vue";

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
      if (data) credentials.value = data.credentials
    }
  } catch (err) {
    console.error("Error fetcing credentials:", err)
    error.value = err instanceof Error ? err.message : "Failed to load credentials"
  } finally {
    isLoading.value = false
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
            <th>Created</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="cred in credentials" :key="cred.Id">
            <td>{{cred.AccessKey}}</td>
            <td>{{cred.CreatedAt}}</td>
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
  </div>
</template>

<style scoped>
.error-message {
  color: #ff4444;
  padding: 15px;
  text-align: center;
  background-color: #ffeeee;
  border-radius: 4px;
  margin-top: 10px;
}

.empty-message {
  padding: 15px;
  text-align: center;
  color: #888;
}
</style>