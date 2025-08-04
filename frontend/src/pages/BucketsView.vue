<script setup lang="ts">
import {onMounted, ref} from "vue";
import {formatSize, getLocaleDateTime} from "../utils";
import {useRouter} from "vue-router";

const router = useRouter()
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

    if(res.status === 401) {
      await router.push('/login')
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
      <RouterLink to="/buckets/create-bucket">Create Bucket</RouterLink>
    </div>
    <div class="table-content">
      <div class="table-wrapper">
        <table class="table">
          <thead>
          <tr>
            <th>Name</th>
            <th>Objects</th>
            <th>Size</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="bucket in buckets" :key="bucket.Name">
            <td>{{ bucket.Name }}</td>
            <td>{{ bucket.NumberObjects }}</td>
            <td>{{ formatSize(bucket.Size) }}</td>
            <td>{{ getLocaleDateTime(bucket.CreatedAt) }}</td>
            <td>
              <RouterLink :to="'/buckets/' + bucket.Name">View</RouterLink>
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

</style>