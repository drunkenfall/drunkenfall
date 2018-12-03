<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' / Players - DrunkenFall'"></headful>
    <tournament-controls />

    <player-toggle
      :on="joined" :onLabel="'In for the showdown! |o/'"
      :off="notJoined" :offLabel="'Booooooo 😧'"
      :toggle="toggle">
    </player-toggle>
  </div>
</template>

<script>
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"
import PlayerToggle from "./players/PlayerToggle.vue"

export default {
  name: 'Participants',
  mixins: [DrunkenFallMixin],
  components: {
    PlayerToggle,
    TournamentControls,
  },

  methods: {
    toggle (e) {
      let $vue = this
      let person = e.target
      this.api.toggle({ id: this.tournament.id, person: person.id }).then((res) => {
        console.log("toggle response:", res)
      }, (err) => {
        $vue.$alert("Toggle failed. See console.")
        console.error(err)
      })
    },
  },

  computed: {
    joined () {
      return _.sortBy(_.map(this.playerSummaries, (p) => p.person), ['name'])
    },
    notJoined () {
      let $vue = this

      return _.sortBy(_.filter(this.people, function (o) {
        let p = _.find($vue.playerSummaries, function (p) {
          return p.person_id === o.id
        })
        return p === undefined
      }), "name")
    },
  },

  watch: {
    tournament (nt, ot) {
      if (nt) {
        let $vue = this
        if (this.playerSummaries === undefined) {
          let id = this.tournament.id
          console.log(`Getting players for ${id}`)
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
        }
      }
    }
  },

  created () {
    this.api = this.$resource("/api", {}, {
      toggle: { method: "GET", url: "/api/tournaments/{id}/toggle/{person}" },
    })
  }
}
</script>
