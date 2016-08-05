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
        <div class="player {{player.preferred_color}}">
          <div class="button">
            <div><p>-</p></div>
          </div>

          <div class="shots">
            <!--- <div><p>✓ {{player.shots}} ✗</p></div> -->
            <div>
              <div class="mark">✗</div>
              <div class="reason"></div>
            </div>
          </div>

          <div class="slider {{player.preferred_color}}">
            <div><p>{{player.name}}</p></div>
          </div>

          <div class="scores">
            <div v-for="n in 10">
              <p>{{n+1}}</p>
            </div>
          </div>

          <div class="button">
            <div><p>+</p></div>
          </div>
        </div>
      </template>
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

<style lang="scss" >

@import "../style.scss";

.control {
  height: 85vh;
  padding: 0.8%;

  .player {
    width: 100%;
    position: relative;

    >div {
      float: left;
      height: 100%;
    }

    .button, .shots {
      width: 10%;
      position: relative;
      text-shadow: 2px 2px 2px rgba(0,0,0,0.8);

      >div {
        width: 80%;
        height: 50%;
        background-color: #333339;
        cursor: pointer;
        text-align: center;
        vertical-align: middle;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translateX(-50%) translateY(-50%);
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
    .button {
      div {
        font-size: 7em;
        p {
          margin-top: -13%;
        }
      }
    }

    .shots div {
      width: 80%;
      display: block;
      margin: 0 auto 1%;

      &.mark {
        padding-top: 2%;
        font-size: 4em;
      }
      &.reason {
        margin-top: -4%;
        font-size: 1.5em;
      }
    }

    .slider {
      width: 15%;
      text-shadow: 2px 2px 2px rgba(0,0,0,0.8);
      font-weight: bold;
      font-size: 1.6vw;

      &.green  { color: $green ; }
      &.blue   { color: $blue  ; }
      &.pink   { color: $pink  ; }
      &.orange { color: $orange; }
      &.white  { color: $white ; }
      &.yellow { color: $yellow; }
      &.cyan   { color: $cyan  ; }
      &.purple { color: $purple; }

      div {
        width: 80%;
        height: 50%;
        font-size: 1.6em;
        cursor: pointer;

        text-align: center;
        vertical-align: middle;
        position: relative;
        top: 50%;
        left: 50%;
        transform: translateX(-50%) translateY(-50%);
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
      }
    }
    .scores {
      // http://stackoverflow.com/questions/6865194/fluid-width-with-equally-spaced-divs
      width: 55%;
      text-align: justify;
      text-justify: distribute;
      // background-color: #333339;

      div {
        // position: relative;
        // top: 25%;
        width: 8%;
        margin: 0 1%;
        height: 8vh;
        vertical-align: top;
        display: inline-block;
        *display: inline;
        zoom: 1;

        background-color: #333339;
        border-radius: 10000px;

        text-align: center;
        vertical-align: middle;
        position: relative;
        top: 50%;
        transform: translateX(0) translateY(-50%);

        p {
          position: relative;
          top: 50%;
          transform: translateX(0) translateY(-50%);
          font-size: 1.5em;
          color: #666672;
        }
      }
    }
  }
}

.player {
  height: 25%;
  display: block;
}

</style>
