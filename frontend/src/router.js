import Vue from 'vue'
import axios from 'axios'
import Router from 'vue-router'
import Login from './views/Login.vue'

import Dashboard from './views/Dashboard.vue'
import Locations from './components/Locations.vue'
import Settings from './components/Settings.vue'
import Battery from './components/Battery.vue'

Vue.use(Router)

const router = new Router({
  routes: [
    { path: '/', component: Login },
    { path: '/dashboard', component: Dashboard,
      children: [
        { path: 'locations', component: Locations },
        { path: 'settings', component: Settings },
        { path: 'battery', component: Battery },
      ] 
    }
  ]
})

axios.post('/api/refresh')
.then(response => {

  if(response.status != 200) {

    router.push({ path: '/' })
  }

})
.catch(() => {

  router.push({ path: '/' })
})

router.beforeEach((to, from, next) => {

  const path = to.path
  const loggedIn = localStorage.getItem('logged_in')

  if(path == '/' && loggedIn) {      // user is logged in, show dashboard

    next({ path: '/dashboard/locations' })

  } else if(path != '/' && !loggedIn) {  // user is logged out, show login form

    next({ path: '/' })

  } else {

    next()
  }
})

export default router
