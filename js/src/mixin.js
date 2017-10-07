import _ from "lodash"

var DrunkenFallMixin = {
  created () {
  },
  methods: {
    $alert (msg) {
      console.error(msg)

      this.$toast(msg, {
        className: ['event', 'alert'],
        horizontalPosition: 'right',
        verticalPosition: 'bottom',
        duration: 5000,
        mode: 'queue',
      })
    },
    $warn (msg) {
      console.warn(msg)

      this.$toast(msg, {
        className: ['event', 'warning'],
        horizontalPosition: 'right',
        verticalPosition: 'bottom',
        duration: 5000,
        mode: 'queue',
      })
    },
    $info (msg) {
      console.log(msg)

      this.$toast(msg, {
        className: ['event'],
        horizontalPosition: 'right',
        verticalPosition: 'bottom',
        duration: 5000,
        mode: 'queue',
      })
    },
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
    userLoaded () {
      return this.$store.state.userLoaded
    },
    match () {
      if (this.$route.params.match === undefined) {
        // Nothing set in params - will fail.
        return
      }
      return this.tournament.matches[this.$route.params.match]
    },
    currentMatch () {
      // TODO(thiderman): This needs to be written as the kind that
      // does not change until the next match is started, as per to be
      // used by ScoreScreen and such
      return
    },
    nextMatch () {
      if (this.tournament.current === undefined) {
        // Nothing set - will fail.
        return
      }
      return this.tournament.matches[this.tournament.current]
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
        index: this.match.index
      }
    },
    chartMatch () {
      if (!this.tournament.current) {
        // Nothing set - will fail.
        return
      }

      let match = this.tournament.currentMatch

      // We don't want to update until the next match has been
      // started. If we do, the graphs are removed as soon as the
      // judges end the previous match.
      // Also, if we're on the first match there is no previous one,
      // so don't try to grab the previous one in that case.
      if (!match.isStarted || this.tournament.current === 0) {
        return this.tournament.matches[this.tournament.current - 1]
      }

      return match
    },
  },
  watch: {
    tournament (val, old) {
      if (val === undefined || old === undefined) {
        return
      }

      // If we're watching a stream-only route, suppress the message.
      if (_.includes([
        "match",
        "scores",
        "next",
        "charts",
        "credits"
      ], this.$route.name)) {
        return
      }

      let $vue = this
      let n = val.events
      let o = old.events
      if (n.length !== o.length) {
        $vue.$info(n[0].print)
      }
    }
  }
}

export default DrunkenFallMixin
