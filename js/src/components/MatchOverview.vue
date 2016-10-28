<template>
  <a v-link="{ name: 'match', params: { kind: match.kind, match: match.index }}">
    <template v-for="x in 4">
      <player v-if="match.players.length > x" :player="match.players[x]" :match="match" :index="x"></player>

      <div v-if="match.players.length <= x" class="player prefill">
        <p>???</p>
      </div>
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
@import "../variables.scss";

.match {
  float: left;
  margin: 10px;
  cursor: pointer;
  text-align: center;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
  background-color: #333339;

  .player {
    float: left;
    width: 25%;
    height: 40%;
    line-height: 120%;
    text-align: center;
    overflow: hidden;
    white-space: nowrap;
    padding: 0.2em 0;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);
  }
  img {
    width: 100%;
  }
}
</style>
