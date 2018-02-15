<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' - DrunkenFall'"></headful>
    <tournament-preview v-if="tournament && !tournament.isStarted"></tournament-preview>
    <tournament-overview v-if="tournament && tournament.isStarted"></tournament-overview>
  </div>
</template>

<script>
import TournamentOverview from '../components/TournamentOverview'
import TournamentPreview from '../components/TournamentPreview'
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'TournamentView',
  mixins: [DrunkenFallMixin],

  components: {
    TournamentOverview,
    TournamentPreview,
  },

  data () {
    return {
      id: null,
    }
  },

  computed: {
    runnerups () {
      let t = this.tournament

      if (!t.runnerups) {
        return []
      }

      return _.map(t.runnerups, (runnerupName) => {
        return _.find(t.players, { name: runnerupName })
      })
    }
  },
}
</script>

<style lang="scss">
</style>
