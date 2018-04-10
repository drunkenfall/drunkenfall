import _ from 'lodash'
import Person from '../models/Person.js'

const middleware = (router) =>
router.beforeEach((to, from, next) => {
  window.scrollTo(0, 0)

  // Why in the world does this need a short timeout? The connect
  // isn't set otherwise.
  setTimeout(function () {
    if (!router.app.$store.state.tournamentsLoaded) {
      router.app.$http.get('/api/tournaments/').then(response => {
        let data = JSON.parse(response.data)
        router.app.$store.commit('updateAll', {
          "tournaments": data.tournaments,
          "$vue": router.app,
        })
      }, response => {
        console.log("Failed getting tournaments", response)
      })
    }

    if (!router.app.$store.state.user.authenticated) {
      router.app.$http.get('/api/user/').then(response => {
        let data = JSON.parse(response.data)

        // If we're not signed in, then the backend will return an
        // object with just "false" and nothing else. If this happens,
        // we should just skip this.
        if (!data || data.authenticated === false) {
          // Mark that we tried to load the user. This means that the
          // interface will show the Facebook login button.
          router.app.$store.commit("setUserLoaded", true)
          return
        }

        router.app.$store.commit('setUser', Person.fromObject(
          data,
          router.app.$cookie
        ))
      }, response => {
        console.log("Failed getting user data", response)
      })
    }

    if (!router.app.$store.state.stats) {
      router.app.$http.get('/api/people/stats/').then(response => {
        let data = JSON.parse(response.data)
        router.app.$store.commit('setStats', data)
        // Since the stats also contain the profiles, we can use this
        // data to populate those as well!
        router.app.$store.commit('setPeople', _.map(data, (s) => {
          return s.person
        }))
      }, response => {
        console.log("Failed getting stats", response)
      })
    }
  }, 20)

  // Reset any pulsating lights
  document.getElementsByTagName("body")[0].className = ""
  next()
})

export default middleware
