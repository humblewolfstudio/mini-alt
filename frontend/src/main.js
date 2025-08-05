import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

import {createRouter, createWebHistory} from 'vue-router'

import HomeView from './pages/HomeView.vue'
import BucketsCreateView from "./pages/BucketsCreateView.vue";
import BucketsView from "./pages/BucketsView.vue";
import BucketDetailView from "./pages/BucketDetailView.vue";
import CredentialsView from "./pages/CredentialsView.vue";
import CredentialsCreateView from "./pages/CredentialsCreateView.vue";
import UsersView from "./pages/UsersView.vue";
import UsersCreateView from "./pages/UsersCreateView.vue";
import LoginView from "./pages/LoginView.vue";
import EventsView from "./pages/EventsView.vue";
import EventsCreateView from "./pages/EventsCreateView.vue";

const routes = [
    { path: '/', component: HomeView },
    { path: '/buckets', component: BucketsView },
    { path: '/buckets/create-bucket', component: BucketsCreateView },
    { path: '/buckets/:slug', component: BucketDetailView },
    { path: '/credentials', component: CredentialsView },
    { path: '/credentials/create-credentials', component: CredentialsCreateView },
    { path: '/users', component: UsersView },
    { path: '/users/create-users', component: UsersCreateView },
    { path: '/login', component: LoginView },
    { path: '/events', component: EventsView },
    { path: '/events/create-event', component: EventsCreateView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

createApp(App).use(router).mount('#app')

export default router