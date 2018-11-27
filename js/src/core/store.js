import Vue from 'vue'
import Vuex from 'vuex'

import Person from '../models/Person.js'
import Player from '../models/Player.js'
import Match from '../models/Match.js'
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
    playerSummaries: {},
    runnerups: {},
    matches: {},
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
        ts[t.id] = Tournament.fromObject(t)
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
    setPlayerSummaries (state, data) {
      Vue.set(state.playerSummaries, data.tid, _.map(data.player_summaries, (p) => {
        return Player.fromObject(p)
      }))
    },
    setRunnerups (state, data) {
      Vue.set(state.runnerups, data.tid, _.map(data.player_summaries, (p) => {
        return Player.fromObject(p)
      }))
    },
    setMatches (state, data) {
      Vue.set(state.matches, data.tid, _.map(data.matches, (m) => {
        return Match.fromObject(m)
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
      state.people = data.reduce((c, p) => {
        c[p.id] = Person.fromObject(p)
        return c
      }, {})
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
      let t

      switch (res.type) {
        case 'all':
          let ts = {}
          _.forEach(data.tournaments, (t) => {
            ts[t.id] = Tournament.fromObject(t)
          })
          state.tournaments = ts
          break

        case 'tournament':
          t = Tournament.fromObject(data)
          Vue.set(state.tournaments, t.id, t)
          break

        case 'player':
          t = state.tournaments[data.tournament]
          t.matches[data.match].players[data.player].state = data.state
          Vue.set(state.tournaments, t.id, t)
          break

        case 'player_summaries':
          let ps = _.map(data.player_summaries, (p) => {
            return Player.fromObject(p)
          })
          Vue.set(state.playerSummaries, data.tournament_id, ps)
          break

        default:
          console.error('Unknown websocket update:', res)
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
    playerSummaries: (state, getters) => (id) => {
      return state.playerSummaries[id]
    },
    runnerups: (state, getters) => (id) => {
      return state.runnerups[id]
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
    matchPlayers: (state, getters) => (tid, idx) => {
      let t = state.matches[tid]
      if (!t) {
        return
      }

      let ret = t.find(m => m.index === idx)
      return ret.players
    },
    getPerson: (state, getters) => (id) => {
      console.log("getting person", id)
      return state.people[id]
    },
    isConnected: state => {
      return state.socket.isConnected
    },
  }
})

export default store
