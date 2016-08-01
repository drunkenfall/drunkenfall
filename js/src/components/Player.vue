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
      return this.player.preferred_color
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
