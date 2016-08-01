<template>
  <a v-link="{path: match.kind + '/' + match.index + '/'}">
    <player v-for="player in match.players" :player="player"></player>
  </a>
</template>

<script>
import Player from './Player.vue'
export default {
  name: 'Match',
  components: {
    Player
  },
  props: {
    match: Object
  }
}
</script>

<style lang="scss" scoped>
.match {
  float: left;
  margin: 10px;
  cursor: pointer;
  text-align: center;

  img {
    width: 100%;
  }
}
</style>
