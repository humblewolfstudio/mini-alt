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

const routes = [
    { path: '/', component: HomeView },
    { path: '/buckets', component: BucketsView },
    { path: '/buckets/create-bucket', component: BucketsCreateView },
    { path: '/buckets/:slug', component: BucketDetailView },
    { path: '/credentials', component: CredentialsView },
    { path: '/credentials/create-credentials', component: CredentialsCreateView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

createApp(App).use(router).mount('#app')
