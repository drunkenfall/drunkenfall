<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' / Players - DrunkenFall'"></headful>
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
import PlayerToggle from "./players/PlayerToggle.vue"

export default {
  name: 'Participants',
  mixins: [DrunkenFallMixin],
  components: {
    PlayerToggle
  },

  methods: {
    toggle (e) {
      let $vue = this
      let person = e.srcElement
      this.api.toggle({ id: this.tournament.id, person: person.id }).then((res) => {
        console.log("join response:", res)
      }, (err) => {
        $vue.$alert("Join failed. See console.")
        console.error(err)
      })
    },
  },

  computed: {
    joined () {
      console.log(this.tournament)
      return _.sortBy(_.map(this.tournament.players, (p) => p.person), ['name'])
    },
    notJoined () {
      let $vue = this
      return _.sortBy(_.filter(this.people, function (o) {
        let p = _.find($vue.tournament.players, function (p) {
          return p.person.id === o.id
        })
        return p === undefined
      }), "name")
    },
  },

  created () {
    this.api = this.$resource("/api", {}, {
      toggle: { method: "GET", url: "/api/{id}/toggle/{person}" },
    })
  }
}
</script>
