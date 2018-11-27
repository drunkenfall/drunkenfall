<template>
<div v-if="match">
  <template v-for="(p, x) in match.players">
    <player :player="p" :match="match" :index="x"></player>
  </template>

  <div class="details">
    <div class="level">
      <p class="title">
        {{match.title}}
      </p>
      <p class="map">
        {{match.levelTitle}}
      </p>
    </div>

    <div class="clock">
      <span>Starts in</span>
      <p>00:00</p>
    </div>
  </div>
</div>
</template>

<script>
import Player from './Player.vue'
export default {
  name: 'Match',
  components: {
    Player
  },
  props: {
    match: {}
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.details {
  display: flex;
  flex-direction: row;
  flex-grow: 1;
  text-align: center;
  flex-basis: 0;
  justify-content: space-around;
  align-items: center;

  .level {
    font-size: 2em;
    .title {
      color: $fg-secondary;
    }
    .map {
      font-size: 2em;
    }
  }
  .clock {
    p {
      font-size: 4em;
    }
    span {
      font-size: 2em;
      color: $fg-secondary;
    }
  }
}
</style>
