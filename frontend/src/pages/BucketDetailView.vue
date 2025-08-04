<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {useRoute, useRouter} from 'vue-router'
import Breadcrumb from "../components/files/Breadcrumb.vue";
import FileItem from "../components/files/FileItem.vue";
import UploadModal from "../components/files/UploadModal.vue";
import CreateFolderModal from "../components/files/CreateFolderModal.vue";
import DeleteModal from "../components/files/DeleteModal.vue";
import RenameModal from "../components/files/RenameModal.vue";
import MoveModal from "../components/files/MoveModal.vue";

const router = useRouter()
const route = useRoute()
const bucketName = route.params.slug as string

const currentPath = ref('')
const files = ref<any[]>([])
const isLoading = ref(false)
const showUploadModal = ref(false)
const showFolderModal = ref(false)
const showDeleteModal = ref(false)
const showRenameModal = ref(false)
const showMoveModal = ref(false)
const itemToModify = ref<{ key: string, name: string, isFolder?: boolean } | null>(null)
const errorMessage = ref('')
const successMessage = ref('')

const fetchFiles = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const res = await fetch(`/api/files/list?bucket=${bucketName}&prefix=${currentPath.value}`)

    if(res.status === 401) {
      await router.push('/login')
    }

    const data = await res.json()
    files.value = data.data || []
  } catch (error) {
    console.error('Error fetching files:', error)
    errorMessage.value = 'Failed to load files. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const navigateToFolder = (prefix: string) => {
  currentPath.value = prefix
  fetchFiles()
}

const goUp = () => {
  if (!currentPath.value) return
  const pathParts = currentPath.value.split('/').filter(part => part !== '')
  pathParts.pop()
  currentPath.value = pathParts.length > 0 ? pathParts.join('/') + '/' : ''
  fetchFiles()
}

const downloadFile = async (key: string) => {
  try {
    window.open(`/api/files/download?bucket=${bucketName}&key=${encodeURIComponent(key)}`, '_blank')
  } catch (error) {
    errorMessage.value = 'Failed to download file.'
  }
}

const promptDelete = (item: { key: string, name: string, isFolder?: boolean }) => {
  itemToModify.value = item
  showDeleteModal.value = true
}

const promptRename = (item: { key: string, name: string, isFolder?: boolean }) => {
  itemToModify.value = item
  showRenameModal.value = true
}

const promptMove = (item: { key: string, name: string, isFolder?: boolean }) => {
  itemToModify.value = item
  showMoveModal.value = true
}

const handleDelete = async () => {
  if (!itemToModify.value) return

  try {
    await fetch('/api/files/delete', {
      method: 'POST',
      body: JSON.stringify({
        bucket: bucketName,
        key: itemToModify.value.key
      })
    })
    successMessage.value = `${itemToModify.value.isFolder ? 'Folder' : 'File'} deleted successfully.`
    await fetchFiles()
  } catch (error) {
    errorMessage.value = `Failed to delete ${itemToModify.value.isFolder ? 'folder' : 'file'}.`
  } finally {
    showDeleteModal.value = false
  }
}

const handleRename = async ({ newName }: { newName: string }) => {
  if (!itemToModify.value) return

  const oldKey = itemToModify.value.key
  const pathParts = oldKey.split('/')
  pathParts.pop()
  pathParts.pop()
  const newKey = pathParts.join('/') + (pathParts.length > 0 ? '/' : '') + newName

  try {
    await fetch('/api/files/rename', {
      method: 'PUT',
      body: JSON.stringify({
        bucket: bucketName,
        oldKey,
        newKey
      })
    })
    successMessage.value = 'Renamed successfully.'
    await fetchFiles()
  } catch (error) {
    errorMessage.value = 'Rename failed. Please try again.'
  } finally {
    showRenameModal.value = false
  }
}

const handleMove = async ({ destinationPath }: { destinationPath: string }) => {
  if (!itemToModify.value) return

  try {
    await fetch('/api/files/move', {
      method: 'PUT',
      body: JSON.stringify({
        bucket: bucketName,
        sourceKey: itemToModify.value.key,
        destinationPath
      })
    })
    successMessage.value = 'Moved successfully.'
    await fetchFiles()
  } catch (error) {
    errorMessage.value = 'Move failed. Please try again.'
  } finally {
    showMoveModal.value = false
  }
}

onMounted(() => {
  fetchFiles()
})
</script>

<template>
  <div class="container">
    <div class="header">
      <h1>Bucket: {{ bucketName }}</h1>
      <div class="action-buttons">
        <button class="img-button tooltip" @click="showUploadModal = true">
          <img src="/icons/document-new.svg" width="24" height="24" alt="Upload file">
          <span class="tooltiptext">Upload file</span>
        </button>
        <button class="img-button tooltip" @click="showFolderModal = true">
          <img src="/icons/folder-new.svg" width="24" height="24" alt="New folder">
          <span class="tooltiptext">Create folder</span>
        </button>
      </div>
    </div>

    <Breadcrumb :bucket="bucketName" :path="currentPath" @navigate="navigateToFolder" />

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>

    <div class="table-wrapper">
      <div class="table-content">
        <table class="table">
          <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Last Modified</th>
            <th>Actions</th>
          </tr>
          </thead>
          <tbody>
          <template v-if="isLoading">
            <tr>
              <td colspan="4" class="loading-row">
                <div class="spinner-container">
                  <div class="spinner"></div>
                  <div>Loading files...</div>
                </div>
              </td>
            </tr>
          </template>
          <template v-else>
            <tr v-if="currentPath" @click="goUp" class="clickable-row back-row">
              <td colspan="4">↩ Go up</td>
            </tr>
            <FileItem
                v-for="file in files"
                :key="file.key"
                :item="file"
                @navigate="navigateToFolder"
                @download="downloadFile"
                @delete="promptDelete"
                @rename="promptRename"
                @move="promptMove"
            />
          </template>
          </tbody>
        </table>
      </div>
    </div>

    <UploadModal
        v-if="showUploadModal"
        :bucket="bucketName"
        :current-path="currentPath"
        @close="showUploadModal = false"
        @upload-success="fetchFiles"
    />

    <CreateFolderModal
        v-if="showFolderModal"
        :bucket="bucketName"
        :current-path="currentPath"
        @close="showFolderModal = false"
        @folder-created="fetchFiles"
    />

    <DeleteModal
        v-if="showDeleteModal"
        :item="itemToModify"
        @close="showDeleteModal = false"
        @confirm="handleDelete"
    />

    <RenameModal
        v-if="showRenameModal"
        :item="itemToModify"
        @close="showRenameModal = false"
        @rename="handleRename"
    />

    <MoveModal
        v-if="showMoveModal"
        :bucket="bucketName"
        :item="itemToModify"
        :current-path="currentPath"
        @close="showMoveModal = false"
        @move="handleMove"
    />
  </div>
</template>

<style scoped>
.tooltiptext {
  margin-left: 5px;
}
</style>