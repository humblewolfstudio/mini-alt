<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const bucketName = ref('');
const isLoading = ref(false);
const error = ref<string | null>(null);
const success = ref<string | null>(null);

const createBucket = async () => {
  if (!bucketName.value.trim()) {
    error.value = 'Bucket name is required';
    return;
  }

  try {
    isLoading.value = true;
    error.value = null;
    success.value = null;

    const res = await fetch('/api/buckets', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name: bucketName.value }),
    });

    if (res.ok) {
      success.value = 'Bucket created successfully!';
      bucketName.value = '';
      setTimeout(() => {
        router.push('/buckets');
      }, 1500);
    } else {
      const data = await res.json();
      error.value = data.message || 'Failed to create bucket';
    }
  } catch (err) {
    console.error("Error creating bucket:", err);
    error.value = err instanceof Error ? err.message : "Failed to create bucket";
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Create New Bucket</h1>
      <button @click="router.push('/buckets')">Back to Buckets</button>
    </div>

    <div class="form-container">
      <form @submit.prevent="createBucket" class="bucket-form">
        <div class="form-group">
          <label for="bucketName">Bucket Name</label>
          <input
              id="bucketName"
              v-model="bucketName"
              type="text"
              placeholder="Enter bucket name"
              :disabled="isLoading"
          />
          <p class="hint">Bucket names should be lowercase and can contain only letters, numbers, and hyphens.</p>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="isLoading">
            <span v-if="isLoading" class="spinner"></span>
            {{ isLoading ? 'Creating...' : 'Create Bucket' }}
          </button>
          <button type="button" @click="router.push('/buckets')" :disabled="isLoading">
            Cancel
          </button>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-if="success" class="success-message">
          {{ success }}
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>