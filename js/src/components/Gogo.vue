<template>

<div v-if="tournament && match && user && user.isJudge">
  <headful :title="match.title + ' / Round ' + round"></headful>
  <div class="content">
    <p class="title">
      {{match.title}} - Round {{round}} @ {{match.levelTitle}}
    </p>

    <a @click="play" id="gogo">
      <div class="icon">
        <icon name="play"></icon>
      </div>
      <div class="clear"></div>
    </a>

    <div class="subheader" v-if="user.isCommentator && tournament.nextMatch && !tournament.isEnded">
      <div v-if="!tournament.nextMatch.isScheduled">
        <p>
          Pause until
          <span>{{tournament.nextMatch.title}}</span>
        </p>
        <div class="links">
          <a @click="setTime(10)">10 min</a>
          <a @click="setTime(7)">7 min</a>
          <a @click="setTime(5)">5 min</a>
          <a @click="setTime(3)">3 min</a>
        </div>
        <div class="clear"></div>
      </div>
      <div v-if="tournament.nextMatch.isScheduled">
        <p class="center">
          <span>{{tournament.nextMatch.title}}</span> scheduled at
          {{tournament.nextMatch.scheduled.format("HH:mm")}}
        </p>
        <div class="clear"></div>
      </div>
    </div>
  </div>
</div>

</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'GogoInterface',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      sending: false,
    }
  },

  computed: {
    match () {
      return this.tournament.nextMatch
    }
  },

  methods: {
    play () {
      let $vue = this

      let go = document.getElementById("gogo")
      go.className = "disabled"
      setTimeout(function () {
        let go = document.getElementById("gogo")
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
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"

    this.api = this.$resource("/api", {}, {
      play: { method: "POST", url: "/api/tournaments/{id}/play/" },
      setTime: { method: "GET", url: "/api/tournaments/{id}/time/{time}" },
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.content {
  .title {
    @include display1();
    text-align: center;
    margin: 5%;
  }

  a {
    @include button();

    padding: 0.5em 0.7em;
    font-size: 6em;
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
  width: 80%;

  @media screen and ($desktop: $desktop-width) {
    p {
      float: left;
    }
    .links {
      float: right;
      a {
        float: right;
      }
    }
  }
  @media screen and ($device: $device-width) {
    & {
      text-align: center;
      padding: 0.5em;
    }
    .links {
      a:last-child {
        margin-bottom: 1em;
      }
    }
  }

  margin: 30px auto;
  background-color: $bg-default;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
  text-shadow: 2px 2px 3px rgba(0,0,0,0.5);

  p {
    font-size: 2em;
    padding: 0.3em 0.5em;

    span {
      text-transform: capitalize;
    }
  }

  .links {
    a, .action {
      @include button();
      margin: 15px  !important;
      background-color: $button-bg;
      border-left: 3px solid $accent;
      color: $fg-default;
      display: block;
      font-weight: bold;
      padding: 7px 30px;
      text-align: center;
      text-decoration: none;
      margin: 10px auto;
      min-width: 60px;

      padding: 0.5em 0.7em;
    }
  }
}

</style>
