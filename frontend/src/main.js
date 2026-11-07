import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import Chart from 'chart.js/auto'

Chart.defaults.color = '#587067'
Chart.defaults.borderColor = '#d7e5df'

createApp(App)
.use(router)
.mount('#app')