<template>
<div v-if="tournament">
  <tournament-preview v-if="!tournament.isStarted"></tournament-preview>
  <next-screen class="local" v-else-if="tournament.betweenMatches"></next-screen>
  <statusbar v-else-if="match && match.isStarted && !tournament.isEnded"></statusbar>
  <tournament-summary v-else-if="tournament.isEnded"></tournament-summary>
</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import NextScreen from './NextScreen'
import TournamentPreview from './TournamentPreview'
import TournamentSummary from './TournamentSummary'
import Statusbar from './Statusbar'

export default {
  name: 'HUD',
  mixins: [DrunkenFallMixin],
  components: {
    NextScreen,
    TournamentPreview,
    TournamentSummary,
    Statusbar,
  },
  computed: {
    tournament () {
      return this.tournaments[this.$route.params.tournament]
    },
    match () {
      return this.tournament.currentMatch
    }
  },
  created () {
    this.loadAll()
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"
  },
}
</script>

<style lang="scss" scoped>

</style>
