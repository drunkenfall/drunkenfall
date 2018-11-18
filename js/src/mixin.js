import _ from "lodash"
import moment from "moment"

var DrunkenFallMixin = {
  created () {
  },
  methods: {
    $alert (msg) {
      console.error(msg)

      // this.$toast(msg, {
      //   className: ['event', 'alert'],
      //   horizontalPosition: 'right',
      //   verticalPosition: 'bottom',
      //   duration: 5000,
      //   mode: 'queue',
      // })
    },
    $warn (msg) {
      console.warn(msg)

      // this.$toast(msg, {
      //   className: ['event', 'warning'],
      //   horizontalPosition: 'right',
      //   verticalPosition: 'bottom',
      //   duration: 5000,
      //   mode: 'queue',
      // })
    },
    $info (msg) {
      console.log(msg)

      // this.$toast(msg, {
      //   className: ['event'],
      //   horizontalPosition: 'right',
      //   verticalPosition: 'bottom',
      //   duration: 5000,
      //   mode: 'queue',
      // })
    },
    getOrdinal (n) {
      var s = ["th", "st", "nd", "rd"]
      var v = n % 100
      return s[(v - 20) % 10] || s[v] || s[0]
    },
    ordinal (n) {
      return `${n}${this.getOrdinal(n)}`
    },
    arrowImage (n) {
      let arrows = [
        "singleArrow",
        "bombArrows",
        "superBombArrows",
        "laserArrows",
        "brambleArrows",
        "drillArrows",
        "boltArrows",
        "toyArrows",
        "featherArrows",
        "triggerArrows",
        "prismArrows",
      ]

      return `/static/img/arrows/${arrows[n]}.png`
    },
    shieldImage () {
      return `/static/img/arrows/shield.png`
    },
    lavaOrbImage () {
      return `/static/img/arrows/lavaOrb.png`
    },
  },

  computed: {
    tournament () {
      return this.tournaments[this.$route.params.tournament]
    },
    nextTournament () {
      return _.head(this.$store.getters.upcoming)
    },
    tournaments () {
      return this.$store.state.tournaments
    },
    upcomingTournament () {
      let up = this.$store.getters.upcoming
      if (up) {
        return up[0]
      }
    },
    latestTournament () {
      return this.$store.getters.latest
    },
    runningTournament () {
      return this.$store.getters.running
    },
    trackingTournament () {
      if (this.runningTournament) {
        return this.runningTournament
      } else if (this.upcomingTournament) {
        return this.upcomingTournament
      } else if (this.latestTournament && this.latestTournament.endedRecently) {
        return this.latestTournament
      }
    },
    currentLeague () {
      let ts = _.filter(this.tournaments, (t) => {
        return t.scheduled.year() === moment().year() && !t.isTest
      })
      return _.reverse(_.sortBy(ts, 'scheduled'))
    },
    user () {
      return this.$store.state.user
    },
    userLoaded () {
      return this.$store.state.userLoaded
    },
    people () {
      return this.$store.state.people
    },
    playerSummaries () {
      return this.$store.getters.playerSummaries(this.tournament.id)
    },
    combatants () {
      return _.sortBy(_.filter(this.stats, (p) => {
        return p.total.score > 0
      }), 'rank')
    },
    unfought () {
      return _.sortBy(_.filter(this.stats, (p) => {
        return p.total.score === 0
      }), 'person.displayName')
    },
    stats () {
      return _.filter(this.$store.state.stats, (p) => {
        return !p.person.disabled
      })
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
      return this.tournament.matches[this.tournament.current]
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
    showSidebar () {
      // If the fullscreen GET parameter is set, we should _not_ show
      // the sidebar at all.
      return this.$route.query.fullscreen === undefined
    }
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
