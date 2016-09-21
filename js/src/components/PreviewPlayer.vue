<template>
  <div class="player {{player.color}}">
    <img alt="{{player.displayName}}" :src="player.avatar"/>

    <div class="protector">
      <div class="super-ribbon">
      {{player.firstName}}
      </div>

      <div class="ribbon {{player.color}}">
        <strong class="ribbon-content">
          {{player.displayName}}
        </strong>
      </div>
    </div>
  </div>
</template>

<script>
import Match from "../models/Match.js"

export default {
  name: 'ControlPlayer',

  props: {
    player: {},
    match: {
      coerce: (val) => { return Match.fromObject(val) }
    },
    index: 0,
  },

  computed: {
    classes: function () {
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
  },
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";
@import "../ribbon.scss";

.player {
  float: left;
  height: 50%;
  width: 50%;
  margin-top: 4em;

  img {
    display: block;
    object-fit: cover;
    border-radius: 100%;
    width:  350px;
    height: 350px;
    margin: 0px auto;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    background-color: rgba(10,12,14,0.3);
  }

  .super-ribbon {
    margin: -3.5em auto 2.2em;
  }
  .ribbon {
    font-size: 2em;
    width: 45%;
  }
}

</style>
