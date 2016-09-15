<template>
  <div>
    <tournament-preview v-if="!tournament.isStarted" :tournament="tournament" :user="user"></tournament-preview>
    <tournament-overview v-if="tournament.isStarted" :tournament="tournament" :user="user"></tournament-overview>
  </div>
</template>

<script>
import Tournament from '../models/Tournament'
import TournamentOverview from '../components/TournamentOverview'
import TournamentPreview from '../components/TournamentPreview'
import * as levels from "../models/Level"
import _ from 'lodash'

export default {
  name: 'Tournament',

  components: {
    TournamentOverview,
    TournamentPreview,
  },

  data () {
    return {
      tournament: new Tournament(),
      user: this.$root.user,
      levels: levels,
    }
  },

  computed: {
    runnerups: function () {
      let t = this.tournament

      if (!t.runnerups) {
        return []
      }

      return _.map(t.runnerups, (runnerupName) => {
        return _.find(t.players, { name: runnerupName })
      })
    }
  },

  methods: {
    start: function () {
      if (this.tournament) {
        this.api.start({ id: this.tournament.id }).then((res) => {
          console.log("start response:", res)
          let j = res.json()
          this.$route.router.go('/towerfall' + j.redirect)
        }, (err) => {
          console.error(`start for ${this.tournament} failed`, err)
        })
      } else {
        console.error("start called with no tournament")
      }
    },
    next: function () {
      if (this.tournament) {
        this.api.next({ id: this.tournament.id }).then((res) => {
          console.debug("next response:", res)
          let j = res.json()
          this.$route.router.go('/towerfall' + j.redirect)
        }, (err) => {
          console.error(`next for ${this.tournament} failed`, err)
        })
      } else {
        console.error("next called with no tournament")
      }
    }
  },

  created: function () {
    console.debug("Creating API resource")
    let customActions = {
      start: { method: "GET", url: "/api/towerfall{/id}/start/" },
      next: { method: "GET", url: "/api/towerfall{/id}/next/" },
      getData: { method: "GET", url: "/api/towerfall/tournament{/id}/" }
    }
    this.api = this.$resource("/api/towerfall", {}, customActions)
  },

  route: {
    data ({ to }) {
      // listen for tournaments from App
      this.$on(`tournament${to.params.tournament}`, (tournament) => {
        console.debug("New tournament from App:", tournament)
        this.$set('tournament', tournament)
      })

      // TODO perhaps use $root.tournaments again?
      return this.api.getData({ id: to.params.tournament }).then((res) => {
        let tournament = Tournament.fromObject(res.data.tournament)
        console.debug("loaded tournament", tournament)
        return {
          tournament: tournament
        }
      }, (error) => {
        console.error('error when getting tournament', error)
        return { tournament: new Tournament() }
      })
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

    &:nth-child(1), &:nth-child(4) {
      background-color: #333;
    }
    &:nth-child(2), &:nth-child(3) {
      background-color: #383838;
    }
    &.prefill {
      color: #555;
    }

    &.green  { color: $green ; }
    &.blue   { color: $blue  ; }
    &.pink   { color: $pink  ; }
    &.orange { color: $orange; }
    &.white  { color: $white ; }
    &.yellow { color: $yellow; }
    &.cyan   { color: $cyan  ; }
    &.purple { color: $purple; }
    &.red    { color: $red; }

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
