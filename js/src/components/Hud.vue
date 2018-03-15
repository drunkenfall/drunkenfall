<template>
<div v-if="tournament">
  <tournament-preview v-if="!tournament.isStarted"></tournament-preview>
  <next-screen class="local" v-else-if="tournament.betweenMatches"></next-screen>
  <statusbar v-else-if="match.isStarted"></statusbar>
</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import NextScreen from './NextScreen'
import TournamentPreview from './TournamentPreview'
import Statusbar from './Statusbar'

export default {
  name: 'HUD',
  mixins: [DrunkenFallMixin],
  components: {
    NextScreen,
    TournamentPreview,
    Statusbar,
  },
  computed: {
    tournament () {
      if (this.runningTournament) {
        return this.runningTournament
      } else if (this.upcomingTournament && this.upcomingTournament.isToday) {
        return this.upcomingTournament
      }
    },
    match () {
      return this.tournament.currentMatch
    }
  },
  created () {
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"
  },
}
</script>

<style lang="scss" scoped>

</style>
