<template>
<div v-if="tournament">
  <headful :title="tournament.name"></headful>
  <div class="top">
    <div class="title subtitle-logo">
      <img alt="" src="/static/img/oem.svg"/>
      <div class="text">
        <p class="header">DrunkenFall</p>
        <p class="subtitle" :class="tournament.color">{{tournament.subtitle}}</p>
      </div>
    </div>
    <div class="time">
      <p>{{clock.time}}</p>
    </div>
  </div>

  <h1>{{match.title}}</h1>

  <div class="timer">
    {{countdown.time}}
  </div>

    <div class="players">
      <template v-for="(player, index) in playersReversed" ref="players">
        <preview-player :index="index + 1" :player="player" :match="match"></preview-player>
      </template>
    </div>
  </div>
</template>

<script>
import PreviewPlayer from './PreviewPlayer.vue'
import {Countdown, Clock} from '../models/Timer'
import _ from 'lodash'
// import moment from 'moment'
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

  methods: {
    setTime (x) {
      // We need an extra two seconds because 1) one interval has to
      // pass 2) by the time it renders the clock a few milliseconds
      // has passed and there is actually less time left.
      // this.countdown.start(moment().add(x, 'minutes').add(2, 'seconds'))
      let $vue = this
      this.api.setTime({ id: this.tournament.id, time: x }).then((res) => {
        console.debug("settime response:", res)
      }, (err) => {
        $vue.$alert("Setting time failed. See console.")
        console.error(err)
      })
    },
    keyPress (e) {
      let code = e.keyCode
      if (code >= 48 && code <= 57) {
        // https://www.cambiaresearch.com/articles/15/javascript-char-codes-key-codes
        this.setTime(code - 48)
      } else {
        console.log(code)
      }
    },
  },

  computed: {
    playersReversed () {
      // Work on a clone, not the original data object.
      return _.reverse(_.map(this.match.players, _.clone))
    },
    tournament () {
      return this.runningTournament
    },
    match () {
      return this.tournament.upcomingMatch
    },
  },

  mounted () {
    document.getElementsByTagName("body")[0].className = "sidebar-less"

    document.addEventListener("keydown", this.keyPress, false)
  },

  watch: {
    tournament (nt, ot) {
      if (nt) {
        console.log("starting clocks")
        this.countdown.start(this.match.scheduled)
        this.clock.start()
      }
    }
  },

  created () {
    this.api = this.$resource("/api", {}, {
      setTime: { method: "GET", url: "/api/{/id}/time{/time}" },
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/fonts.scss";
@import "../css/ribbon.scss";

.top {
  display: flex;
  width: 100%;
  justify-content: space-between;
  flex-direction: row;

  .title {
    width: 580px;
  }

  .time {
    font-size: 1.5em;
    padding: 16px 40px;
    @include button();
  }
}

.players {
  width: 100%;

  .player {
    width: 25%;
    display: block;
    float: left;
  }
}

h1 {
  /* margin-top: -50px; */
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

</style>
