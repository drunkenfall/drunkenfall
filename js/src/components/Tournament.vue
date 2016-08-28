<template>
  <div>
    <header>
      <div class="content">
        <div class="title">{{tournament.name}}</div>
      </div>
      <div class="links">
        <a v-if="tournament.canStart" v-link="{path: 'join/'}">Join</a>
        <div class="action" v-if="tournament.canStart" @click="start">Start</div>
        <div class="action" v-if="tournament.isRunning" @click="next">Next match</div>
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

      <div v-if="runnerups.length > 0">
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
import Tournament from '../models/Tournament.js'
import _ from 'lodash'

export default {
  name: 'Tournament',

  components: {
    MatchOverview
  },

  data () {
    return {
      tournament: new Tournament(),
    }
  },

  computed: {
    runnerups: function () {
      var ret = []
      var t = this.tournament

      if (!t.runnerups) {
        return ret
      }

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
        // XXX: Worst hack of all time
        this.$data.tournament.started = 'hehe'
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
      to.router.app.$watch('tournaments', (newVal, oldVal) => {
        let thisTournament = _.find(newVal, { id: to.params.tournament })
        console.debug('update tournament with new data', thisTournament)
        this.$set('tournament', Tournament.fromObject(thisTournament))
      })

      if (to.router.app.tournaments.length === 0) {
        // Nothing is set - we're reloading the page and we need to get the
        // data manually
        to.router.app.loadInitial(to.params.tournament)
      } else {
        // Something is set - we're clicking on a link and can reuse the
        // already existing data immediately
        this.$set('tournament', to.router.app.get(to.params.tournament))
      }
    }
  }
}
</script>

<style lang="scss">
@import "../style.scss";

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

    &.green  { background-color: $green ; }
    &.blue   { background-color: $blue  ; }
    &.pink   { background-color: $pink  ; }
    &.orange { background-color: $orange; }
    &.white  { background-color: $white ; }
    &.yellow { background-color: $yellow; }
    &.cyan   { background-color: $cyan  ; }
    &.purple { background-color: $purple; }
    &.red    { background-color: $red; }

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
