<script setup lang="ts">
import {useRouter} from "vue-router";
import {ref} from "vue";
import ShowCredentialsModal from "../components/credentials/ShowCredentialsModal.vue";

const router = useRouter()
const accessKey = ref('')
const secretKey = ref('')

const showCredentialsModal = ref(false)
const isLoading = ref(false)
const error = ref<string | null>(null)

const createCredentials = async () => {
    try {
      isLoading.value = true
      error.value = null
      showCredentialsModal.value = false

      const res = await fetch('/api/credentials', {
        method: 'POST'
      })

      const data = await res.json()
      if (res.ok) {
        accessKey.value = data.access_key
        secretKey.value = data.secret_key
        showCredentialsModal.value = true
      } else {
        error.value = data.message || 'Failed to create credentials'
      }
    } catch (err) {
      console.error('Error creating credentials:', err)
      error.value = err instanceof Error ? err.message : 'Failed to create credentials'
    } finally {
      isLoading.value = false
    }
}

</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Create Credentials</h1>
      <button @click="router.push('/credentials')">Back to Credentials</button>
    </div>

    <div class="form-container">
      <form @submit.prevent="createCredentials" class="form">
        <div class="form-actions">
          <button type="submit" :disabled="isLoading">
            Create Credentials
          </button>
          <button type="button" @click="router.push('/credentials')" :disabled="isLoading">Cancel</button>
        </div>

        <div v-if="error" class="error-message">{{error}}</div>
      </form>
    </div>

    <ShowCredentialsModal
      v-if="showCredentialsModal"
      :accessKey="accessKey"
      :secretKey="secretKey"
      @continue="router.push('/credentials')"
      />
  </div>
</template>

<style scoped>

</style>