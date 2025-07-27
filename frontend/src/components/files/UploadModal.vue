<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps({
  bucket: {
    type: String,
    required: true
  },
  currentPath: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['close', 'upload-success']);

const fileInput = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const handleFileChange = async (event: Event) => {
  const files = (event.target as HTMLInputElement).files;
  if (!files || files.length === 0) return;

  if (!files || files.length === 0) return;

  isUploading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    const formData = new FormData();
    for (let i = 0; i < files.length; i++) {
      const safeName = files[i].name.replace(/\s+/g, '_');
      formData.append('files', files[i], safeName);
    }
    formData.append('bucket', props.bucket);
    formData.append('prefix', props.currentPath);

    const response = await fetch('/api/files/upload', {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error('Upload failed');
    }

    successMessage.value = 'Files uploaded successfully!';
    emit('upload-success');
  } catch (error) {
    console.error('Upload error:', error);
    errorMessage.value = 'Failed to upload files. Please try again.';
  } finally {
    isUploading.value = false;
  }
};

const triggerFileInput = () => {
  fileInput.value?.click();
};
</script>

<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Upload Files</h2>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>

      <div class="modal-body">
        <div v-if="isUploading" class="uploading-state">
          <div class="spinner"></div>
          <p>Uploading files...</p>
        </div>

        <div v-else class="upload-area" @click="triggerFileInput">
          <input
              ref="fileInput"
              type="file"
              multiple
              @change="handleFileChange"
              style="display: none"
          />
          <img src="/icons/upload.svg" width="48" height="48" alt="Upload">
          <p>Click to select files or drag and drop</p>
          <p class="hint">Files will be uploaded to: {{ currentPath || 'root' }}</p>
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <div v-if="successMessage" class="success-message">
          {{ successMessage }}
        </div>
      </div>

      <div class="modal-footer">
        <button class="cancel-btn" @click="emit('close')">Close</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.upload-area {
  border: 2px dashed #e0e0e0;
  border-radius: 8px;
  padding: 30px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #42b983;
  background-color: #f0fdf4;
}

.upload-area p {
  margin: 10px 0 0;
  color: #34495e;
}

.hint {
  font-size: 0.85rem;
  color: #7f8c8d;
  margin-top: 5px;
}

.uploading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 30px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #42b983;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>