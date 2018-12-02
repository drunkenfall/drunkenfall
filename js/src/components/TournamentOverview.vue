<template>

<div v-if="tournament && matchesLoaded" class="main">
  <headful :title="tournament.subtitle + ' - DrunkenFall'"></headful>
  <tournament-controls />

  <div class="overview">
    <div class="ongoing">
      <h1>Next up!</h1>

      <match :match="match" class="match"></match>

      <div class="logo">
        <img :class="{ded: !isConnected}" alt="One-Eye" src="/static/img/oem.svg"/>
      </div>
    </div>

    <div class="players" v-if="nextMatch">
      <h1 v-if="tournament.qualifyingOpen">Next scheduled</h1>
      <h1 v-else>Last qualifying</h1>
      <template v-for="(p, x) in nextMatch.players">
        <list-player :player="p" :match="tournament.nextNextMatch" :index="x"></list-player>
      </template>

      <div class="active" v-if="tournament.qualifyingOpen">
        <h1>In queue</h1>
        <template v-for="(p, x) in tournament.runnerups">
          <list-player :player="p" :index="x"></list-player>
        </template>
      </div>
    </div>
  </div>
</div>

</template>

<script>
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"
import Match from "./Match"
import ListPlayer from "./ListPlayer"
// import _ from 'lodash'

export default {
  name: 'TournamentOverview',
  mixins: [DrunkenFallMixin],

  components: {
    TournamentControls,
    Match,
    ListPlayer,
  },

  methods: {
    setTime (x) {
      let $vue = this
      this.api.setTime({ id: this.tournament.id, time: x }).then((res) => {
        console.debug("settime response:", res)
      }, (err) => {
        $vue.$alert("Setting time failed. See console.")
        console.error(err)
      })
    },
    usurp () {
      this.tournament.usurp()
    },
  },

  computed: {
    match () {
      return this.tournament.nextMatch
    },
    nextMatch () {
      return this.tournament.nextNextMatch
    },
  },

  created () {
    let $vue = this
    let id = this.tournament.id

    this.$http.get(`/api/tournaments/${id}/matches/`).then(function (res) {
      let data = JSON.parse(res.data)
      this.$store.commit('setMatches', {
        tid: id,
        matches: data.matches,
      })
    }, function (res) {
      $vue.$alert("Getting players failed. See console.")
      console.error(res)
    })

    this.$http.get(`/api/tournaments/${id}/players/`).then(function (res) {
      let data = JSON.parse(res.data)
      this.$store.commit('setPlayerSummaries', {
        tid: id,
        player_summaries: data.player_summaries,
      })
    }, function (res) {
      $vue.$alert("Getting players failed. See console.")
      console.error(res)
    })

    this.$http.get(`/api/tournaments/${id}/runnerups/`).then(function (res) {
      let data = JSON.parse(res.data)
      this.$store.commit('setRunnerups', {
        tid: id,
        player_summaries: data.player_summaries,
      })
    }, function (res) {
      $vue.$alert("Getting players failed. See console.")
      console.error(res)
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.main {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow: hidden;
}

.overview {
  display: flex;
  flex-grow: 1;
  justify-content: space-between;

  .ongoing, .players {
    padding: 0 5em;
  }

  .ongoing {
    flex-basis: 60%;

    display: flex;
    flex-direction: column;

    .match {
      display: flex;
      flex-grow: 1;
      flex-direction: column;
    }

    .logo {
      img {
        height: 150px;
      }
    }
  }

  .players {
    flex-basis: 35%;
  }
}

</style>
