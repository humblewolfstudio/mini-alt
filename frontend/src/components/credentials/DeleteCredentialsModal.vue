<script setup lang="ts">
import {ref} from "vue";

defineProps({
  accessKey: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'confirm'])

const isLoading = ref(false)

const handleDelete = async () => {
  isLoading.value = true
  try {
    await emit('confirm')
    emit('close')
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Confirm Delete</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
        <p>Are you sure you want to delete <strong>{{ accessKey }}</strong>?</p>
        <p class="warning">This action cannot be undone.</p>
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
            class="delete-btn"
            @click="handleDelete"
            :disabled="isLoading"
        >
          <span v-if="isLoading" class="spinner"></span>
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.warning {
  color: #e74c3c;
  font-weight: 500;
}

.delete-btn {
  padding: 8px 16px;
  margin-left: 10px;
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.delete-btn:hover:not(:disabled) {
  background-color: #c0392b;
}

.delete-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
</style>