<template>
  <div>
    <div v-for="player in tournament.players" class="player {{player.person.color_preference[0]}}">
      <img alt="{{player.person.nick}}" :src="player.person.avatar_url"/>
      <div class="ribbon {{player.person.color_preference[0]}}">
        <strong class="ribbon-content">{{player.person.nick}}</strong>
      </div>

    </div>
  </div>
</template>

<script>
import Tournament from '../models/Tournament'
import User from '../models/User'
import * as levels from "../models/Level"

export default {
  name: 'TournamentPreview',

  props: {
    tournament: new Tournament(),
    user: new User(),
    levels: levels,
  },
}
</script>

<style lang="scss">
@import "../style.scss";

.player {
  position: relative;
  float: left;
  width: 150px;
  margin: 1% 2% 0%;
  height: 20vh;

  /* background-color: rgba(0,0,0,0.3); */

  img {
    object-fit: cover;
    border-radius: 100%;
    width:  150px;
    height: 150px;
  }
  .ribbon {
    width: 88%;
  }

  &.green  img { border: 1px solid $green ; }
  &.blue   img { border: 1px solid $blue  ; }
  &.pink   img { border: 1px solid $pink  ; }
  &.orange img { border: 1px solid $orange; }
  &.white  img { border: 1px solid $white ; }
  &.yellow img { border: 1px solid $yellow; }
  &.cyan   img { border: 1px solid $cyan  ; }
  &.purple img { border: 1px solid $purple; }
  &.red    img { border: 1px solid $red; }

}

</style>
