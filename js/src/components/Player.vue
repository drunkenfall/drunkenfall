<template>
  <div class="player {{classes}}">
    <p>{{display_name}}</p>
  </div>
</template>

<script>
import Match from "../models/Match.js"

export default {
  name: 'Player',

  props: {
    player: Object,
    match: {
      coerce: (val) => { return Match.fromObject(val) }
    },
    index: 0
  },

  computed: {
    avatar: function () {
      return this.player.avatar
    },
    display_name: function () {
      return this.player.person.nick
    },
    classes: function () {
      if (this.match.isEnded) {
        if (this.index === 0) {
          return 'gold'
        } else if (this.index === 1) {
          return 'silver'
        } else if (this.index === 2 && this.match.kind === 'final') {
          return 'bronze'
        }

        return 'out'
      }

      return this.player.person.color_preference[0]
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
