<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { onMounted } from "vue";
import Navbar from "./components/Navbar.vue";
import './assets/tables.css'
import './assets/forms.css'
import './assets/modals.css'

const route = useRoute();
const router = useRouter();

const authenticate = async () => {
  try {
    const res = await fetch('/api/users/authenticate')

    if (!res.ok) {
      await router.push('/login')
    }
  } catch (err) {
    await router.push('/login')
  }
}

onMounted(() => {
  authenticate()
})
</script>

<template>
  <Navbar v-if="route.path !== '/login'" />
  <main class="content">
    <RouterView />
  </main>
</template>

<style scoped>
.content {
  flex: 1;
  padding: 0;
  overflow: auto;
  height: 100vh;
}

@media (min-width: 600px) {
  .content {
    padding: 2rem;
  }
}
</style>
