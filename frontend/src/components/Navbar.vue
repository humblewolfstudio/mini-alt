<script setup lang="ts">
import {computed, ref} from "vue";
import {useRoute, useRouter} from 'vue-router';
import Cookies from "js-cookie";

const isCollapsed = ref(false)
const isAdmin = ref(Cookies.get("admin") === "true");
const route = useRoute();
const router = useRouter()

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

const adminRoutes = [
  {path: '/users', name: 'Users', icon: 'users'},
  {path: '/events', name: 'Events', icon: 'events'},
  {path: '/metrics', name: 'Metrics', icon: 'metrics'}
];

const userRoutes = [
  {path: '/', name: 'Home', icon: 'home'},
  {path: '/buckets', name: 'Buckets', icon: 'buckets'},
  {path: '/credentials', name: 'Credentials', icon: 'credentials'},
];

const filteredAdminRoutes = computed(() => {
  return isAdmin.value ? adminRoutes : [];
});

const filteredUserRoutes = computed(() => {
  return userRoutes;
});

const isActive = (navItem: any) => {
  if (navItem.path.includes(':')) {
    const regex = new RegExp('^' + navItem.path.replace(/:[^/]+/g, '[^/]+') + '($|/)');
    return regex.test(route.path);
  }
  return route.path === navItem.path || route.path.startsWith(navItem.path + '/');
};

const logout = async () => {
  try {
    const res = await fetch('/api/users/logout')

    if (res.ok) {
      await router.push('/login')
    }
  } catch (err) {
    console.log("Error logging out:", err)
  }
}
</script>

<template>
  <div class="navbar-container" :class="{ collapsed: isCollapsed }">
    <button class="toggle-btn" @click="toggleCollapse">
      {{ isCollapsed ? '☰' : '×' }}
    </button>

    <div class="navbar-content">
      <nav class="nav-links">
        <div class="nav-section-title">User</div>

        <RouterLink
            v-for="route in filteredUserRoutes"
            :key="route.path"
            :to="route.path"
            class="nav-link"
            :class="{ 'router-link-exact-active': isActive(route) }"
        >
          <div class="icon-container">
            <img class="nav-icon" :src="'/icons/' + route.icon + '.svg'" width="25" height="25" :alt="route.name"/>
          </div>
          <div class="text-container">
            <span class="nav-text">{{ route.name }}</span>
          </div>
        </RouterLink>

        <template v-if="filteredAdminRoutes.length">
          <div class="nav-section-title administrator-section">Administration</div>

          <RouterLink
              v-for="route in filteredAdminRoutes"
              :key="route.path"
              :to="route.path"
              class="nav-link"
              :class="{ 'router-link-exact-active': isActive(route) }"
          >
            <div class="icon-container">
              <img class="nav-icon" :src="'/icons/' + route.icon + '.svg'" width="25" height="25" :alt="route.name"/>
            </div>
            <div class="text-container">
              <span class="nav-text">{{ route.name }}</span>
            </div>
          </RouterLink>
        </template>
      </nav>

      <button @click="logout" class="nav-link logout-btn">
        <div class="icon-container">
          <img class="nav-icon" src="/icons/logout.svg" width="25" height="25" alt="Logout"/>
        </div>
        <div class="text-container">
          <span class="nav-text">Logout</span>
        </div>
      </button>
    </div>
  </div>
</template>

<style scoped>
.navbar-container {
  --navbar-width: 200px;
  --navbar-collapsed-width: 60px;
  --transition-duration: 0.3s;

  width: var(--navbar-width);
  height: 100vh;
  background-color: #2c3e50;
  color: white;
  transition: width var(--transition-duration) ease;
  flex-shrink: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.navbar-container.collapsed {
  width: var(--navbar-collapsed-width);
}

.toggle-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0.5rem;
  text-align: center;
  width: 100%;
}

.toggle-btn:hover {
  background-color: #34495e;
}

.navbar-content {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-grow: 1;
  padding: 0 8px;
  overflow: hidden;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-link {
  color: white;
  text-decoration: none;
  display: flex;
  align-items: center;
  border-radius: 4px;
  padding: 8px;
  background-color: transparent;
  border: none;
  transition: background-color 0.2s;
  overflow: hidden;
  text-align: left;
  cursor: pointer;
}

.nav-section-title {
  font-size: 0.75rem;
  font-weight: bold;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.6);
  padding: 6px 8px;
  letter-spacing: 1px;
}

.administrator-section {
  margin-top: 50px;
}

.icon-container {
  width: 25px;
  height: 25px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
}

.text-container {
  overflow: hidden;
  margin-left: 12px;
  transition: opacity var(--transition-duration) ease,
  margin var(--transition-duration) ease;
}

.nav-text {
  white-space: nowrap;
  opacity: 1;
}

.navbar-container.collapsed .text-container {
  opacity: 0;
  margin-left: 0;
  width: 0;
}

.nav-link:hover {
  background-color: #34495e;
}

.nav-link.router-link-exact-active {
  background-color: #42b983;
  font-weight: bold;
}

.logout-btn {
  margin-bottom: 10px;
}
</style>