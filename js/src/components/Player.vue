<template>
  <p class="player {{classes}}">
    {{display_name}}
  </p>
</template>

<script>
export default {
  name: 'Player',
  data: () => {
    return {
      name: '',
      index: 0,
      preferred_color: ''
    }
  },
  computed: {
    display_name: function () {
      return this.prefill === true ? '???' : this.$data.name
    },
    prefill: function () {
      return this.$data.name === ''
    },
    classes: function () {
      return this.prefill === true ? 'prefill' : this.$data.preferred_color
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
