<template>
  <div class="player {{classes}}">
    <img id="{{player.person.id}}" alt="{{player.person.nick}}" src="https://graph.facebook.com/{{player.person.facebook_id}}/picture?width=9999"/>
    <p class="{{player.color}}">{{display_name}}</p>
    <div class="clear"></div>

  </div>
</template>

<script>
import Match from "../models/Match.js"

export default {
  name: 'Player',

  props: {
    player: Object,
    match: {
      coerce: (val) => { return Match.fromObject(val) }
    },
    index: 0
  },

  computed: {
    avatar: function () {
      return this.player.avatar
    },
    display_name: function () {
      return this.player.person.nick
    },
    classes: function () {
      if (this.match.isEnded) {
        if (this.index === this.match.kill_order[0]) {
          return 'gold'
        } else if (this.index === this.match.kill_order[1]) {
          return 'silver'
        } else if (this.index === this.match.kill_order[2] && this.match.kind === 'final') {
          return 'bronze'
        }

        return 'out'
      }

      return this.player.color
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";

.player {
  display: inline-block;
  width: 100px;
  margin-top: 0px;
  cursor: pointer;

  img {
    object-fit: cover;
    border-radius: 100%;
    width:  125px;
    height: 125px;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    margin-bottom: -40px;
    z-index: 10;
  }
  p {
    z-index: 100;
    width: 80%;
    text-align: center;
    padding: 0.05em 0.1em;
    margin: 0.5em auto;
    display: block;
    font-weight: bold;
    text-shadow: 2px 2px 3px rgba(0,0,0,0.9);

    &.green  { color: $green ; }
    &.blue   { color: $blue  ; }
    &.pink   { color: $pink  ; }
    &.orange { color: $orange; }
    &.white  { color: $white ; }
    &.yellow { color: $yellow; }
    &.cyan   { color: $cyan  ; }
    &.purple { color: $purple; }
    &.red    { color: $red; }

  }
}

</style>
