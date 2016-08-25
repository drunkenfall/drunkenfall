<template>
  <p class="player {{classes}}">
    {{display_name}}
  </p>
</template>

<script>
export default {
  name: 'Player',

  props: {
    player: Object,
    match: Object,
    index: 0
  },

  computed: {
    prefill: function () {
      return this.player === undefined || this.player.name === ''
    },
    display_name: function () {
      if (this.prefill) {
        return '???'
      }
      return this.player.name
    },
    classes: function () {
      if (this.prefill) {
        return 'prefill'
      }

      if (!this.match.isEnded) {
        if (this.index === 0) {
          return 'gold'
        } else if (this.index === 1) {
          return 'silver'
        } else if (this.index === 2 && this.match.kind === 'final') {
          return 'bronze'
        }

        return 'out'
      }

      return this.player.preferred_color
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
