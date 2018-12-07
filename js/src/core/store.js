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
    initialLoaded: false,
    userLoaded: false,
    stats: undefined,
    people: undefined,
    playerSummaries: {},
    playerStates: {},
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
      let m = t.currentMatch
      console.log("updateplayer tournament", t)
      console.log("updateplayer currentmatch", m)
      m.players[data.player].state = data.state
      Vue.set(state.tournaments, t.id, t)
    },
    setInitial (state, data) {
      let t = Tournament.fromObject(data.tournament, data.$vue)
      Vue.set(state.tournaments, t.id, t)

      Vue.set(state.playerSummaries, t.id, _.map(data.player_summaries, (p) => {
        return Player.fromObject(p)
      }))

      Vue.set(state.runnerups, t.id, _.map(data.player_summaries, (p) => {
        return Player.fromObject(p)
      }))

      Vue.set(state.matches, t.id, _.map(data.matches, (m) => {
        return Match.fromObject(m)
      }))

      let sts = {}
      _.forEach(data.player_states, (st) => {
        sts[st.index] = st
      })
      state.playerStates = sts

      state.initialLoaded = true
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
      console.debug("Incoming websocket update", state, res)

      let data = res.data
      let t, ps, ms, m

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
          Vue.set(state.playerStates, data.state.index, data.state)
          break

        case 'player_summaries':
          ps = _.map(data.player_summaries, (p) => {
            return Player.fromObject(p)
          })
          Vue.set(state.playerSummaries, data.tournament_id, ps)
          break

        case 'runnerups':
          ps = _.map(data.runnerups, (p) => {
            return Player.fromObject(p)
          })
          Vue.set(state.runnerups, data.tournament_id, ps)
          break

        case 'match':
          t = state.tournaments[data.tournament]
          m = Match.fromObject(data.match)
          Vue.set(state.matches[t.id], m.index, m)
          break

        case 'matches':
          ms = _.map(data.matches, (m) => {
            return Match.fromObject(m)
          })
          Vue.set(state.matches, data.tournament_id, ms)
          break

        case 'match_end':
          t = Tournament.fromObject(data.tournament)

          ms = _.map(data.matches, (m) => {
            return Match.fromObject(m)
          })

          ps = _.map(data.player_summaries, (p) => {
            return Player.fromObject(p)
          })

          let rups = _.map(data.runnerups, (p) => {
            return Player.fromObject(p)
          })

          Vue.set(state.runnerups, t.id, rups)
          Vue.set(state.matches, t.id, ms)
          Vue.set(state.playerSummaries, t.id, ps)
          Vue.set(state.tournaments, t.id, t)
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

      let m = t.find(m => m.index === idx)
      return _.map(_.sortBy(m.players, "id"), Player.fromObject)
    },
    getPerson: (state, getters) => (id) => {
      if (state.people === undefined) {
        return
      }

      return state.people[id]
    },
    getPlayerSummary: (state, getters) => (tid, id) => {
      let ps = state.playerSummaries[tid]
      if (ps === undefined) {
        console.log("player summaries undefined", ps)
        return
      }
      return ps.find(s => s.person_id === id)
    },
    getPlayerState: (state, getters) => (idx) => {
      return state.playerStates[idx]
    },
    getMatches: (state, getters) => (id) => {
      return state.matches[id]
    },
    isConnected: state => {
      return state.socket.isConnected
    },
  }
})

export default store
