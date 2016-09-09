<template>
  <a v-link="{path: match.kind + '/' + match.index + '/'}">
    <template v-for="player in match.players">
      <player :player="player" :match="match" :index="$index">
    </template>
  </a>
</template>

<script>
import Player from './Player.vue'
import Match from '../models/Match.js'

export default {
  name: 'MatchOverview',
  components: {
    Player
  },
  props: {
    match: {
      coerce: (val) => { return Match.fromObject(val) }
    }
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
