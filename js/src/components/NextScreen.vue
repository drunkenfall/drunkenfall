<template>
  <div v-if="tournament">
    <header>
      <div class="content">
        <div class="title">
          {{tournament.name}} / {{nextMatch.title}}
        </div>
      </div>
      <div class="links">
        <p class="time">
          Local time: {{clock.time}}
        </p>
      </div>
      <div class="clear"></div>
    </header>

    <h1>Starting in</h1>
    <div class="timer">
      {{countdown.time}}
    </div>

    <div class="players">
      <template v-for="(player, index) in playersReversed" ref="players">
        <preview-player :index="index + 1" :player="player" :match="nextMatch"></preview-player>
      </template>
    </div>
  </div>
</template>

<script>
import PreviewPlayer from './PreviewPlayer.vue'
import {Countdown, Clock} from '../models/Timer'
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'NextScreen',
  mixins: [DrunkenFallMixin],
  components: {
    PreviewPlayer,
  },

  data () {
    return {
      countdown: new Countdown(),
      clock: new Clock(),
    }
  },

  computed: {
    playersReversed () {
      // Work on a clone, not the original data object.
      return _.reverse(_.map(this.nextMatch.players, _.clone))
    },
  },

  mounted () {
    document.getElementsByTagName("body")[0].className = "sidebar-less"
  },

  watch: {
    tournament (nt, ot) {
      if (nt) {
        console.log("starting clocks")
        this.countdown.start(this.nextMatch.scheduled)
        this.clock.start()
      }
    }
  },

}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/ribbon.scss";

.players {
  width: 100%;

  .player {
    width: 25%;
    display: block;
    float: left;
  }
}

h1 {
  margin-top: 100px;
  margin-bottom: -1em;
}

.timer {
  margin: 0.5em auto 0.25em;
  width: 3em;
  font-size: 12em;
  text-align: center;
  padding: 0.08em 0.4em;

  text-shadow: 3px 3px 3px rgba(0,0,0,0.7);
}

.time {
  font-size: 1.5em;
  padding: 16px 40px;
}

</style>
