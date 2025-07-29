<script setup lang="ts">
import {ref} from "vue";

const props = defineProps({
  bucket: {
    type: String,
    required: true
  },
  currentPath: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'folder-created'])

const folderName = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

const createFolder = async () => {
  if (!folderName.value.trim()) {
    errorMessage.value = 'Folder name cannot be empty'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await fetch('/api/files/create-folder', {
      method: 'POST',
      body: JSON.stringify({
        bucket: props.bucket,
        prefix: props.currentPath,
        folderName: folderName.value.trim()
      })
    })

    emit('folder-created')
    emit('close')
  } catch (error) {
    errorMessage.value = 'Failed to create folder. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Create New Folder</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label for="folderName">Folder Name</label>
          <input
              id="folderName"
              v-model="folderName"
              type="text"
              placeholder="Enter folder name"
              :disabled="isLoading"
          />
          <p class="hint">Folder will be created in: {{ currentPath || 'root' }}</p>
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
      </div>

      <div class="modal-footer">
        <button
            class="cancel-btn"
            @click="emit('close')"
            :disabled="isLoading"
        >
          Cancel
        </button>
        <button
            class="create-btn"
            @click="createFolder"
            :disabled="isLoading"
        >
          <span v-if="isLoading" class="spinner"></span>
          Create Folder
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.create-btn {
  padding: 8px 16px;
  margin-left: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.create-btn:hover:not(:disabled) {
  background-color: #3aa876;
}

.create-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>