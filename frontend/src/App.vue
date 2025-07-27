<script setup>
import Navbar from "./components/Navbar.vue";
import './assets/tables.css'
import './assets/forms.css'
import './assets/modals.css'

</script>

<template>
  <Navbar/>
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
