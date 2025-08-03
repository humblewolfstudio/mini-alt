<script setup lang="ts">

import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";

const router = useRouter()

const name = ref('')
const selectedBucket = ref('0')

const buckets = ref([])

const isLoading = ref(false)
const error = ref<string | null>(null)

const createEvents = async () => {
  try {
    if (selectedBucket.value === '0') return

    const bucket = Number.parseInt(selectedBucket.value)

    isLoading.value = false
    error.value = null

    const res = await fetch('/api/events', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: name.value,
        bucket: bucket
      })
    })

    const data = await res.json()
    if (res.ok) {
      await router.push('/events')
    } else {
      error.value = data.message || 'Failed to create event'
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to create event'
  } finally {
    isLoading.value = false
  }
}

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
  fetchBuckets()
})
</script>

<template>
<div class="container">
  <div class="header">
    <h1>Create Event</h1>
    <button @click="router.push('/events')">Back to Events</button>
  </div>

  <div class="form-container">
    <form @submit.prevent="createEvents" class="form">
      <div class="form-group">
        <label for="name">Name</label>
        <input
            id="name"
            v-model="name"
            type="text"
            placeholder="Event name"
            :disabled="isLoading"
        />
        <p class="hint">The name for the event is optional.</p>
      </div>

      <div class="form-group">
        <label for="category">Category</label>
        <select id="category" v-model="selectedBucket" :disabled="isLoading">
          <option value="0">Select Bucket</option>
          <option v-for="bucket in buckets" :value="bucket.Id" :key="bucket.Id">{{ bucket.Name }}</option>
        </select>
        <p class="hint">Choose a bucket for your event.</p>
      </div>

      <div class="form-actions">
        <button type="submit" :disabled="isLoading">
          Create Event
        </button>
        <button type="button" @click="router.push('/events')" :disabled="isLoading">Cancel</button>
      </div>

      <div v-if="error" class="error-message">{{ error }}</div>
    </form>
  </div>
</div>
</template>

<style scoped>

</style>