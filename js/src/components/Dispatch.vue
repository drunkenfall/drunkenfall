<template>
  <div>

  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Dispatch',
  mixins: [DrunkenFallMixin],
  methods: {
    dispatch () {
      if (this.upcomingTournament !== undefined) {
        this.$router.push({
          name: 'tournament',
          params: {
            tournament: this.upcomingTournament.id
          },
        })
      } else {
        this.$router.push({name: "tournaments"})
      }
    },
  },
  mounted () {
    if (this.tournaments !== undefined) {
      this.dispatch()
    }
  },
  watch: {
    // Wait for data to be set...
    tournaments (val, old) {
      if (val === undefined || old === undefined) {
        return
      }

      this.dispatch()
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
