<script setup lang="ts">
import {useRouter} from "vue-router";
import {ref} from "vue";
import ShowCredentialsModal from "../components/credentials/ShowCredentialsModal.vue";
import {fetchCreateCredentials} from "../sources/CredentialsDataSource";

const router = useRouter()
const accessKey = ref('')
const secretKey = ref('')

const name = ref('')
const description = ref('')
const expiresAt = ref<string>('');

const showCredentialsModal = ref(false)
const isLoading = ref(false)
const error = ref<string | null>(null)

const createCredentials = async () => {
    try {
      isLoading.value = true
      error.value = null
      showCredentialsModal.value = false

      const data = await fetchCreateCredentials({
        expiresAt: expiresAt.value ? new Date(expiresAt.value).toISOString() : null,
        name: name.value,
        description: description.value
      })

      accessKey.value = data.access_key
      secretKey.value = data.secret_key
      showCredentialsModal.value = true
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
        <div class="form-group">
          <label for="name">Name</label>
          <input
              id="name"
              v-model="name"
              type="text"
              placeholder="Credentials name"
              :disabled="isLoading"
          />
          <p class="hint">The name for the credentials is optional.</p>
        </div>
        <div class="form-group">
          <label for="description">Description</label>
          <input
              id="description"
              v-model="description"
              type="text"
              placeholder="Credentials description"
              :disabled="isLoading"
          />
          <p class="hint">The description for the credentials is optional.</p>
        </div>
        <div class="form-group">
          <label for="expiresAt">Expiration Date (optional)</label>
          <input
              type="date"
              id="expiresAt"
              v-model="expiresAt"
              :disabled="isLoading"
          />
          <p class="hint">Leave blank for no expiration.</p>
        </div>
        
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