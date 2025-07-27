<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue";

const props = defineProps({
  item: {
    type: Object,
    required: true
  }
});

const emit = defineEmits([
  'navigate',
  'download',
  'preview',
  'delete',
  'rename',
  'move'
]);

const isMenuOpen = ref(false);
const menuPosition = ref({ top: 0, left: 0 });
const triggerRef = ref<HTMLElement | null>(null);

const handleAction = (action: string) => {
  isMenuOpen.value = false;
  switch (action) {
    case 'download':
      emit('download', props.item.key);
      break;
    case 'delete':
      emit('delete', props.item);
      break;
    case 'rename':
      emit('rename', props.item);
      break;
    case 'move':
      emit('move', props.item);
      break;
  }
};

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString();
};

const toggleMenu = async (e: MouseEvent) => {
  isMenuOpen.value = !isMenuOpen.value;

  if (isMenuOpen.value && triggerRef.value) {
    await nextTick();
    const rect = triggerRef.value.getBoundingClientRect();
    menuPosition.value = {
      top: rect.bottom + window.scrollY,
      left: rect.right + window.scrollX
    };
  }
};

const handleClickOutside = (event: MouseEvent) => {
  if (!triggerRef.value?.contains(event.target as Node)) {
    isMenuOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", handleClickOutside);
});
</script>

<template>
  <tr @click="item.isFolder ? emit('navigate', item.key) : null" :class="{ 'clickable-cell': item.isFolder, 'folder-row': item.isFolder }">
    <td class="clickable-name">
      <span v-if="item.isFolder">📁</span>
      <span v-else>📄</span>
      {{ item.name }}
    </td>
    <td>{{ item.isFolder ? '-' : formatSize(item.size) }}</td>
    <td>{{ item.isFolder ? '-' : formatDate(item.lastModified) }}</td>
    <td>
      <div class="file-actions">
        <button v-if="!item.isFolder" class="action-btn" @click="handleAction('download')" title="Download">
          <img src="/icons/download.svg" width="32" height="32" alt="Download">
        </button>
        <div class="dropdown">
          <button
              class="action-btn"
              ref="triggerRef"
              @click.stop="toggleMenu"
              title="More actions"
          >
            <img src="/icons/dots-vertical.svg" width="32" height="32" alt="More">
          </button>
        </div>
      </div>
    </td>
  </tr>

  <Teleport to="body">
    <div
        v-if="isMenuOpen"
        class="dropdown-menu"
        :style="{
      position: 'absolute',
      top: menuPosition.top + 'px',
      left: menuPosition.left + 'px'
    }"
    >
      <button @click="handleAction('rename')">Rename</button>
      <button @click="handleAction('move')">Move</button>
      <button @click="handleAction('delete')">Delete</button>
    </div>
  </Teleport>
</template>

<style scoped>
.folder-row {
  font-weight: 500;
}

.clickable-cell {
  cursor: pointer;
}
.clickable-name:hover {
  text-decoration: underline;
}

.file-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  background-color: #f0f0f0;
}

.dropdown-menu {
  min-width: 140px;
  max-width: 220px;
  background-color: white;
  box-shadow: 0 2px 5px rgba(0,0,0,0.2);
  border-radius: 4px;
  z-index: 10000;
  padding: 4px 0;
  transform-origin: top right;
  transform: translateX(-100%);
}

.dropdown-menu button {
  width: 100%;
  text-align: left;
  font-size: 1em;
  padding: 8px 12px;
  border: none;
  background: none;
  cursor: pointer;
}

.dropdown-menu button:hover {
  background-color: #f5f5f5;
}
</style>
