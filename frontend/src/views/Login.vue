<template>

<v-app>
    <v-content>
      <v-container fluid fill-height>

        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
              <v-toolbar dark color="primary">
                <v-toolbar-title>Login</v-toolbar-title>
                <v-spacer></v-spacer>
              </v-toolbar>
              <v-card-text>
                <v-form ref="form" lazy-validation>
                  <v-text-field ref="username_field" id="username" name="username" prepend-icon="person" label="Username" type="text" 
                                v-model="username"
                                :rules="[v => !!v || 'Username required']"></v-text-field>
                  <v-text-field id="password" name="password" prepend-icon="lock" label="Password" type="password"
                                v-model="password"
                                @keyup.enter="login"
                                :error-messages="errors"
                                :rules="[v => !!v || 'Password required']"></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="login" :loading="loading">Login</v-btn>
              </v-card-actions>
              <v-progress-linear v-if="loading" indeterminate></v-progress-linear>
            </v-card>
          </v-flex>
        </v-layout>

      </v-container>
    </v-content>
</v-app>  

</template>

<script>
export default {
  name: 'login',
  mounted() {
    this.$refs.username_field.focus()
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      if(from.path != '/') vm.previousPath = from.path
    })
  },
  data () {
    return {
      loading: false,
      errors: [],
      username: '',
      password: '',
      previousPath: null,
    }
  },
  props: {
    source: String
  },
  methods: {
    login() {

      if(this.$refs.form.validate()) {

        this.errors = []
        this.loading = true

        window.axios.post('/api/login', window.qs.stringify({
          'uuid': this.username,
          'password': this.password
        }))
        .then(() => {

          this.loading = false

          // token automatically set by server in cookie
          localStorage.setItem('logged_in', true)
          localStorage.setItem('uuid', this.username)

          this.$refs.form.reset()
          if(this.previousPath) {
            this.$router.push({ path: this.previousPath })
          } else {
            this.$router.push({ path: '/dashboard/locations' })
          }

        })
        .catch(() => {

          this.loading = false

          this.password = ''
          this.errors.push('Wrong username or password')
        })
      }
    },
  }
}
</script>