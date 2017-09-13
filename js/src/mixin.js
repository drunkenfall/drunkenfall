var DrunkenFallMixin = {
  created () {
  },
  methods: {
  },
  computed: {
    tournament () {
      return this.$store.getters.getTournament(
        this.$route.params.tournament
      )
    },
    tournaments () {
      return this.$store.state.tournaments
    },
    user () {
      return this.$store.state.user
    },
    match () {
      if (!this.$route.params.kind) {
        // Nothing set in params - will fail.
        return
      }
      let kind = this.$route.params.kind
      let idx = this.$route.params.match

      if (kind === 'final') {
        return this.tournament.final
      }
      kind = kind + 's'
      return this.tournament[kind][idx]
    },
    currentMatch () {
      // TODO(thiderman): This needs to be written as the kind that
      // does not change until the next match is started, as per to be
      // used by ScoreScreen and such
      return
    },
    nextMatch () {
      if (!this.tournament.current) {
        // Nothing set - will fail.
        return
      }
      let kind = this.tournament.current.kind
      let idx = this.tournament.current.index

      if (kind === 'final') {
        return this.tournament.final
      }
      kind = kind + 's'
      return this.tournament[kind][idx]
    },
    round () {
      if (!this.match || !this.match.commits) {
        return 1
      }
      return this.match.commits.length + 1
    },
    match_id () {
      if (!this.match) {
        return {}
      }
      return {
        id: this.tournament.id,
        kind: this.match.kind,
        index: this.match.index
      }
    },
    chartMatch () {
      if (!this.tournament.current) {
        // Nothing set - will fail.
        return
      }

      let match
      let index = this.tournament.current.index
      let kind = this.tournament.current.kind

      if (kind === 'final') {
        return this.tournament.final
      }

      kind = kind + 's'

      match = this.tournament[kind][index]

      // We don't want to update until the next match has been
      // started. If we do, the graphs are removed as soon as the
      // judges end the previous match.
      // Also, if we're on the first match there is no previous one,
      // so don't try to grab the previous one in that case.
      if (!match.isStarted || (kind === 'tryouts' && index === 0)) {
        index = this.tournament.previous.index
        kind = this.tournament.previous.kind + 's'
        console.log([index, kind])
        return this.tournament[kind][index]
      }

      return match
    },
  }
}

export default DrunkenFallMixin
