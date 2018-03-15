<template>
<div v-if="tournament && match && user && user.isJudge">
  <headful :title="match.title + ' / Round ' + round"></headful>
  <div class="control">
    <template v-for="(player, index) in match.players" ref="players">
      <control-player :index="index"></control-player>
    </template>
  </div>

  <div class="bottom">
    <div class="content">
      <p class="title">
        {{match.title}} / Round {{round}} @ {{match.level}}
      </p>
      </div>

      <div class="links">
        <a v-if="match.canStart && user.isCommentator" @click="start"
           :class="{ disabled: tournament.shouldBackfill}">
          <div class="icon positive">
            <icon name="play"></icon>
          </div>
          <p>Start match</p>
          <p class="tooltip">Semis need to be backfilled.</p>
          <div class="clear"></div>
        </a>

        <a v-if="match.isRunning && match.canEnd" @click="end">
          <div class="icon positive">
            <icon name="check"></icon>
          </div>
          <p>End match</p>
          <div class="clear"></div>
        </a>

        <a v-if="match.isRunning && !match.canEnd" @click="commit"  :class="{sending: sending}">
          <div class="icon warning">
            <icon name="gavel"></icon>
          </div>
          <p>End round</p>
          <div class="clear"></div>
        </a>
      </div>
      <div class="clear"></div>
    </div>

    <div class="clear"></div>
  </div>
</template>

<script>
import ControlPlayer from './ControlPlayer.vue'
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'JudgeInterface',
  mixins: [DrunkenFallMixin],
  components: {
    ControlPlayer,
  },

  data () {
    return {
      sending: false,
    }
  },

  computed: {
    players () {
      return _.filter(this.$children, (o) => {
        return o.$options._componentTag === "control-player"
      })
    },
    match () {
      return this.tournament.upcomingMatch
    }
  },

  methods: {
    commit () {
      let $vue = this

      if ($vue.sending) {
        console.log("Avoiding sending...")
        return
      }

      $vue.sending = true

      let data = _.map(this.players, (controlPlayer) => {
        return _.pick(controlPlayer, ['ups', 'downs', 'shot', 'reason'])
      })

      let payload = { 'state': data }

      let hasShots = _.some(data, ['shot', true])
      let hasKills = _.sumBy(data, (o) => { return o.ups }) > 0
      let hasSelfs = _.sumBy(data, (o) => { return o.downs }) < 0

      if (!hasShots && !hasKills && !hasSelfs) {
        console.log("Nothing to commit. Doing nothing.")
        return
      }

      console.log(payload)
      this.api.commit(this.match_id, payload).then(function (res) {
        console.log("Round committed.")
        _.each(this.players, (controlPlayer) => { controlPlayer.reset() })
        $vue.sending = false
      }, function (res) {
        $vue.$alert("Setting score failed. See console.")
        console.error(res)
        $vue.sending = false
      })
    },
    end () {
      this.match.end(this)
    },
    start () {
      this.match.start()
    },
    reset () {
      this.match.reset()
    },
  },

  created () {
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"

    this.api = this.$resource("/api", {}, {
      commit: { method: "POST", url: "/api/tournament{/id}{/index}/commit/" },
    })
  },
}
</script>

<style lang="scss" >
@import "../css/colors.scss";

.control {
  height: 80vh;
  padding: 0.8%;
}

.player {
  @include display1();
  height: 25%;
  display: block;
}

.bottom {
  .content {
    @include button();
    font-size: 0.9em;
    font-weight: bold;
    margin: 10px;
    padding: 0.5em 1.4em;
    background-color: $bg-default;
    text-shadow: 2px 2px 3px rgba(0,0,0,0.5);
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
    float: left;
  }

  .links {
    @include button();
    float: right;
    a, .action {
      margin: 10px !important;
      float: right;
      display: block;
      font-weight: bold;
      text-align: center;
      text-decoration: none;
      min-width: 3em;
      font-size: 0.7em;

      p {
        min-width: 15px;
      }

      &.disabled {
        background-color: $bg-disabled;
        color: $fg-disabled;
        cursor: default;
      }

      &.danger {
        background-color: #604040;
        color: $fg-default;
        margin-right: 200px !important;
      }

      &.sending p {
        color: $fg-disabled !important;
      }
    }
    .tooltip {
      display: none;
    }
  }
}

</style>
