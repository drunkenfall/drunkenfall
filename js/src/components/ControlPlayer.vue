<template>

<div :class="'player ' + player.preferred_color">
  <div class="button" @click="score(index, -1)">
    <div>
      <icon name="minus"></icon>
    </div>
  </div>

  <div class="shots">
    <div @click="manual_shot" v-bind:class="{'give': shot}">
      <div v-if="shot" class="mark">{{player.shots + 1}}</div>
      <div v-if="!shot" class="mark">{{player.shots}}</div>
      <p v-if="reason !== ''" class="reason">{{reason}}</p>
    </div>
  </div>

  <div :class="'slider ' + player.color" @click="reset">
    <div ref="name">{{player.displayName}}</div>
  </div>

  <div v-if="match.endScore == 10" class="scores">
    <div v-for="n in 10"
      :class="bulletClass(player, index, n)">
      <p>{{bulletDisplay(player, index, n)}}</p>
    </div>
  </div>

  <div v-if="match.endScore == 20" class="scores final">
    <div v-for="n in 20"
      :class="bulletClass(player, index, n)">
      <p>{{bulletDisplay(player, index, n)}}</p>
    </div>
  </div>

  <div class="button" @click="score(index, 1)">
    <div>
      <icon name="plus"></icon>
    </div>
  </div>
</div>

</template>

<script>
import fitText from "../util/fittext.js"
import DrunkenFallMixin from "../mixin"

export default {
  name: 'ControlPlayer',
  mixins: [DrunkenFallMixin],

  props: {
    index: 0,
  },

  data () {
    return {
      shot: false,
      ups: 0,
      downs: 0,
      reason: ''
    }
  },

  computed: {
    player () {
      return this.match.players[this.index]
    },
    match () {
      return this.tournament.upcomingMatch
    }
  },

  methods: {
    bulletClass (player, playerIndex, n) {
      if (n === this.match.endScore && (player.kills + this.ups) > this.match.endScore) {
        return 'overkill'
      } else if (n > player.kills && n <= (player.kills + this.ups)) {
        return 'up'
      } else if (n === player.kills && this.downs === -1) {
        return 'down'
      } else if (n <= player.kills) {
        return 'kill'
      }
      return ''
    },

    bulletDisplay (player, playerIndex, n) {
      // Change the display to +1 for every overkill the player makes
      // at the end of the match.
      if (n === this.match.endScore && (player.kills + this.ups) > this.match.endScore) {
        return "+" + (player.kills + this.ups - this.match.endScore)
      }
      return n
    },

    score (playerIndex, score) {
      if (!this.match.isStarted) {
        return
      }

      if (score === 1 && this.ups < 3) {
        // We can only allow up to three kills per round...
        this.ups += 1
        if (this.ups === 3) {
          this.shot = true
          this.reason = 'sweep'
        }
      } else if (score === -1 && this.downs === 0) {
        // ...and only one suicide.
        this.downs -= 1
        this.shot = true
        this.reason = 'suicide'
      }
    },

    reset () {
      this.ups = 0
      this.downs = 0
      this.shot = false
      this.reason = ''
    },

    manual_shot () {
      if (!this.shot) {
        this.shot = true
        this.reason = 'manual'
      } else {
        this.shot = false
        this.reason = ''
      }
    }
  },

  mounted () {
    fitText(this.$refs.name, this.$refs.name.innerText, "'Teko', sans-serif", 0.8)
  }
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.player {
  width: 100%;
  position: relative;

  >div {
    float: left;
    height: 100%;
  }

  .button, .shots {
    width: 10%;
    position: relative;
    text-shadow: 3px 3px 5px rgba(0,0,0,0.8);

    svg {
      filter: drop-shadow(3px 3px 3px rgba(0,0,0,0.5));
    }

    >div {
      box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

      width: 80%;
      height: 50%;
      background-color: $bg-default;
      cursor: pointer;
      text-align: center;
      vertical-align: middle;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translateX(-50%) translateY(-50%);
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
  .button {
    transition: 0.2s;
    div {
      background-color: $button-bg;
      border-left: 3px solid $accent;
      font-size: 5vh;
      user-select: none;
      -ms-user-select: none;
      -moz-user-select: none;
      -webkit-user-select: none;
      p {
        margin-top: -13%;
      }
    }
  }

  .shots {
    transition: 0.2s;

    >div {
      border-left: 3px solid $accent;
    }
    div {
      transition: 0.2s;
      background-color: $button-bg;
      width: 80%;
      display: block;
      margin: 0 auto 1%;
      user-select: none;
      -ms-user-select: none;
      -moz-user-select: none;
      -webkit-user-select: none;

      >div {
        transition: 0.2s;
      }

      &.give {
        background-color: $accent;
        p {
          color: #fff;
          font-size: 0.3em;
          margin-top: -0.3em;
        }
        >div {
          background-color: $accent;
        }
      }

      &.mark {
        padding-top: 10%;
        font-size: 6vh;
      }
      &.reason {
        margin-top: -4%;
        font-size: 0.5em;
      }
    }
  }

  .slider {
    width: 15%;
    text-shadow: 2px 2px 2px rgba(0,0,0,0.8);
    font-weight: bold;
    font-size: 1.6vw;

    &.green  { color: $green ; }
    &.blue   { color: $blue  ; }
    &.pink   { color: $pink  ; }
    &.orange { color: $orange; }
    &.white  { color: $white ; }
    &.yellow { color: $yellow; }
    &.cyan   { color: $cyan  ; }
    &.purple { color: $purple; }
    &.red    { color: $red; }

    div {
      width: 80%;
      height: 50%;
      font-size: 1.6em;
      cursor: pointer;

      text-align: center;
      vertical-align: middle;
      position: relative;
      top: 50%;
      left: 50%;
      transform: translateX(-50%) translateY(-50%);
      display: flex;
      align-items: center;
      justify-content: center;
      // overflow: hidden;
    }
  }
  .scores {
    // http://stackoverflow.com/questions/6865194/fluid-width-with-equally-spaced-divs
    width: 55%;
    text-align: justify;
    text-justify: distribute;
    // background-color: $bg-default;

    &.final {
      div {
        width: 4%;
        margin: 0 0.5%;
      }
    }

    div {
      // position: relative;
      // top: 25%;
      width: 8%;
      margin: 0 1%;
      height: 8vh;
      vertical-align: top;
      display: inline-block;
      zoom: 1;

      background-color: $bg-default;
      border-radius: 10000px;
      box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

      text-align: center;
      vertical-align: middle;
      position: relative;
      top: 50%;
      transform: translateX(0) translateY(-50%);

      p {
        position: relative;
        top: 50%;
        transform: translateX(0) translateY(-50%);
        font-size: 0.5em;
        color: #666672;
      }

      transition: 0.2s;

      &.kill {
        background-color: $fg-default;
        p {color: $bg-bottom;}
      }
      &.overkill {
        background-color: #daa520;
        text-shadow: 2px 2px 2px rgba(0,0,0,0.5);
        p {color: #fff;}
      }
      &.up {
        background-color: #508850;
        p {color: #fff;}
      }
      &.down {
        background-color: #885050;
        p {color: #fff;}
      }
    }
  }
}
</style>
