<template>
  <div v-if="isVisible">
    <router-link
      :to="{ name: 'tournament', params: { tournament: tournament.id }}">

      <div class="title">
        <div class="numeral" :style="numeralColor">{{tournament.numeral}}</div>
        <div class="subtitle">{{tournament.subtitle}}</div>
        <div class="clear"></div>
      </div>
    </router-link>

    <div v-if="isSelected" class="content">
      <div v-if="viewing(['tournament', 'edit', 'log', 'participants'])">
        <div class="matches links">
          <div class="items">
            <router-link class="action" v-if="user.isJudge"
              :to="{ name: 'log', params: { tournament: tournament.id }}">
              <div class="icon">
                <icon name="book"></icon>
              </div>
              <p>Log</p>
              <div class="clear"></div>
            </router-link>

            <router-link class="action" v-if="user.isProducer"
              :to="{ name: 'participants', params: { tournament: tournament.id }}">
              <div class="icon"
                :class="{ warning: tournament.isStarted }">
                <icon name="users"></icon>
              </div>
              <p>Participants</p>
              <div class="clear"></div>
            </router-link>

            <router-link class="action" v-if="user.isProducer && tournament.canShuffle"
              :to="{ name: 'edit', params: { tournament: tournament.id }}">
              <div class="icon danger">
                <icon name="pencil"></icon>
              </div>
              <p>Edit</p>
              <div class="clear"></div>
            </router-link>
          </div>
        </div>
      </div>

      <div v-if="viewing(['tournament'])">
        <h1>Actions</h1>
        <div class="actions links">
          <a v-if="tournament.canStart && user.isCommentator" @click="start">
            <div class="icon positive">
              <icon name="play"></icon>
            </div>
            <p>Start</p>
            <div class="clear"></div>
          </a>

          <a class="action" @click="next" v-if="user.isJudge && tournament.isRunning">
            <div class="icon positive">
              <icon name="play"></icon>
            </div>
            <p>Next match</p>
            <div class="clear"></div>
          </a>

          <a class="action" @click="reshuffle" v-if="user.isProducer && tournament.canShuffle">
            <div class="icon warning">
              <icon name="random"></icon>
            </div>
            <p>Reshuffle</p>
            <div class="clear"></div>
          </a>

          <a class="action" @click="usurp"
            :class="{ disabled: !tournament.isUsurpable}"
            v-if="user.isProducer && tournament.isTest && tournament.canStart">
            <div class="icon warning">
              <icon name="user-plus"></icon>
            </div>
            <p>Add testing players</p>
            <p class="tooltip">Tournament is full.</p>
            <div class="clear"></div>
          </a>
        </div>
      </div>

      <div class="actions links">
        <div v-if="match">
          <a v-if="match.canStart && user.isCommentator" @click="startMatch">
            <div class="icon positive">
              <icon name="play"></icon>
            </div>
            <p>Start match</p>
            <div class="clear"></div>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import _ from "lodash"
import Tournament from '../models/Tournament.js'

export default {
  name: 'SidebarTournament',

  props: {
    tournament: new Tournament()
  },

  computed: {
    user () {
      return this.$store.state.user
    },
    match () {
      if (!this.viewing(["match"])) {
        return undefined
      }

      let kind = this.$route.params.kind
      let idx = this.$route.params.match

      if (kind === 'final') {
        return this.tournament.final
      }
      kind = kind + 's'
      return this.tournament[kind][idx]
    },
    numeralColor () {
      return "background-color: #405060;"
    },
    isVisible () {
      if (this.$route.name === 'start') {
        return true
      }
      return this.isSelected
    },
    isSelected () {
      return this.$route.params.tournament === this.tournament.id
    }
  },

  methods: {
    viewing (names) {
      // Returns true if currently viewing any of the route names.
      return _.includes(names, this.$route.name)
    },
    start () {
      this.api.start({ id: this.tournament.id }).then((res) => {
        console.log("start response:", res)
        let j = res.json()
        this.$route.router.push('/towerfall' + j.redirect)
      }, (err) => {
        console.error(`start for ${this.tournament} failed`, err)
      })
    },
    join () {
      this.api.join({ id: this.tournament.id }).then((res) => {
        console.log("join response:", res)
        var j = res.json()
        this.$route.router.push('/towerfall' + j.redirect)
      }, (err) => {
        console.error(`joining tournament ${this.tournament} failed`, err)
      })
    },
    next () {
      this.api.next({ id: this.tournament.id }).then((res) => {
        console.debug("next response:", res)
        let j = res.json()
        this.$router.push('/towerfall' + j.redirect)
      }, (err) => {
        console.error(`next for ${this.tournament} failed`, err)
      })
    },
    reshuffle () {
      this.api.reshuffle({ id: this.tournament.id }).then((res) => {
        console.debug("reshuffle response:", res)
      }, (err) => {
        console.error(`reshuffle for ${this.tournament} failed`, err)
      })
    },
    usurp () {
      this.api.usurp({ id: this.tournament.id }).then((res) => {
        console.log("usurp response:", res)
      }, (err) => {
        console.error(`usurp for ${this.tournament} failed`, err)
      })
    },
  },

  created () {
    this.api = this.$resource("/api/towerfall", {}, {
      start: { method: "GET", url: "/api/towerfall{/id}/start/" },
      usurp: { method: "GET", url: "/api/towerfall{/id}/usurp/" },
      next: { method: "GET", url: "/api/towerfall{/id}/next/" },
      reshuffle: { method: "GET", url: "/api/towerfall{/id}/reshuffle/" },

      commitMatch: { method: "POST", url: "/api/towerfall/tournament{/id}{/kind}{/index}/commit/" },
      startMatch: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/start/" },
      endMatch: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/end/" },
      resetMatch: { method: "GET", url: "/api/towerfall/tournament{/id}{/kind}{/index}/reset/" },
    })
  },
}

</script>

<style lang="scss" scoped>
@import "../variables.scss";

.content {
  margin-top: 1.5rem;
}

.tournament {
  margin-bottom: 0.5em;

  .title {
    cursor: pointer;
    font-size: 1.3em;
    background-color: $button-bg;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
    text-shadow: 2px 2px 2px rgba(0,0,0,0.3);

    overflow: hidden;
    white-space: nowrap;

    .numeral, .subtitle {
      padding: 0.6em 0.5em 0.3em;
    }

    .numeral {
      float: left;
      min-width: 1.2em;
      text-align: center;
    }
    .subtitle {
      display: inline-block;
      /* padding-left: 2.6em; */
    }
  }
}
</style>
