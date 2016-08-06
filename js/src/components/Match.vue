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
        <a v-if="!can_start" v-bind:class="{'disabled': !can_end}" @click="end">End match</a>
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
    }

  },

  methods: {
    commit: function () {
      var url = '/api/towerfall/tournament/'
      url += this.$data.tournament.id + '/'
      url += this.$data.match.kind + '/'
      url += this.$data.match.index + '/commit/'

      this.$http.get(url).then(function (res) {
        console.log(res)
        var target = this.$data.match.kind
        var match = this.$data.match.index

        if (target === 'tryout') {
          target = 'tryouts'
        } else if (target === 'semi') {
          target = 'semis'
        }

        if (target === 'final') {
          this.$set('match', res.data.tournament[target])
        } else {
          this.$set('match', res.data.tournament[target][match])
        }

        this.$set('tournament', res.data.tournament)
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
        var target = this.$data.match.kind
        var match = this.$data.match.index

        if (target === 'tryout') {
          target = 'tryouts'
        } else if (target === 'semi') {
          target = 'semis'
        }

        if (target === 'final') {
          this.$set('match', res.data.tournament[target])
        } else {
          this.$set('match', res.data.tournament[target][match])
        }

        this.$set('tournament', res.data.tournament)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    }
  },

  route: {
    data ({ to }) {
      this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
        console.log(res)
        var target = to.params.kind
        var match = parseInt(to.params.match)

        if (target === 'tryout') {
          target = 'tryouts'
        } else if (target === 'semi') {
          target = 'semis'
        }

        if (target === 'final') {
          this.$set('match', res.data.Tournament[target])
        } else {
          this.$set('match', res.data.Tournament[target][match])
        }

        this.$set('tournament', res.data.Tournament)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
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
