<template>
  <div v-if="tournament">
    <next-screen class="local" v-if="tournament.betweenMatches"></next-screen>
    <statusbar v-else-if="match.isStarted"></statusbar>
  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import NextScreen from './NextScreen'
import Statusbar from './Statusbar'

export default {
  name: 'HUD',
  mixins: [DrunkenFallMixin],
  components: {
    NextScreen,
    Statusbar,
  },
  computed: {
    tournament () {
      return this.runningTournament
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
