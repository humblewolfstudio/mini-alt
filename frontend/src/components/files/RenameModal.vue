<script setup lang="ts">
import {ref} from "vue";

const props = defineProps({
  item: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'rename'])

const newName = ref(props.item?.name || '')
const isLoading = ref(false)
const errorMessage = ref('')

const handleRename = async () => {
  if (!newName.value.trim()) {
    errorMessage.value = 'Name cannot be empty'
    return
  }

  if (newName.value === props.item?.name) {
    emit('close')
    return
  }

  let trimmedName = newName.value.trim()

  if (props.item.isFolder && !trimmedName.endsWith('/')) {
    trimmedName += '/'
  }

  if (trimmedName === props.item?.name) {
    emit('close')
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await emit('rename', { newName: trimmedName })
    emit('close')
  } catch (error) {
    errorMessage.value = 'Failed to rename. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Rename {{ item?.isFolder ? 'Folder' : 'File' }}</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>New Name</label>
          <input
              v-model="newName"
              type="text"
              :placeholder="item?.name"
              :disabled="isLoading"
          />
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
            class="rename-btn"
            @click="handleRename"
            :disabled="isLoading"
        >
          <span v-if="isLoading" class="spinner"></span>
          Rename
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rename-btn {
  padding: 8px 16px;
  margin-left: 10px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.rename-btn:hover:not(:disabled) {
  background-color: #2980b9;
}
</style>