<template>
  <div v-if="tournament">
    <header>
      <div class="content">
        <div class="title">
          {{match.title}} / Round {{round}}
        </div>
      </div>
      <div class="clear"></div>
    </header>

    <div class="control" v-if="user.isJudge">
      <template v-for="(player, index) in match.players" ref="players">
        <control-player :index="index"></control-player>
      </template>
    </div>

    <footer v-if="user.isJudge">

      <div class="content">
        <div class="title">
          Actions
        </div>
      </div>
      <div class="links" v-if="user.isJudge">
        <a v-if="match.canStart && user.isCommentator" @click="start">
          <div class="icon positive">
            <icon name="play"></icon>
          </div>
          <p>Start match</p>
          <div class="clear"></div>
        </a>

        <a v-if="match.isRunning && match.canEnd" @click="end">
          <div class="icon positive">
            <icon name="check"></icon>
          </div>
          <p>End match</p>
          <div class="clear"></div>
        </a>

        <a v-if="match.isRunning && !match.canEnd" @click="commit">
          <div class="icon warning">
            <icon name="gavel"></icon>
          </div>
          <p>End round</p>
          <div class="clear"></div>
        </a>
      </div>
      <div class="clear"></div>
    </footer>

    <div class="control" v-if="!user.isJudge">
      <template v-if="!match.isStarted" v-for="(player, index) in match.players" ref="players">
        <preview-player :index="index" :player="player" :match="match"></preview-player>
      </template>

      <template v-if="match.isStarted" v-for="(player, index) in match.players" ref="players">
        <live-player :index="index" :player="player" :match="match"></live-player>
      </template>
    </div>

    <div class="clear"></div>
  </div>
</template>

<script>
import ControlPlayer from './ControlPlayer.vue'
import PreviewPlayer from './PreviewPlayer.vue'
import LivePlayer from './LivePlayer.vue'
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Match',
  mixins: [DrunkenFallMixin],
  components: {
    ControlPlayer,
    PreviewPlayer,
    LivePlayer,
  },

  computed: {
    players () {
      return _.filter(this.$children, (o) => {
        return o.$options._componentTag === "control-player"
      })
    }
  },

  methods: {
    commit () {
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
      }, function (res) {
        console.log('error when setting score')
        console.log(res)
      })
    },
    end () {
      this.api.end(this.match_id).then(function (res) {
        console.log(res)
        this.$router.push('/towerfall/' + this.tournament.id + '/')
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    start () {
      this.api.start(this.match_id).then(function (res) {
        console.log(res)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
    reset () {
      this.api.reset(this.match_id).then(function (res) {
        console.log(res)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    },
  },

  created () {
    document.getElementsByTagName("body")[0].className = "scroll-less"

    this.api = this.$resource("/api/towerfall", {}, {
      commit: { method: "POST", url: "/api/towerfall/tournament{/id}{/kind}{/index}/commit/" },
      start: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/start/" },
      end: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/end/" },
      reset: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/reset/" },
    })
  },
}
</script>

<style lang="scss" >

.control {
  height: 85vh;
  padding: 0.8%;
}

.player {
  height: 25%;
  display: block;
}

</style>
