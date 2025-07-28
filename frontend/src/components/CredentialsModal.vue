<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
  accessKey: {
    type: String,
    default: ''
  },
  secretKey: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['continue'])

const copiedField = ref<string | null>(null)

function copyToClipboard(value: string, field: string) {
  navigator.clipboard.writeText(value).then(() => {
    copiedField.value = field
    setTimeout(() => {
      copiedField.value = null
    }, 2000)
  })
}
</script>


<template>
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>New Credentials</h2>
        <button class="close-btn" @click="emit('continue')">×</button>
      </div>

      <div class="modal-body">
        <div class="credentials">
          <div class="credential-item">
            <label>Access Key:</label>
            <div class="credential-box">
              <pre class="credential">{{ props.accessKey }}</pre>
              <div class="copy-actions">
                <button class="copy-btn" @click="copyToClipboard(props.accessKey, 'access')">{{copiedField === 'access' ? 'Copied!' : 'Copy'}}</button>
              </div>
            </div>
          </div>

          <div class="credential-item">
            <label>Secret Key:</label>
            <div class="credential-box">
              <pre class="credential">{{ props.secretKey }}</pre>
              <div class="copy-actions">
                <button class="copy-btn" @click="copyToClipboard(props.secretKey, 'secret')">{{copiedField === 'secret' ? 'Copied!' : 'Copy'}}</button>
              </div>
            </div>
            <p class="warning">You will only be able to see this secret key once. Please copy it now.</p>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="action-btn" @click="emit('continue')">
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

.credentials {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.credential-item {
  display: flex;
  flex-direction: column;
}

.credential-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #f5f5f5;
  padding: 0.5rem;
  border-radius: 5px;
  font-family: monospace;
  word-break: break-all;
}

.credential {
  margin: 0;
  flex: 1;
}

.copy-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.copy-btn {
  background: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.copy-btn:hover {
  background: #3aa876;
}

.copied-msg {
  font-size: 0.75rem;
  color: #16a34a;
}

.warning {
  color: #d97706;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}
</style>

