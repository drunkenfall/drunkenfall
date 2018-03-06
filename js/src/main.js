/* eslint-env browser */

import Vue from 'vue'
import Vuex from 'vuex'
import Router from 'vue-router'
import Resource from 'vue-resource'
import Cookie from 'vue-cookie'
import Vue2Filters from 'vue2-filters'
import Toast from 'vue-easy-toast'
import vueHeadful from 'vue-headful'
import * as Icon from 'vue-awesome'
import App from './App'

import _ from 'lodash'

import About from './components/About.vue'
import Admin from './components/Admin.vue'
import Credits from './components/Credits.vue'
// import Dispatch from './components/Dispatch.vue'
import Disable from './components/Disable.vue'
import Edit from './components/Edit.vue'
import Join from './components/Join.vue'
import JudgeInterface from './components/Judge.vue'
import Log from './components/Log.vue'
import Match from './components/Match.vue'
import New from './components/New.vue'
import NextScreen from './components/NextScreen.vue'
import Participants from './components/Participants.vue'
import Archers from './components/Archers.vue'
import PostMatch from './components/PostMatch.vue'
import Runnerups from './components/Runnerups.vue'
import ScoreScreen from './components/ScoreScreen.vue'
import Settings from './components/Settings.vue'
import Sidebar from './components/Sidebar.vue'
import Stream from './components/Stream'
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
Vue.component('headful', vueHeadful)

// routing
var router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'about',
      component: About
    },
    {
      path: '/facebook/finalize',
      name: 'facebook',
      component: Settings
    },
    {
      path: '/tournaments/',
      name: 'tournaments',
      component: TournamentList
    },
    {
      path: '/tournaments/new/',
      name: 'new',
      component: New,
    },
    {
      path: '/tournaments/new/drunkenfall/',
      name: 'newDrunkenfall',
      component: DrunkenFallNew,
    },
    {
      path: '/tournaments/new/group/',
      name: 'newGroup',
      component: GroupNew,
    },
    {
      path: '/settings/',
      name: 'settings',
      component: Settings,
    },
    {
      path: '/archers/',
      name: 'archers',
      component: Archers,
    },
    {
      path: '/archers/:id',
      name: 'archer',
      component: Archers,
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin,
    },
    {
      path: '/admin/disable',
      name: 'disable',
      component: Disable,
    },
    {
      path: '/tournaments/:tournament/',
      name: 'tournament',
      component: TournamentView
    },
    {
      path: '/tournaments/:tournament/join/',
      name: 'join',
      component: Join
    },
    {
      path: '/tournaments/:tournament/participants/',
      name: 'participants',
      component: Participants
    },
    {
      path: '/tournaments/:tournament/runnerups/',
      name: 'runnerups',
      component: Runnerups
    },
    {
      path: '/tournaments/:tournament/edit/',
      name: 'edit',
      component: Edit
    },
    {
      path: '/tournaments/:tournament/scores/',
      name: 'scores',
      component: ScoreScreen
    },
    {
      path: '/tournaments/:tournament/next/',
      name: 'next',
      component: NextScreen
    },
    {
      path: '/tournaments/:tournament/judge/',
      name: 'judge',
      component: JudgeInterface
    },
    {
      path: '/tournaments/:tournament/charts/',
      name: 'charts',
      component: PostMatch
    },
    {
      path: '/tournaments/:tournament/log/',
      name: 'log',
      component: Log
    },
    {
      path: '/live/',
      name: 'live',
      component: Stream
    },
    {
      path: '/tournaments/:tournament/credits/',
      name: 'credits',
      component: Credits
    },
    {
      path: '/tournaments/:tournament/:match/',
      name: 'match',
      component: Match
    },
    {
      path: '/tournaments/:tournament/*',
      redirect: '/tournaments/:tournament/',
    },
  ],
})

router.beforeEach((to, from, next) => {
  window.scrollTo(0, 0)

  // Why in the world does this need a short timeout? The connect
  // isn't set otherwise.
  setTimeout(function () {
    router.app.connect()

    if (_.isEmpty(router.app.$store.state.tournaments)) {
      router.app.$http.get('/api/tournament/').then(response => {
        let data = JSON.parse(response.data)
        console.debug("tournament state", data)
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
        console.debug("user state", data)

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
        console.debug("stats state", data)
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
    tournaments: {},
    user: new Person(),
    userLoaded: false,
    stats: undefined,
    people: undefined,
    credits: {}
  },
  mutations: {
    updateAll (state, data) {
      let ts = {}
      _.forEach(data.tournaments, (t) => {
        ts[t.id] = Tournament.fromObject(t, data.$vue)
      })
      state.tournaments = ts
    },
    updateTournament (state, data) {
      let t = Tournament.fromObject(data.tournament, data.$vue)
      Vue.set(state.tournaments, t.id, t)
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
      return state.tournaments[id]
    },
    upcoming: state => {
      return _.filter(_.sortBy(state.tournaments, 'scheduled'), 'isUpcoming')
    },
    running: state => {
      return _.filter(_.sortBy(state.tournaments, 'scheduled'), 'isRunning')[0]
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
