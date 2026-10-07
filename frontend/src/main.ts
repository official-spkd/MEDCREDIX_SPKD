// SPKD MedCredix — entry point Vue 3 + Vite.
// "From Hospital Data to Action." — PT Sistem Pelayanan Kesehatan dan Data (SPKD)
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/css/main.css'

createApp(App).use(createPinia()).mount('#app')
