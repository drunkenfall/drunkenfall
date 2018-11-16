import Vue from 'vue'
import Vuex from 'vuex'

import Person from '../models/Person.js'
import Stats from '../models/Stats.js'
import {Credits as CreditsModel} from '../models/Credits.js'
import Tournament from '../models/Tournament.js'

import _ from 'lodash'

Vue.use(Vuex)

const store = new Vuex.Store({ // eslint-disable-line
  state: {
    tournaments: {},
    user: new Person(),
    tournamentsLoaded: false,
    userLoaded: false,
    stats: undefined,
    people: undefined,
    credits: {},
    socket: {
      isConnected: false,
      reconnectError: false,
    },
  },
  mutations: {
    updateAll (state, data) {
      let ts = {}
      _.forEach(data.tournaments, (t) => {
        ts[t.slug] = Tournament.fromObject(t)
      })
      state.tournaments = ts
      state.tournamentsLoaded = true
    },
    updateTournament (state, data) {
      let t = Tournament.fromObject(data.tournament, data.$vue)
      Vue.set(state.tournaments, t.id, t)
    },
    updatePlayer (state, data) {
      let t = state.tournaments[data.tournament]
      t.matches[data.match].players[data.player].state = data.state
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
    SOCKET_ONOPEN (state, event) {
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state, event) {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event) {
      console.error(state, event)
    },
    // default handler called for all methods

    SOCKET_ONMESSAGE (state, res) {
      let data = res.data
      if (res.type === 'all') {
        let ts = {}
        _.forEach(data.tournaments, (t) => {
          ts[t.id] = Tournament.fromObject(t)
        })
        state.tournaments = ts
      } else if (res.type === 'tournament') {
        let t = Tournament.fromObject(data)
        Vue.set(state.tournaments, t.id, t)
      } else if (res.type === 'player') {
        let t = state.tournaments[data.tournament]
        t.matches[data.match].players[data.player].state = data.state
        Vue.set(state.tournaments, t.id, t)
      } else {
        console.log('Unknown websocket update:', res)
      }
    },

    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true
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
    latest: state => {
      return _.reverse(_.filter(_.sortBy(state.tournaments, 'scheduled'), 'isEnded'))[0]
    },
    getPerson: (state, getters) => (id) => {
      if (!state.people) {
        return undefined
      }
      return state.people.find(p => p.id === id)
    },
  }
})

export default store
