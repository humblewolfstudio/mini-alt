<script setup lang="ts">
import {onMounted, ref, watch} from "vue";

const props = defineProps({
  bucket: {
    type: String,
    required: true
  },
  item: {
    type: Object,
    default: null
  },
  currentPath: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'move'])

const destinationPath = ref('')
const isLoading = ref(false)
const errorMessage = ref('')
const folders = ref<string[]>([])

const fetchFolders = async () => {
  try {
    const params = new URLSearchParams({
      bucket: props.bucket,
      excludePrefix: props.item?.key || '',
      currentPath: props.currentPath // Pass current path
    })

    const res = await fetch(`/api/files/list-folders?${params}`)
    const data = await res.json()
    folders.value = data.data || []
  } catch (error) {
    console.error('Error fetching folders:', error)
  }
}

watch(() => props.currentPath, () => {
  fetchFolders()
})

const handleMove = async () => {
  if (!destinationPath.value.trim()) {
    errorMessage.value = 'Please select a destination'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await emit('move', { destinationPath: destinationPath.value.trim() })
    emit('close')
  } catch (error) {
    errorMessage.value = 'Failed to move. Please try again.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchFolders()
})
</script>

<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Move {{ item?.isFolder ? 'Folder' : 'File' }}</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>Current Location</label>
          <input
              type="text"
              :value="item?.key || currentPath"
              disabled
          />
        </div>

        <div class="form-group">
          <label>Destination Path</label>
          <select v-model="destinationPath" :disabled="isLoading">
            <option value="">Select destination...</option>
            <option v-for="folder in folders" :key="folder" :value="folder">
              {{ folder === '' ? '/' : '/' + folder }}
            </option>
          </select>
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
            class="move-btn"
            @click="handleMove"
            :disabled="isLoading || !destinationPath"
        >
          <span v-if="isLoading" class="spinner"></span>
          Move
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.move-btn {
  padding: 8px 16px;
  margin-left: 10px;
  background-color: #9b59b6;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.move-btn:hover:not(:disabled) {
  background-color: #8e44ad;
}

select {
  padding: 12px 15px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  font-size: 1rem;
  width: 100%;
  background-color: white;
}

select:focus {
  outline: none;
  border-color: #42b983;
  box-shadow: 0 0 0 2px rgba(66, 185, 131, 0.2);
}
</style>