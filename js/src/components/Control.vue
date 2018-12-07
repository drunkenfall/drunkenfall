<template>

<div v-if="tournament && match && user && user.isJudge">
  <headful :title="match.title + ' / Round ' + round"></headful>
  <tournament-controls />

  <div class="content">
    <p class="title">
      {{match.title}} - Round {{round}} @ {{match.levelTitle}}
    </p>

    <a @click="play" id="god">
      <div class="icon">
        <icon name="play"></icon>
      </div>
    </a>

    <div class="subheader" v-if="match && !match.isStarted">
      <div v-if="!match.isScheduled">
        <p>
          Pause until
          <span>{{match.title}}</span>
        </p>
        <div class="links">
          <a @click="setTime(x)" v-for="x in pauses">{{x}} min</a>
        </div>
      </div>
      <p v-else>
        Scheduled at {{match.scheduled.format("HH:mm")}}
      </p>
    </div>

    <template v-for="(p, x) in match.players">
      <player :player="p" :match="match" :index="x"></player>
    </template>

  </div>
</div>

</template>

<script>
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"
import Player from './Player.vue'

export default {
  name: 'Control',
  mixins: [DrunkenFallMixin],
  components: {
    TournamentControls,
    Player,
  },

  data () {
    return {
      sending: false,
      pauses: [2, 3, 5, 7, 10],
    }
  },

  computed: {
    match () {
      return this.tournament.currentMatch
    },
  },

  methods: {
    play () {
      let $vue = this

      let go = document.getElementById("god")
      go.className = "disabled"
      setTimeout(function () {
        let go = document.getElementById("god")
        go.className = ""
      }, 3000)

      console.log("Sending startplay match...", this.match_id)
      this.api.play(this.match_id, {}).then((res) => {
        console.log("Match started.", res)
      }, (res) => {
        $vue.$alert("Starting failed. See console.")
        console.error(res)
      })
    },
    setTime (x) {
      let $vue = this
      this.api.setTime({ id: this.tournament.id, time: x }).then((res) => {
        console.debug("settime response:", res)
      }, (err) => {
        $vue.$alert("Setting time failed. See console.")
        console.error(err)
      })
    },
  },

  created () {
    this.loadAll()

    this.api = this.$resource("/api", {}, {
      play: { method: "POST", url: "/api/tournaments/{id}/play/" },
      setTime: { method: "GET", url: "/api/tournaments/{id}/time/{time}" },
    })
  },
}
</script>

<style lang="scss">
@import "../css/colors.scss";

.content {
  padding: 1em;
  .title {
    @include headline();
    text-align: center;
  }

  #god {
    @include button();

    padding: 0.5em 0.7em;
    font-size: 3em;
    margin: 5%;
    display: block;
    font-weight: bold;
    text-align: center;
    text-decoration: none;
    background-color: $positive;
    transition: 0.35s;

    &.disabled {
      background-color: $bg-disabled;
      color: $fg-disabled;
    }
  }
}

.subheader {
  @include subheading();

  text-align: center;
  display: flex;
  flex-direction: row;
  margin: 25px auto;
  background-color: $bg-default;
  padding: 0.3em 0.5em;

  p {
    font-size: 2em;
    padding: 0.3em 0.5em;
    text-align: center;
    flex-grow: 1;
  }

  .links {
    display: flex;
    flex-basis: content;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-evenly;

    a {
      @include button();
      background-color: $button-bg;
      border-left: 3px solid $accent;
      color: $fg-default;
      font-weight: bold;
      text-align: center;
      text-decoration: none;
      width: 6em;

      margin: 2%;
      padding: 0.5em 0.7em;
    }
  }
}

.players {
  img {
    height: 5vh;
  }
}

</style>
