import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import 'bootstrap/dist/css/bootstrap.css'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Notifications from 'vue-notification'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import VueNativeSock from 'vue-native-websocket'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.use(Notifications)
Vue.use(BootstrapVue)
Vue.use(VueNativeSock, 'ws://localhost:8010/ws/', {
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
  reconnectionAttempts: 5, // (Number) number of reconnection attempts before giving up (Infinity),
  reconnectionDelay: 3000 // (Number) how long to initially wait before attempting a new (1000)
})

axios.defaults.baseURL = 'http://localhost:8010/api/v1/go-sb'
axios.defaults.headers.common.Accept = 'application/json'

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
