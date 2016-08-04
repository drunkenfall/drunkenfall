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
      <div id="control" class="match">

        <template v-for="player in match.players">
          <div class="player {{player.preferred_color}}">
            <p class="name">
              {{player.name}}
            </p>

            <div class="control-group">

              <div class="control kills">
                <div @click="score($index, 'kills', 'down')" class="sign">
                  <p v-show="is_running">-</p>
                </div >
                <div class="count">{{player.kills}} kills</div>
                <div @click="score($index, 'kills', 'up')" class="sign">
                  <p v-show="is_running">+</p>
                </div>
                <div class="clear"></div>
              </div>

              <div class="control shots">
                <div @click="score($index, 'shots', 'down')" class="sign">
                  <p v-show="is_running">-</p>
                </div>
                <div class="count">{{player.shots}} shots</div>
                <div @click="score($index, 'shots', 'up')" class="sign">
                  <p v-show="is_running">+</p>
                </div>
                <div class="clear"></div>
              </div>

              <div class="control sweeps">
                <div @click="score($index, 'sweeps', 'down')" class="sign">
                  <p v-show="is_running">-</p>
                </div>
                <div class="count">{{player.sweeps}} sweeps</div>
                <div @click="score($index, 'sweeps', 'up')" class="sign">
                  <p v-show="is_running">+</p>
                </div>
                <div class="clear"></div>
              </div>

              <div class="control self">
                <div @click="score($index, 'self', 'down')" class="sign">
                  <p v-show="is_running">-</p>
                </div>
                <div class="count">{{player.self}} self</div>
                <div @click="score($index, 'self', 'up')" class="sign">
                  <p v-show="is_running">+</p>
                </div>
                <div class="clear"></div>
              </div>

              <div class="control explosions">
                <div @click="score($index, 'explosions', 'down')" class="sign">
                  <p v-show="is_running">-</p>
                </div>
                <div class="count">{{player.explosions}} explosions</div>
                <div @click="score($index, 'explosions', 'up')" class="sign">
                  <p v-show="is_running">+</p>
                </div>
                <div class="clear"></div>
              </div>

              <div class="clear"></div>
            </div>
          </div>
        </template>

      </div>
    </div>
    <div class="clear"></div>

</div>

</template>

<script>
import Player from './Player.vue'
export default {
  name: 'Match',
  components: {
    Player
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
    score: function (player, action, direction) {
      var url = '/api/towerfall/tournament/'
      url += this.$data.tournament.id + '/'
      url += this.$data.match.kind + '/'
      url += this.$data.match.index + '/'
      url += player + '/'
      url += action + '/'
      url += direction

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

<style lang="scss" scoped>
.player {
  min-height: 2.2em;

  .name {
    margin-top: 0.3em;
    font-size: 1.5em;
  }
  .control-group {
    margin: 3% auto;
    height: 16%;
    width: 80%;
    position: relative;
  }
  .control {
    height: 16%;
    display: block;
    position: relative;
    line-height: 1.7em;

    div, a {
      float: left;
      display: block;
      margin: 1% 0;
      height: 100%;
    }
    .count {
      width: 70%;
      background-color: #404040;
    }
    .sign {
      cursor: pointer;
      background-color: #353535;
      width: 15%;
      cursor: pointer;
      display: block;
    }
  }
}
.shots .count, .self .count {
  background-color: #484848 !important;
}
.shots .sign, .self .sign {
  background-color: #404040 !important;
}

</style>
