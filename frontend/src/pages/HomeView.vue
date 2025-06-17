<script setup lang="ts">
import {onMounted, ref} from "vue";

const isLoading = ref(false);
const buckets = ref([]);
const error = ref<string | null>(null);

const fetchBuckets = async () => {
  try {
    isLoading.value = true;
    error.value = null;

    const res = await fetch('/api/buckets')

    if(res.ok) {
      const data = await res.json()
      if(data) buckets.value = data
    }

  } catch (err) {
    console.error("Error fetching buckets:", err);
    error.value = err instanceof Error ? err.message : "Failed to load buckets";
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchBuckets();
});
</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Buckets</h1>
      <button @click="fetchBuckets" :disabled="isLoading">
        {{ isLoading ? 'Refreshing...' : 'Refresh' }}
      </button>
    </div>
    <div class="table-content">
      <div class="table-wrapper">
        <table class="table">
          <thead>
          <tr>
            <th>Name</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="bucket in buckets" :key="bucket.Name">
            <td>{{ bucket.Name }}</td>
            <td>{{bucket.CreationDate}}</td>
            <td>
              <button @click="">View</button>
            </td>
          </tr>
          </tbody>
        </table>

        <div v-if="isLoading" class="loading-row">
          <div class="spinner-container">
            <div class="spinner"></div>
            <div>Loading buckets...</div>
          </div>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-if="!isLoading && buckets.length === 0 && !error" class="empty-message">
          No buckets found
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