<template>
  <div v-if="tournament">
    <h1>
      Starting soon
    </h1>

    <div class="players">
      <div v-for="player in tournament.players" class="player">
        <img :alt="player.person.nick" :src="player.avatar"/>
      </div>
      <div class="clear"></div>
    </div>

    <div class="protector">
      <div class="super-ribbon">
        drunkenfall.com
      </div>

      <div class="ribbon">
        <strong class="ribbon-content">
          {{ countdown.time }}
        </strong>
      </div>
    </div>
  </div>
</template>

<script>
import {Countdown} from '../models/Timer.js'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'TournamentPreview',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      countdown: new Countdown(),
    }
  },

  watch: {
    tournament (nt, ot) {
      if (nt) {
        console.log("starting clock")
        this.countdown.start(nt.scheduled)
      }
    }
  },
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";
@import "../ribbon.scss";

h1 {
  font-size: 6em;
  margin-top: 0;
  margin-bottom: 0.4em;
  padding-top: 0.2em;
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
.super-ribbon {
  margin: -3em auto 2.5em;
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
      box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
      background-color: rgba(10,12,14,0.3);
    }
    .ribbon {
      width: 88%;
    }
  }
}
</style>
