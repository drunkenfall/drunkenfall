<template>
  <div>
    <header>
      <div class="content">
        <div class="title">{{tournament.name}}</div>
      </div>
      <div class="links">
        <a v-if="can_join" v-link="{path: 'join/'}">Join</a>
        <div class="action" v-if="can_start" @click="start">Start</div>
        <div class="action" v-if="is_running" @click="next">Next match</div>
      </div>
      <div class="clear"></div>
    </header>

    <div class="category tryouts">
      <h3>Tryouts</h3>
      <div class="matches">
        <template v-for="m in tournament.tryouts">
          <match-overview :match="m" class="match {{m.kind}}">
        </template>
      </div>
      <div class="clear"></div>
    </div>

    <div class="category semis">
      <h3>Semi-finals</h3>
      <div class="matches">
        <template v-for="m in tournament.semis">
          <match-overview :match="m" class="match {{m.kind}}">
        </template>
      </div>
      <div class="clear"></div>

      <h3>Runnerups</h3>
      <div class="runnerups">
        <template v-for="player in runnerups">
          <div class="runnerup">
            <p class="name">{{player.name}}</p>
            <p class="score">
              <b>{{player.score}}</b> points
              /
              <b>{{player.matches}}</b> matches
            </p>

            <div class="clear"></div>
          </div>
        </template>
      </div>

    </div>
    <div class="category final">
      <h3>Final</h3>
      <div class="matches">
        <match-overview :match="tournament.final" class="match final">
      </div>
    </div>
  </div>
</template>

<script>
import MatchOverview from './MatchOverview.vue'

export default {
  name: 'Tournament',

  components: {
    MatchOverview
  },

  data () {
    return {
      tournament: {
        players: [],
        runnerups: []
      },
      can_join: false,
      can_start: true,
      is_running: false
    }
  },

  computed: {
    can_start: function () {
      // Such is the default nil format in Go
      return this.tournament.started === '0001-01-01T00:00:00Z'
    },
    is_running: function () {
      return this.tournament.started !== '0001-01-01T00:00:00Z' && this.tournament.ended === '0001-01-01T00:00:00Z'
    },
    runnerups: function () {
      var ret = []
      var t = this.tournament

      for (var i = 0; i < t.runnerups.length; i++) {
        for (var j = 0; j < t.players.length; j++) {
          var runnerupName = t.runnerups[i]
          var player = t.players[j]

          if (runnerupName === player.name) {
            ret.push(player)
          }
        }
      }

      return ret
    }
  },

  methods: {
    start: function () {
      this.$http.get('/api/towerfall/' + this.$data.tournament.id + '/start/').then((res) => {
        console.log(res)
        var j = res.json()
        this.$route.router.go('/towerfall' + j.redirect)
      }, (res) => {
        console.log('fail')
        console.log(res)
      })
    },
    next: function () {
      this.$http.get('/api/towerfall/' + this.$data.tournament.id + '/next/').then((res) => {
        console.log(res)
        var j = res.json()
        this.$route.router.go('/towerfall' + j.redirect)
      }, (res) => {
        console.log('fail')
        console.log(res)
      })
    }
  },

  route: {
    data ({ to }) {
      this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
        this.$set('tournament', res.data.Tournament)
        this.$set('can_join', res.data.CanJoin)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    }
  }
}
</script>

<style lang="scss">
.tournament {
  position: relative;
}

.tryouts, .semis, .final {
  width: 29%;
  float: left;
  margin-left: 3%;
  position: relative;
}

.category h3 {
  text-align: center;
  font-size: 200%;
  margin: 4%;
}

.match {
  width: 100%;
  font-size: 150%;
  display: block;

  .player {
    float: left;
    width: 50%;
    height: 40%;
    line-height: 120%;
    text-align: center;
    overflow: hidden;
    white-space: nowrap;
    padding: 0.2em 0;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);

    &.prefill:nth-child(1), &.prefill:nth-child(4) {
      background-color: #333;
    }
    &.prefill:nth-child(2), &.prefill:nth-child(3) {
      background-color: #383838;
    }
    &.prefill {
      color: #555;
    }

    &.green  { background-color: #4E9110; }
    &.blue   { background-color: #4C7CBA; }
    &.pink   { background-color: #E39BB5; }
    &.orange { background-color: #CF9648; }
    &.white  { background-color: #dbdbdb; }
    &.yellow { background-color: #D1BD66; }
    &.cyan   { background-color: #59C2C1; }
    &.purple { background-color: #762c7a; }

    &.gold {
      background-color: #daa520;
    }
    &.silver {
      background-color: #999;
    }
    &.bronze {
      background-color: #8C7853;
    }
    &.out {
      color: #777;
      text-shadow: 1px 1px 5px rgba(0,0,0,0.3);
    }
    &.out:nth-child(1), &.out:nth-child(4) {
      background-color: #433;
    }
    &.out:nth-child(2), &.out:nth-child(3) {
      background-color: #4a3838;
    }

  }

  &.gold, &.silver, &.bronze {
    color: #fff;
    text-shadow: 1px 1px 5px rgba(0,0,0,0.3);
  }

}

.runnerups {
  width: 100%;
  margin: 10px;

  .runnerup {
    padding: 0.1em 0.3em;
    font-size: 24px;
    color: #aaa;

    p {
      margin: 1px;
      &.name {
        float: left;
        // font-weight: bold;
      }
      &.score {
        float: right;
      }
      b {
        text-shadow: 1px 1px 1px rgba(0,0,0,0.4);
      }
    }
  }

  .runnerup:nth-child(odd) {
    background-color: #333;
  }
  .runnerup:nth-child(even) {
    background-color: #272727;
  }
}

</style>
