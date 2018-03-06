<template>
<div>
  <div class="level">
    {{match.level | levelTitle}}
  </div>
  <div class="clear"></div>

  <router-link :to="{ name: 'match', params: { match: match.index }}">
    <template v-for="x in 4">
      <player v-if="match.players.length >= x" :player="match.players[x-1]" :match="match" :index="x-1"></player>

      <div v-if="match.players.length < x" class="player prefill">
        <p>???</p>
      </div>
    </template>

  </router-link>
  <div class="clear"></div>


</div>
</template>

<script>
import Player from './Player.vue'

export default {
  name: 'MatchOverview',
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

.match {
  float: left;

  @media screen and ($desktop: $desktop-width) {
    margin: 10px;
  }
  @media screen and ($device: $device-width) {
    margin: 0.3rem 0;
  }

  cursor: pointer;
  text-align: center;

  .player {
    z-index: 100;
    float: left;
    width: 50%;
    height: 40%;
    line-height: 120%;
    text-align: center;
    overflow: hidden;
    white-space: nowrap;
    padding: 0.2em 0;
    font-size: 30px;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);

    &:nth-child(1), &:nth-child(4) {
      background-color: $bg-default;
      img {
        left: -20px;
      }
    }
    &:nth-child(2), &:nth-child(3) {
      background-color: $bg-default-alt;
      img {
        right: -20px;
      }
    }

    img {
      object-fit: cover;
      border-radius: 100%;
      width:  60px;
      height: 60px;
      box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
      background-color: rgba(10,12,14,0.3);
      position: absolute;
      top: 0;
    }

    &.prefill {
      color: #555;
    }

    &.green  { color: $green ; }
    &.blue   { color: $blue  ; }
    &.pink   { color: $pink  ; }
    &.orange { color: $orange; }
    &.white  { color: $white ; }
    &.yellow { color: $yellow; }
    &.cyan   { color: $cyan  ; }
    &.purple { color: $purple; }
    &.red    { color: $red; }

    &.gold {
      background-color: #daa520;
    }
    &.silver {
      background-color: #999;
    }
    &.bronze {
      background-color: #8C7853;
    }
    &.out {
      color: #777;
      text-shadow: 1px 1px 5px rgba(0,0,0,0.3);
    }
    &.out:nth-child(1), &.out:nth-child(4) {
      background-color: #433;
    }
    &.out:nth-child(2), &.out:nth-child(3) {
      background-color: #4a3838;
    }

  }

  img {
    width: 100%;
  }

  .level {
    background-color: $bg-default-hover;
    padding: 0.3em 0.5em;
    text-align: center;
    margin: 0 auto;
    font-size: 0.9em;
    height: 1em;
    width: 60%;
  }
}
</style>
