<template>
  <router-view
    keep-alive
    transition>
  </router-view>
</template>

<script>
/* eslint-env browser */
import _ from 'lodash'

export default {
  data () {
    return {
      // All the tournaments
      // TODO(thiderman): It's practically pointless to always update all of them
      // but this was the quickest way. In the future, we can probably just update
      // either the match, or just the running tournament.
      tournaments: [],
      // The main websocket object
      ws: null,
    }
  },
  methods: {
    // (Re-)Connect the websocket.
    // Is safe to run when the connection is already up - then it will be a noop.
    connect: function () {
      if (this.ws) {
        console.log('Setting up new websocket')
        this.$set('ws', new WebSocket('ws://' + window.location.host + '/api/towerfall/auto-updater'))

        this.ws.onmessage = (event) => {
          let res
          try {
            res = JSON.parse(event.data)
          } catch (e) {
            console.error("Failed to parse message data:", e)
            return
          }

          // If p is set, this is a ping message that only serves to keep the connection open.
          // Break immediately.
          if (res.p !== undefined) {
            return
          }

          console.log(event)
          if (res.data) {
            if (res.data.tournaments) {
            // The main bulk update. This contains the latest state.
              console.log('Updating tournaments')
              this.$set('tournaments', res.data.tournaments)
              return
            }

            console.log('no tournaments received, did not update anything')
            return
          }

          console.log('Unknown websocket update:', res)
        }

        this.ws.onerror = (errorEvent) => { console.error("websocket error:", errorEvent) }
        this.ws.onclose = (closeEvent) => { console.warn("websocket closed", closeEvent) }

        console.log("WebSocket:", this.ws)
      }
    },

    populate: function () {
      if (!this.tournaments || this.tournaments.length === 0) {
        console.log('Grabbing initial set of tournament data')
        this.$http.get('/api/towerfall/tournament/').then((res) => {
          this.$set('tournaments', res.data)
        }, (error) => {
          console.log('error when getting tournaments:', error)
        })
      }
    },

    loadInitial: function (tid) {
      console.log("tournament", tid)
      this.$http.get('/api/towerfall/tournament/' + tid + '/').then((res) => {
        console.log("returned tournament")
        console.log(res.data.Tournament)
        this.$set('tournament', res.data.Tournament)
      }, (error) => {
        console.log('error when getting tournament', error)
      })
    },

    get: function (tid) {
      return _.find(this.tournaments, { id: tid })
    }
  }
}
</script>

<style lang="scss">
@import "./style.scss";

@font-face { font-family: Archer; src: url('/static/Archer.ttf'); }

</style>
