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

import Admin from './components/Admin.vue'
import Credits from './components/Credits.vue'
import Dispatch from './components/Dispatch.vue'
import Disable from './components/Disable.vue'
import Edit from './components/Edit.vue'
import Join from './components/Join.vue'
import Log from './components/Log.vue'
import Match from './components/Match.vue'
import New from './components/New.vue'
import NextScreen from './components/NextScreen.vue'
import Participants from './components/Participants.vue'
import Archers from './components/Archers.vue'
import PostMatch from './components/PostMatch.vue'
import Profile from './components/Profile.vue'
import Runnerups from './components/Runnerups.vue'
import ScoreScreen from './components/ScoreScreen.vue'
import Settings from './components/Settings.vue'
import Sidebar from './components/Sidebar.vue'
import TournamentList from './components/TournamentList.vue'
import TournamentView from './components/Tournament.vue'

import DrunkenFallNew from './components/new/DrunkenFall.vue'
import GroupNew from './components/new/Group.vue'

import Person from './models/Person.js'
import Stats from './models/Stats.js'
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
    {
      path: '/',
      name: 'dispatch',
      component: Dispatch
    },
    {
      path: '/facebook/finalize',
      name: 'facebook',
      component: Settings
    },
    {
      path: '/towerfall/tournaments/',
      name: 'start',
      component: TournamentList
    },
    {
      path: '/towerfall/tournaments/new/',
      name: 'new',
      component: New,
    },
    {
      path: '/towerfall/tournaments/new/drunkenfall/',
      name: 'newDrunkenfall',
      component: DrunkenFallNew,
    },
    {
      path: '/towerfall/tournaments/new/group/',
      name: 'newGroup',
      component: GroupNew,
    },
    {
      path: '/towerfall/settings/',
      name: 'settings',
      component: Settings,
    },
    {
      path: '/towerfall/archers/',
      name: 'archers',
      component: Archers,
    },
    {
      path: '/towerfall/archers/:id',
      name: 'profile',
      component: Profile,
    },
    {
      path: '/towerfall/admin',
      name: 'admin',
      component: Admin,
    },
    {
      path: '/towerfall/admin/disable',
      name: 'disable',
      component: Disable,
    },
    {
      path: '/towerfall/tournaments/:tournament/',
      name: 'tournament',
      component: TournamentView
    },
    {
      path: '/towerfall/tournaments/:tournament/join/',
      name: 'join',
      component: Join
    },
    {
      path: '/towerfall/tournaments/:tournament/participants/',
      name: 'participants',
      component: Participants
    },
    {
      path: '/towerfall/tournaments/:tournament/runnerups/',
      name: 'runnerups',
      component: Runnerups
    },
    {
      path: '/towerfall/tournaments/:tournament/edit/',
      name: 'edit',
      component: Edit
    },
    {
      path: '/towerfall/tournaments/:tournament/scores/',
      name: 'scores',
      component: ScoreScreen
    },
    {
      path: '/towerfall/tournaments/:tournament/next/',
      name: 'next',
      component: NextScreen
    },
    {
      path: '/towerfall/tournaments/:tournament/charts/',
      name: 'charts',
      component: PostMatch
    },
    {
      path: '/towerfall/tournaments/:tournament/log/',
      name: 'log',
      component: Log
    },
    {
      path: '/towerfall/tournaments/:tournament/credits/',
      name: 'credits',
      component: Credits
    },
    {
      path: '/towerfall/tournaments/:tournament/:match/',
      name: 'match',
      component: Match
    },
    {
      path: '/towerfall/tournaments/:tournament/*',
      redirect: '/towerfall/tournaments/:tournament/',
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
      router.app.$http.get('/api/towerfall/people/stats/').then(response => {
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

const store = new Vuex.Store({ // eslint-disable-line
  state: {
    tournaments: [],
    user: new Person(),
    userLoaded: false,
    stats: undefined,
    people: undefined,
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
    setStats (state, stats) {
      state.stats = Stats.fromObject(stats)
    },
    setPeople (state, data) {
      state.people = _.map(data, (p) => {
        return Person.fromObject(p)
      })
    },
  },
  getters: {
    getTournament: (state, getters) => (id) => {
      return state.tournaments.find(t => t.id === id)
    },
    upcoming: state => {
      return _.reverse(_.filter(state.tournaments, 'isUpcoming'))
    },
    getPerson: (state, getters) => (id) => {
      if (!state.people) {
        return undefined
      }
      return state.people.find(p => p.id === id)
    },
  }
})

var Root = Vue.extend(App)
new Root({ // eslint-disable-line
  router: router,
  store: store,
}).$mount("#app")
