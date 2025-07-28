import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

import {createRouter, createWebHistory} from 'vue-router'

import HomeView from './pages/HomeView.vue'
import CreateBucketView from "./pages/CreateBucketView.vue";
import BucketsView from "./pages/BucketsView.vue";
import BucketDetailView from "./pages/BucketDetailView.vue";
import CredentialsView from "./pages/CredentialsView.vue";
import CreateCredentialsView from "./pages/CreateCredentialsView.vue";

const routes = [
    { path: '/', component: HomeView },
    { path: '/buckets', component: BucketsView },
    { path: '/buckets/create-bucket', component: CreateBucketView },
    { path: '/buckets/:slug', component: BucketDetailView },
    { path: '/credentials', component: CredentialsView },
    { path: '/credentials/create-credentials', component: CreateCredentialsView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

createApp(App).use(router).mount('#app')
