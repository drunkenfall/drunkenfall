<template>
  <div v-if="tournament">
    <h1>
      Starting soon
    </h1>

    <div class="players">
      <transition-group name="list" tag="div">
        <div v-for="player in tournament.players" v-bind:key="player.person.id" class="player">
          <img :alt="player.person.nick" :src="player.avatar"/>
        </div>
      </transition-group>
      <div class="clear"></div>
    </div>

    <div class="protector">
      <div class="super-ribbon">
        drunkenfall.com
      </div>

      <div class="ribbon">
        <strong class="ribbon-content">
          {{ countdown.time }}
        </strong>
      </div>
    </div>

    <div id="join" v-if="userLoaded && user.authenticated">
      <div class="links standalone">
        <a v-if="!isJoined" @click="join">
          <div class="icon positive">
            <icon name="user-plus"></icon>
          </div>
          <p>Join the showdown</p>
          <div class="clear"></div>
        </a>
        <a v-else @click="join">
          <div class="icon warning">
            <icon name="times"></icon>
          </div>
          <p>Leave tournament :(</p>
          <div class="clear"></div>
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import {Countdown} from '../models/Timer.js'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'TournamentPreview',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      countdown: new Countdown(),
    }
  },

  methods: {
    join () {
      let $vue = this
      console.log(this.tournament.id)
      this.api.join({id: this.tournament.id}).then((res) => {
        console.debug("join response:", res)
        let j = res.json()
        this.$router.push('/towerfall' + j.redirect)
      }, (err) => {
        $vue.$alert("Join failed. See console.")
        console.error(err)
      })
    },
  },

  computed: {
    isJoined () {
      return this.tournament.playerJoined(this.user)
    }
  },

  watch: {
    tournament (nt, ot) {
      if (nt) {
        console.log("starting clock")
        this.countdown.start(nt.scheduled)
      }
    }
  },

  created () {
    this.api = this.$resource("/api/towerfall", {}, {
      join: { method: "GET", url: "/api/towerfall/{id}/join/" },
    })
  }
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";
@import "../ribbon.scss";

h1 {
  font-size: 6em;
  margin-top: 0;
  margin-bottom: 0.4em;
  padding-top: 0.2em;
  text-shadow: 5px 5px 10px rgba(0,0,0,0.7);
}

h2 {
  margin: -1.5em 0 1em;
  font-size: 2.5em;
  text-align: center;
  strong {
    text-shadow: 2px 2px 4px rgba(0,0,0,0.7);
    color: #a090a0;
  }
}

.ribbon {
  font-size: 40px;
}
.super-ribbon {
  margin: -3em auto 3em;
}

.list-item {
  display: inline-block;
  margin-right: 10px;
}
.list-enter-active, .list-leave-active {
  transition: all 1s;
}
.list-enter, .list-leave-to {
  opacity: 0;
  width: 0px;
}

.players {
  text-align: center;
  width: 80%;
  margin: 100px auto;

  .player {
    display: inline-block;
    width: 130px;
    margin-top: -30px;

    img {
      object-fit: cover;
      border-radius: 100%;
      width:  150px;
      height: 150px;
      box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
      background-color: rgba(10,12,14,0.3);
    }
    .ribbon {
      width: 88%;
    }
  }
}
</style>
