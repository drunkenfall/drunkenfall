<template>
  <div>
    <header>
      <div class="content">
        <div class="title">
          {{tournament.name}} / {{match.kind | capitalize}} {{match.index +1}}
        </div>
      </div>
      <div class="links">
        <a v-if="can_start" @click="start">Start match</a>
        <a v-if="is_running" v-bind:class="{'disabled': !can_commit}" @click="commit">End round</a>
        <a v-if="can_end"@click="end">End match</a>
      </div>
      <div class="clear"></div>
    </header>

    <div class="control">
      <template v-for="player in match.players">
        <control-player :index="$index" :player="player" :match="match"
                        :downs="0" :ups="0">
      </template>
    </div>
    <div class="clear"></div>
  </div>
</template>

<script>
import ControlPlayer from './ControlPlayer.vue'

export default {
  name: 'Match',
  components: {
    ControlPlayer
  },

  data () {
    return {
      match: {},
      tournament: {}
    }
  },

  computed: {
    can_end: function () {
      var m = this.$data.match
      var end = 10

      if (m.ended !== '0001-01-01T00:00:00Z') {
        return false
      }

      if (m.kind === 'final') {
        end = 20
      }

      for (var i = 0; i < m.players.length; i++) {
        var p = m.players[i]
        if (p.kills >= end) {
          return true
        }
      }
      return false
    },
    can_start: function () {
      var m = this.$data.match

      if (m.started === '0001-01-01T00:00:00Z') {
        return true
      }
      return false
    },
    is_running: function () {
      var m = this.$data.match

      if (m.started !== '0001-01-01T00:00:00Z' && m.ended === '0001-01-01T00:00:00Z') {
        return true
      }
      return false
    },
    can_commit: function () {
      return true
    }

  },

  methods: {
    commit: function () {
      var url = '/api/towerfall/tournament/'
      url += this.$data.tournament.id + '/'
      url += this.$data.match.kind + '/'
      url += this.$data.match.index + '/commit/'

      // TODO: pls
      var payload = {
        'state': [
          {
            'ups': this.$children[0].ups,
            'downs': this.$children[0].downs,
            'shot': this.$children[0].shot,
            'reason': this.$children[0].reason
          },
          {
            'ups': this.$children[1].ups,
            'downs': this.$children[1].downs,
            'shot': this.$children[1].shot,
            'reason': this.$children[1].reason
          },
          {
            'ups': this.$children[2].ups,
            'downs': this.$children[2].downs,
            'shot': this.$children[2].shot,
            'reason': this.$children[2].reason
          },
          {
            'ups': this.$children[3].ups,
            'downs': this.$children[3].downs,
            'shot': this.$children[3].shot,
            'reason': this.$children[3].reason
          }
        ]
      }

      console.log(payload)
      this.$http.post(url, payload).then(function (res) {
        console.log(res)
        this.$set('match', res.data.match)

        this.$children[0].reset()
        this.$children[1].reset()
        this.$children[2].reset()
        this.$children[3].reset()
      }, function (res) {
        console.log('error when setting score')
        console.log(res)
      })
    },

    refresh: function () {
      // Hax to make vue refresh the entire page.
      // Since nothing on this page is properly bound to components right now
      // the updates won't trigger properly.
      this.$set('updated', Date.now())
    },
    end: function () {
      var url = '/api/towerfall/tournament/'
      url += this.$data.tournament.id + '/'
      url += this.$data.match.kind + '/'
      url += this.$data.match.index + '/toggle/'

      this.$http.get(url).then(function (res) {
        console.log(res)
        this.$route.router.go('/towerfall/' + this.$data.tournament.id + '/')
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    start: function () {
      var url = '/api/towerfall/tournament/'
      url += this.$data.tournament.id + '/'
      url += this.$data.match.kind + '/'
      url += this.$data.match.index + '/toggle/'

      this.$http.get(url).then(function (res) {
        console.log(res)
        this.setData(
          res.data.tournament,
          this.$data.match.kind,
          this.$data.match.index
        )
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    setData: function (tournament, kind, match) {
      if (kind === 'tryout') {
        kind = 'tryouts'
      } else if (kind === 'semi') {
        kind = 'semis'
      }

      if (kind === 'final') {
        this.$set('match', tournament[kind])
      } else {
        this.$set('match', tournament[kind][match])
      }

      this.$set('tournament', tournament)
    }
  },

  route: {
    data ({ to }) {
      // We need a reference here because `this` inside the callback will be
      // the main App and not this one.
      var $vue = this

      to.router.app.$watch('tournaments', function (newVal, oldVal) {
        for (var i = 0; i < newVal.length; i++) {
          if (newVal[i].id === to.params.tournament) {
            console.log("Match.vue - watch update")
            console.log(newVal[i])
            $vue.setData(
              newVal[i],
              to.params.kind,
              parseInt(to.params.match)
            )
          }
        }
      })

      if (to.router.app.$data.tournaments.length === 0) {
        // Nothing is set - we're reloading the page and we need to get the
        // data manually
        this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
          console.log(res)
          this.setData(
            res.data.Tournament,
            to.params.kind,
            parseInt(to.params.match)
          )
        }, function (res) {
          console.log('error when getting tournament')
          console.log(res)
        })
      } else {
        // Something is set - we're clicking on a link and can reuse the
        // already existing data immediately
        this.setData(
          to.router.app.get(to.params.tournament),
          to.params.kind,
          parseInt(to.params.match)
        )
      }
    }
  }
}
</script>

<style lang="scss" >

@import "../style.scss";

.control {
  height: 85vh;
  padding: 0.8%;
}

.player {
  height: 25%;
  display: block;
}

</style>
