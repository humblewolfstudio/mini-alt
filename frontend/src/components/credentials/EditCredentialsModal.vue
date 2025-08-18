<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
  accessKey: {
    type: String,
    default: ''
  },
  data: {
    type: Object
  }
})

const formatDate = (iso: string | undefined): string => {
  if (!iso) return ''
  return iso.split('T')[0]
}

const emit = defineEmits(['continue', 'close'])

const name = ref(props.data.Name ?? '')
const description = ref(props.data.Description ?? '')
const status = ref(props.data.Status ?? true)
const expiresAt = ref(formatDate(props.data.ExpiresAt))

const isLoading = ref(false)

const handleContinue = () => {
  emit('continue', {
    name: name.value,
    description: description.value,
    status: status.value,
    expiresAt: expiresAt.value ? new Date(expiresAt.value).toISOString() : null
  })
}
</script>


<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Edit Credentials</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
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
        <div class="form-group">
          <label for="statusToggle">Active Status</label>
          <label class="switch">
            <input type="checkbox" id="statusToggle" v-model="status" :disabled="isLoading" />
            <span class="slider"></span>
          </label>
          <p class="hint">Toggle to activate or deactivate these credentials.</p>
        </div>
      </div>

      <div class="modal-footer">
        <button class="action-btn" @click="handleContinue">
          Continue
        </button>
      </div>
    </div>
  </div>
</template>


<style scoped>
.modal-body {
  padding: 1rem;
}
</style>

