<template>
  <div>
    <header>
      <div class="content">
        <div class="title">
          {{tournament.name}} / {{match.kind | capitalize}} {{match.index +1}}
        </div>
      </div>
      <div class="links">
        <a v-if="match.canStart" @click="start">Start match</a>
        <a v-if="match.isRunning" v-bind:class="{'disabled': !can_commit}" @click="commit">End round</a>
        <a v-if="match.canEnd"@click="end">End match</a>
      </div>
      <div class="clear"></div>
    </header>

    <div class="control">
      <template v-for="player in match.players" v-ref:players>
        <control-player :index="$index" :player="player" :match="match"
                        :downs="0" :ups="0">
      </template>
    </div>
    <div class="clear"></div>
  </div>
</template>

<script>
import ControlPlayer from './ControlPlayer.vue'
import Match from '../models/Match.js'
import Tournament from '../models/Tournament.js'
import _ from 'lodash'

export default {
  name: 'Match',
  components: {
    ControlPlayer
  },

  data () {
    return {
      match: new Match(),
      tournament: new Tournament()
    }
  },

  computed: {
    can_commit: function () {
      return true
    }

  },

  methods: {
    commit: function () {
      let url = `/api/towerfall/tournament/${this.tournament.id}/${this.match.kind}/${this.match.index}/commit/`

      // TODO this could potentially be a class
      let payload = {
        'state': _.map(this.$refs.players, (controlPlayer) => {
          return _.pick(controlPlayer, ['ups', 'downs', 'shot', 'reason'])
        })
      }

      console.log(payload)
      this.$http.post(url, payload).then(function (res) {
        console.log(res)
        this.$set('match', Match.fromObject(res.data.match))

        _.each(this.$refs.players, (controlPlayer) => { controlPlayer.reset() })
      }, function (res) {
        console.log('error when setting score')
        console.log(res)
      })
    },

    refresh: function () {
      // Hax to make vue refresh the entire page.
      // Since nothing on this page is properly bound to components right now
      // the updates won't trigger properly.
      this.$set('updated', Date.now())
    },
    end: function () {
      let url = `/api/towerfall/tournament/${this.tournament.id}/${this.match.kind}/${this.match.index}/toggle/`

      this.$http.get(url).then(function (res) {
        console.log(res)
        this.$route.router.go('/towerfall/' + this.tournament.id + '/')
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    start: function () {
      let url = `/api/towerfall/tournament/${this.$data.tournament.id}/${this.$data.match.kind}/${this.$data.match.index}/toggle/`
      this.$http.get(url).then(function (res) {
        console.log(res)
        this.setData(
          res.data.tournament,
          this.match.kind,
          this.match.index
        )
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    setData: function (tournament, kind, match) {
      if (kind === 'tryout') {
        kind = 'tryouts'
      } else if (kind === 'semi') {
        kind = 'semis'
      }

      if (kind === 'final') {
        this.$set('match', Match.fromObject(tournament[kind]))
      } else {
        this.$set('match', Match.fromObject(tournament[kind][match]))
      }

      this.$set('tournament', Tournament.fromObject(tournament))
    }
  },

  route: {
    data ({ to }) {
      to.router.app.$watch('tournaments', (newVal, oldVal) => {
        let thisTournament = _.find(newVal, { id: to.params.tournament })
        console.debug('Match.vue - watch update', thisTournament)
        this.setData(thisTournament, to.params.kind, parseInt(to.params.match))
      })

      if (to.router.app.$data.tournaments.length === 0) {
        // Nothing is set - we're reloading the page and we need to get the
        // data manually
        this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
          console.log(res)
          this.setData(
            res.data.Tournament,
            to.params.kind,
            parseInt(to.params.match)
          )
        }, function (res) {
          console.log('error when getting tournament')
          console.log(res)
        })
      } else {
        // Something is set - we're clicking on a link and can reuse the
        // already existing data immediately
        this.setData(
          to.router.app.get(to.params.tournament),
          to.params.kind,
          parseInt(to.params.match)
        )
      }
    }
  }
}
</script>

<style lang="scss" >

@import "../style.scss";

.control {
  height: 85vh;
  padding: 0.8%;
}

.player {
  height: 25%;
  display: block;
}

</style>
