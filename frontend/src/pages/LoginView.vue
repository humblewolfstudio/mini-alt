<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
const success = ref('');

const router = useRouter();

const login = async () => {
  error.value = '';
  success.value = '';
  loading.value = true;

  try {
    const res = await fetch('/api/users/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });

    const data = await res.json();

    if (!res.ok) {
      error.value = data.error || 'Login failed';
    }

    success.value = 'Logged in successfully!';
    await router.push('/');
  } catch (err: any) {
    error.value = err.message || 'Something went wrong.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Login</h1>
    </div>

    <div class="form-container">
      <form class="form" @submit.prevent="login">
        <div class="form-group">
          <label for="username">Username</label>
          <input
              id="username"
              v-model="username"
              type="text"
              placeholder="Enter your username"
              required
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
              id="password"
              v-model="password"
              type="password"
              placeholder="Enter your password"
              required
          />
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            <span v-else>Login</span>
          </button>
          <button type="button" @click="() => { username = ''; password = '' }">
            Clear
          </button>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 500px;
  margin: 50px auto;
}
</style>