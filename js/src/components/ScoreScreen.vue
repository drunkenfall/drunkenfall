<template>
  <div v-if="tournament">
    <template v-for="(player, index) in nextMatch.players" ref="players">
      <live-player :index="index" :player="player" :match="nextMatch"></live-player>
    </template>
    <div class="clear"></div>
  </div>
</template>

<script>
import LivePlayer from './LivePlayer.vue'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'ScoreScreen',
  mixins: [DrunkenFallMixin],
  components: {
    LivePlayer,
  },
  mounted () {
    document.getElementsByTagName("body")[0].className += "scroll-less sidebar-less"
  },
}
</script>

<style lang="scss" >

control {
  height: 85vh;
  padding: 0.8%;
}

.player {
  height: 25%;
  display: block;
}

</style>
