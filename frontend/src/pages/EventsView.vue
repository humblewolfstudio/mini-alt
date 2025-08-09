<script setup lang="ts">
import {onMounted, ref} from "vue";
import {getLocaleDateTime} from "../utils";
import {fetchListEvents} from "../sources/EventsDataSource";

const isLoading = ref(false)
const events = ref([])
const error = ref<string | null>(null)

const fetchEvents = async () => {
  try {
    isLoading.value = true
    error.value = null

    events.value = await fetchListEvents()
    
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load events'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchEvents()
})
</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Events</h1>
      <RouterLink to="/events/create-event">Create Event</RouterLink>
    </div>

    <div class="table-content">
      <div class="table-wrapper">
        <table class="table">
          <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Global</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="event in events" :key="event.Id">
            <td>{{ event.Name }}</td>
            <td>{{ event.Description }}</td>
            <td>{{ event.Global }}</td>
            <td>{{ getLocaleDateTime(event.CreatedAt) }}</td>
            <td></td>
          </tr>
          </tbody>
        </table>

        <div v-if="isLoading" class="loading-row">
          <div class="spinner-container">
            <div class="spinner"></div>
            <div>Loading events...</div>
          </div>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-if="!isLoading && events.length === 0 && !error" class="empty-message">
          No events found
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>