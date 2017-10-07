/* eslint-env browser */

import Vue from 'vue'
import Vuex from 'vuex'
import Router from 'vue-router'
import Resource from 'vue-resource'
import Cookie from 'vue-cookie'
import Vue2Filters from 'vue2-filters'
import Toast from 'vue-easy-toast'
import * as Icon from 'vue-awesome'
import App from './App'

import _ from 'lodash'

import TournamentList from './components/TournamentList.vue'
import TournamentView from './components/Tournament.vue'
import New from './components/New.vue'
import Join from './components/Join.vue'
import Participants from './components/Participants.vue'
import Runnerups from './components/Runnerups.vue'
import Edit from './components/Edit.vue'
import Match from './components/Match.vue'
import PostMatch from './components/PostMatch.vue'
import Log from './components/Log.vue'
import ScoreScreen from './components/ScoreScreen.vue'
import NextScreen from './components/NextScreen.vue'
import Facebook from './components/Facebook.vue'
import FacebookFinalize from './components/FacebookFinalize.vue'
import Credits from './components/Credits.vue'
import Sidebar from './components/Sidebar.vue'

import Person from './models/Person.js'
import {Credits as CreditsModel} from './models/Credits.js'
import Tournament from './models/Tournament.js'

// install packages
Vue.use(Vuex)
Vue.use(Router)
Vue.use(Resource)
Vue.use(Cookie)
Vue.use(Toast)
Vue.use(Vue2Filters)

Vue.component('icon', Icon)
Vue.component('sidebar', Sidebar)

// routing
var router = new Router({
  mode: 'history',
  routes: [
    // As long as we only have Drunken TowerFall on drunkenfall.com, we should
    // always redirect to the towerfall app right away.
    {
      path: '/',
      redirect: '/towerfall/',
    },
    {
      path: '/facebook/',
      component: Facebook,
      name: 'facebook',
    },
    {
      path: '/facebook/finalize',
      component: FacebookFinalize
    },
    {
      path: '/towerfall/',
      name: 'start',
      component: TournamentList
    },
    {
      path: '/towerfall/new/',
      name: 'new',
      component: New,
    },
    {
      path: '/towerfall/:tournament/',
      name: 'tournament',
      component: TournamentView
    },
    {
      path: '/towerfall/:tournament/join/',
      name: 'join',
      component: Join
    },
    {
      path: '/towerfall/:tournament/participants/',
      name: 'participants',
      component: Participants
    },
    {
      path: '/towerfall/:tournament/runnerups/',
      name: 'runnerups',
      component: Runnerups
    },
    {
      path: '/towerfall/:tournament/edit/',
      name: 'edit',
      component: Edit
    },
    {
      path: '/towerfall/:tournament/scores/',
      name: 'scores',
      component: ScoreScreen
    },
    {
      path: '/towerfall/:tournament/next/',
      name: 'next',
      component: NextScreen
    },
    {
      path: '/towerfall/:tournament/charts/',
      name: 'charts',
      component: PostMatch
    },
    {
      path: '/towerfall/:tournament/log/',
      name: 'log',
      component: Log
    },
    {
      path: '/towerfall/:tournament/credits/',
      name: 'credits',
      component: Credits
    },
    {
      path: '/towerfall/:tournament/:match/',
      name: 'match',
      component: Match
    },
    {
      path: '/towerfall/:tournament/*',
      redirect: '/towerfall/:tournament/',
    },
  ],
})

router.beforeEach((to, from, next) => {
  window.scrollTo(0, 0)

  // Why in the world does this need a short timeout? The connect
  // isn't set otherwise.
  setTimeout(function () {
    router.app.connect()

    if (!router.app.$store.state.user.authenticated) {
      router.app.$http.get('/api/towerfall/user/').then(response => {
        let data = JSON.parse(response.data)

        // If we're not signed in, then the backend will return an
        // object with just "false" and nothing else. If this happens,
        // we should just skip this.
        if (data.authenticated === false) {
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
  }, 20)

  // Reset any pulsating lights
  document.getElementsByTagName("body")[0].className = ""
  next()
})

const store = new Vuex.Store({ // eslint-disable-line
  state: {
    tournaments: [],
    user: new Person(),
    userLoaded: false,
    credits: {}
  },
  mutations: {
    updateAll (state, data) {
      state.tournaments = _.reverse(_.map(data.tournaments, (t) => {
        return Tournament.fromObject(t, data.$vue)
      }))
    },
    setUser (state, user) {
      state.user = user
      state.userLoaded = true
    },
    setUserLoaded (state, loaded) {
      state.userLoaded = loaded
    },
    setCredits (state, credits) {
      state.credits = CreditsModel.fromObject(credits)
    },
  },
  getters: {
    getTournament: (state, getters) => (id) => {
      return state.tournaments.find(t => t.id === id)
    }
  }
})

var Root = Vue.extend(App)
new Root({ // eslint-disable-line
  router: router,
  store: store,
}).$mount("#app")
