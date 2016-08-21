<template>
  <router-view
    keep-alive
    transition>
  </router-view>
</template>

<script>
/* eslint-env browser */

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
      if (this.$data.ws === null) {
        console.log('Setting up new websocket')
        this.$data.ws = new WebSocket('ws://' + window.location.host + '/api/towerfall/auto-updater')

        // We need to be able to reference back to the Vue app instance from
        // inside of the websocket.
        this.$data.ws.$vue = this

        this.$data.ws.onmessage = function (e) {
          var res = JSON.parse(e.data)

          // If p is set, this is a ping message that only serves to keep the connection open.
          // Break immediately.
          if (res.p !== undefined) {
            return
          }

          console.log(e)
          if (res.data !== undefined) {
            // The main bulk update. This contains the latest state.
            if (res.data.tournaments !== undefined) {
              console.log('Updating tournaments')
              this.$vue.$set('tournaments', res.data.tournaments)
              return
            }

            console.log('Did not set')
          }

          console.log('Unknown websocket update:')
          console.log(res)
        }

        console.log(this.$data.ws)
      }
    },

    populate: function () {
      if (this.$data.tournaments.length === 0) {
        console.log('Grabbing initial set of tournament data')
        this.$http.get('/api/towerfall/tournament/').then(function (res) {
          this.$set('tournaments', res.data)
        }, function (res) {
          console.log('error when getting tournaments')
          console.log(res)
        })
      }
    },

    loadInitial: function ($vue, tid) {
      this.$http.get('/api/towerfall/tournament/' + tid + '/').then(function (res) {
        console.log("returned tournament")
        console.log(res.data.Tournament)
        $vue.$set('tournament', res.data.Tournament)
      }, function (res) {
        console.log('error when getting tournament')
      })
    },

    get: function (tid) {
      for (var i = 0; i < this.$data.tournaments.length; i++) {
        if (this.$data.tournaments[i].id === tid) {
          return this.$data.tournaments[i]
        }
      }
    }
  }
}
</script>

<style lang="scss">
@import "./style.scss";

@font-face { font-family: Archer; src: url('/static/Archer.ttf'); }

</style>
