<template>
  
<v-app>
  
  <v-navigation-drawer fixed app dark v-model="drawer" class="primary">
    
    <v-list>

      <v-list-tile 
        v-for="item in linksEditor" 
        router :to="item.path" :key="item.title" :value="item.active" 
        active-class="accent--text" ripple>
        <v-list-tile-action>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>{{ item.title }}</v-list-tile-title>
        </v-list-tile-content>
      </v-list-tile>

    </v-list>

  </v-navigation-drawer>

  <v-toolbar absolute app light color="white" elevation-1>
    <v-toolbar-side-icon @click.native="drawer = !drawer"></v-toolbar-side-icon>
    <span class="title ml-3 mr-5"><span class="font-weight-light">Location</span> Tracker</span>

    <v-spacer></v-spacer>

    <v-toolbar-items>
      <v-btn flat large slot="activator" @click="logout">Logout</v-btn>
    </v-toolbar-items>
  </v-toolbar>

  <v-dialog
    v-model="loading"
    hide-overlay
    persistent
    width="300">
    <v-card
    color="primary"
    dark>
    <v-card-text>
      Laden...
      <v-progress-linear
      indeterminate
      color="white"
      class="mb-0"
      ></v-progress-linear>
    </v-card-text>
    </v-card>
  </v-dialog>

  <v-content>
    <v-container fluid>
        <router-view></router-view>
      <transition name="fade" mode="out-in">
      </transition>
    </v-container>
  </v-content>

</v-app>    

</template>

<script>
export default {
  data () {
    return {
      drawer: null,
      loading: false,
      linksEditor: [
        { path: '/dashboard/locations', icon: 'location_on', title: 'Locations'},
        { path: '/dashboard/settings', icon: 'settings', title: 'Settings'},
        { path: '/dashboard/battery', icon: 'battery_charging_full', title: 'Battery'}, 
      ],
    }
  },
  methods: {
    link(path) {
      this.$router.push({ path: path })
    },
    logout() {

      localStorage.removeItem('logged_in')
      this.$router.push({ path: '/' })

      axios.post('/api/logout')
      .then(() => {

      })
    }
  }
}

</script>

<style>

.fade-enter-active, .fade-leave-active {
  transition: .15s;
}
.fade-enter, .fade-leave-to {
  opacity: 0.2;
}

</style>