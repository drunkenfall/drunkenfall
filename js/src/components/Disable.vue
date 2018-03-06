<template>
  <div v-if="people">
    <headful title="Disable - DrunkenFall"></headful>
    <player-toggle
      :on="enabled" :onLabel="'Still fightin\'!'"
      :off="disabled" :offLabel="'...their deeds of valor will be remembered.'"
      :toggle="toggle">
    </player-toggle>
  </div>
</template>

<script>
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"
import PlayerToggle from "./players/PlayerToggle.vue"

export default {
  name: 'Disable',
  mixins: [DrunkenFallMixin],
  components: {
    PlayerToggle
  },

  methods: {
    toggle (e) {
      let $vue = this
      let person = e.target
      this.api.toggle({ person: person.id }).then((res) => {
        $vue.$store.commit('setPeople', JSON.parse(res.data))
      }, (err) => {
        $vue.$alert("Disable failed. See console.")
        console.error(err)
      })
    },
  },

  computed: {
    enabled () {
      return _.sortBy(_.filter(this.people, (p) => !p.disabled), ['name'])
    },
    disabled () {
      return _.sortBy(_.filter(this.people, (p) => p.disabled), ['name'])
    },
  },

  created () {
    this.api = this.$resource("/api", {}, {
      toggle: { method: "GET", url: "/api/user/{person}/disable" },
    })
  }
}
</script>
