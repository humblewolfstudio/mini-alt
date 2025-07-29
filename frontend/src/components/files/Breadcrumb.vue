<script setup lang="ts">
import {computed} from "vue";

const props = defineProps({
  bucket: {
    type: String,
    required: true
  },
  path: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['navigate'])

const navigate = (prefix: string) => {
  emit('navigate', prefix)
}

const pathParts = computed(() => {
  if (!props.path) return []

  const parts = props.path.split('/').filter(part => part !== '')
  const result = []
  let currentPath = ''

  for (const part of parts) {
    currentPath += part + '/'
    result.push({
      name: part,
      path: currentPath
    })
  }

  return result
})
</script>

<template>
  <div class="breadcrumb">
    <span class="breadcrumb-item" @click="navigate('')">
      /{{ bucket }}/
    </span>
    <span
        v-for="(part, index) in pathParts"
        :key="index"
        class="breadcrumb-item"
        @click="navigate(part.path)"
    >
      {{ part.name }}/
    </span>
  </div>
</template>

<style scoped>
.breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 20px;
  font-size: 0.9rem;
}

.breadcrumb-item {
  cursor: pointer;
  color: #42b983;
  padding: 2px 0;
  border-radius: 3px;
}

.breadcrumb-item:hover {
  background-color: #f0fdf4;
  text-decoration: underline;
}

.breadcrumb-item:first-of-type {
  padding-left: 4px;
}
</style>