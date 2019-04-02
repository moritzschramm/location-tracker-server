import Vue from 'vue'
import './plugins/vuetify'
import router from './router'

Vue.config.productionTip = false

window.axios = require('axios')
window.qs = require('querystring')
window.axios.defaults.headers.common['Content-Type'] = "application/x-www-form-urlencoded";
window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';

const App = { template: '<router-view></router-view>' }

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
