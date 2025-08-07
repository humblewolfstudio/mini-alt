<script setup lang="ts">

import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {fetchListBuckets} from "../sources/BucketsDataSource";
import {fetchCreateEvent} from "../sources/EventsDataSource";

const router = useRouter()

const name = ref('')
const description = ref('')
const selectedBucket = ref('0')
const endpoint = ref('')
const token = ref('')

const buckets = ref([])

const isLoading = ref(false)
const error = ref<string | null>(null)

const createEvents = async () => {
  try {
    if (selectedBucket.value === '0' || endpoint.value === '') return

    const bucket = Number.parseInt(selectedBucket.value)

    isLoading.value = false
    error.value = null

    await fetchCreateEvent({
      name: name.value,
      description: description.value,
      bucket,
      endpoint: endpoint.value,
      token: token.value
    })
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

    buckets.value = await fetchListBuckets()

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
        <label for="description">Description</label>
        <input
            id="description"
            v-model="description"
            type="text"
            placeholder="Event description"
            :disabled="isLoading"
        />
        <p class="hint">The description for the event is optional.</p>
      </div>

      <div class="form-group">
        <label for="category">Category</label>
        <select id="category" v-model="selectedBucket" :disabled="isLoading">
          <option value="0">Select Bucket</option>
          <option v-for="bucket in buckets" :value="bucket.Id" :key="bucket.Id">{{ bucket.Name }}</option>
        </select>
        <p class="hint">Choose a bucket for your event.</p>
      </div>

      <div class="form-group">
        <label for="endpoint">Endpoint</label>
        <input
            id="endpoint"
            v-model="endpoint"
            type="text"
            placeholder="Webhook endpoint"
            :disabled="isLoading"
        />
        <p class="hint">The endpoint is required.</p>
      </div>

      <div class="form-group">
        <label for="token">Token</label>
        <input
            id="token"
            v-model="token"
            type="text"
            placeholder="Webhook authentication token"
            :disabled="isLoading"
        />
        <p class="hint">If the webhook requires an auth token or JWT.</p>
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