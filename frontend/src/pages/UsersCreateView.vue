<script setup lang="ts">

import {useRouter} from "vue-router";
import {ref} from "vue";

const router = useRouter()

const username = ref('')
const password = ref('')
const expiresAt = ref('')

const isLoading = ref(false)
const error = ref<string | null>(null)

const createUser = async () => {
  try {
    isLoading.value = true
    error.value = null

    const res = await fetch('/api/users/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
        expiresAt: expiresAt.value
      })
    })

    if(res.status === 401) {
      await router.push('/login')
    }

    const data = await res.json()
    if (res.ok) {
      await router.push('/users')
    } else {
      error.value = data.error || 'Failed to create user'
    }
  } catch (err) {
    console.error('Error creating user:', err)
    error.value = err instanceof Error ? err.message : 'Failed to create user'
  } finally {
    isLoading.value = false
  }
}

</script>

<template>
<div class="container">
  <div class="header">
    <h1>Create User</h1>
    <button @click="router.push('/users')">Back to Users</button>
  </div>

  <div class="form-container">
    <form @submit.prevent="createUser" class="form">
      <div class="form-group">
        <label for="username">Username</label>
        <input
            type="text"
            id="username"
            v-model="username"
            :disabled="isLoading"
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
            type="password"
            id="password"
            v-model="password"
            :disabled="isLoading"
        />
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
          Create User
        </button>
        <button type="button" @click="router.push('/users')" :disabled="isLoading">Cancel</button>
      </div>

      <div v-if="error" class="error-message">{{error}}</div>
    </form>
  </div>
</div>
</template>

<style scoped>

</style>