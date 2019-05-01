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

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.use(Notifications)
Vue.use(BootstrapVue)

axios.defaults.baseURL = 'http://localhost:8010/api/v1/go-sb'
axios.defaults.headers.common.Accept = 'application/json'

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
