<template>
  <div>
    <h1>
      Starting soon
    </h1>
    <div class="players">
      <div v-for="player in tournament.players" class="player {{player.person.color_preference[0]}}">
        <img alt="{{player.person.nick}}" :src="player.person.avatar_url"/>
      </div>
    </div>

    <div class="protector">
      <div class="super-ribbon">
        drunkenfall.com
      </div>

      <div class="ribbon">
        <strong class="ribbon-content">
          00:43:17
        </strong>

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

<style lang="scss" scoped>
@import "../variables.scss";
@import "../ribbon.scss";

h1 {
  font-size: 8em;
  text-shadow: 5px 5px 10px rgba(0,0,0,0.7);
}

h2 {
  margin: -1.5em 0 1em;
  font-size: 2.5em;
  text-align: center;
  strong {
    text-shadow: 2px 2px 4px rgba(0,0,0,0.7);
    color: #a090a0;
  }
}

.ribbon {
  font-size: 40px;
}

.players {
  text-align: center;
  width: 80%;
  margin: 100px auto;

  .player {
    display: inline-block;
    width: 130px;
    margin-top: -30px;

    img {
      object-fit: cover;
      border-radius: 100%;
      width:  150px;
      height: 150px;
      box-shadow: 5px 5px 3px rgba(0,0,0,0.3);
    }
    .ribbon {
      width: 88%;
    }
  }
}
</style>
