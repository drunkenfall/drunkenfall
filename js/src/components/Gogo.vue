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
      return this.tournament.upcomingMatch
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
    }
  },

  created () {
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"

    this.api = this.$resource("/api", {}, {
      play: { method: "POST", url: "/api/tournaments/{id}/play/" },
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

</style>
