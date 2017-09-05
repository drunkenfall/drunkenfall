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

      <h1>Links</h1>
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
  </div>
</template>

<script>
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
    })
  },
}

</script>

<style lang="scss" scoped>
$button-bg: #404049;
$button-hover-bg: #434350;

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

  .content {
    h1 {
      font-size: 1em;
    }

    .links {
      margin-left: 1em;
      font-size: 1.2em;

      a {
        position: relative;
        background-color: $button-bg;
        text-shadow: 2px 2px 2px rgba(0,0,0,0.3);
        width: 100%;
        display: block;
        margin-bottom: 0.5em;
        transition: 0.3s;

        &:hover {
          background-color: $button-hover-bg;
        }

        .tooltip {
          visibility: hidden;
          opacity: 0;
          position: absolute;
          top: 0.4em;
          left: 105%;
          z-index: 1;

          transition: opacity .15s, visibility .15s;

          color: #dbdbdb;
          background-color: #504040;
          box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

          white-space: nowrap;
          font-size: 0.8em;
        }

        &.router-link-active p {
          color: #708090;
        }

        &.disabled {
          color: #555;
          background-color: #373739;

          .icon {
            background-color: #404040 !important;
          }

          &:hover .tooltip {
            visibility: visible;
            opacity: 1;
          }
        }

        .icon, p {
          padding: 0.6em 0.5em 0.2em;
        }

        .icon {
          float: left;
          width: 1em;
          text-align: center;
          background-color: #405060;

          &.positive {background-color: #406040;}
          &.warning {background-color: #774e2e;}
          &.danger {background-color: #772e2e;}

          .fa-icon {
            filter: drop-shadow(1px 1px 1px rgba(0,0,0,0.3));
          }
        }

        p {
          display: inline-block;

        }
      }
    }
  }
}
</style>
