<script setup lang="ts">

import {onMounted, ref} from "vue";
import {fetchSystemInfo, fetchSystemSpecs, SystemInfoResponse, SystemSpecsResponse} from "../sources/SystemDataSource";
import DiskUsageCard from "../components/DiskUsageCard.vue";
import {formatSize} from "../utils";

const isLoading = ref(false)
const error = ref<string | null>(null)

const info = ref<SystemInfoResponse | null>(null)
const specs = ref<SystemSpecsResponse | null>(null)

const fetchData = async () => {
  try {
    isLoading.value = true

    await Promise.all(
        [
          info.value = await fetchSystemInfo(),
          specs.value = await fetchSystemSpecs()
        ]
    )

  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load information"
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})

</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Server Information</h1>
    </div>

    <div v-if="error" class="error-message">{{ error }}</div>

    <div v-if="isLoading" class="spinner-container">
      <div class="spinner"></div>
      <span>Loading...</span>
    </div>

    <div v-else>
      <div style="display: flex; flex-wrap: wrap; gap: 20px; margin-bottom: 20px;">
        <div class="card">
          <h3>Buckets</h3>
          <p>{{ info?.NumberBuckets ?? '-' }}</p>
        </div>
        <div class="card">
          <h3>Objects</h3>
          <p>{{ info?.NumberObjects ?? '-' }}</p>
        </div>
        <div class="card">
          <h3>Reported Usage</h3>
          <p>{{ formatSize(info?.Usage) }}</p>
        </div>
      </div>

      <DiskUsageCard :specs="specs" />
    </div>
  </div>
</template>

<style scoped>
.card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  flex: 1 1 200px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  text-align: center;
}

.card h3 {
  margin: 0 0 10px;
  font-size: 1rem;
  color: #2c3e50;
}

.card p {
  font-size: 1.5rem;
  font-weight: bold;
  color: #42b983;
}
</style>