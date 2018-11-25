<template>
  <div v-if="tournament" class="preview">
    <headful :title="tournament.subtitle + ' - DrunkenFall'"></headful>
    <tournament-controls />

    <h1>
      {{tournament.name}}
    </h1>
    <h3>
      {{tournament.scheduled.local().format("ddd MMMM Do HH:mm")}}
    </h3>


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
          <p>Leave tournament</p>
          <div class="clear"></div>
        </a>
      </div>
    </div>

    <div class="players">
      <transition-group name="list" tag="div">
        <div v-for="player in playerSummaries" v-bind:key="player.person.id" class="player">
          <img :alt="player.person.nick" :src="player.avatar"/>
        </div>
      </transition-group>
      <div class="clear"></div>
    </div>

    <div class="protector">
      <!-- <div class="super-ribbon">
        drunkenfall.com
      </div> -->

      <div class="ribbon">
        <strong v-if="tournament.isToday" class="ribbon-content">
          {{ countdown.time }}
        </strong>
        <strong v-else class="ribbon-content">
          {{tournament.scheduled.local().format("ddd MMMM Do HH:mm")}}
        </strong>
      </div>
    </div>
  </div>
</template>

<script>
import {Countdown} from '../models/Timer.js'
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"

export default {
  name: 'TournamentPreview',
  mixins: [DrunkenFallMixin],
  components: {
    TournamentControls,
  },

  data () {
    return {
      countdown: new Countdown(),
    }
  },

  methods: {
    join () {
      let $vue = this

      this.api.join({id: this.tournament.id}).then((res) => {
        console.debug("join response:", res)
        this.$router.push({name: "tournament", params: {id: this.tournament.id}})
      }, (err) => {
        $vue.$alert("Join failed. See console.")
        console.error(err)
      })
    },
  },

  computed: {
    isJoined () {
      return this.tournament.playerJoined(this.user)
    },
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
    let $vue = this
    this.api = this.$resource("/api", {}, {
      join: { method: "GET", url: "/api/tournaments/{id}/join/" },
    })

    console.log(this.tournament)

    if (this.tournament.isToday) {
      console.log("starting clock")
      this.countdown.start(this.tournament.scheduled)
    }

    let id = this.tournament.id
    this.$http.get(`/api/tournaments/${id}/players/`).then(function (res) {
      let data = JSON.parse(res.data)
      this.$store.commit('setPlayerSummaries', {
        tid: id,
        player_summaries: data.player_summaries,
      })
    }, function (res) {
      $vue.$alert("Getting players failed. See console.")
      console.error(res)
    })
  }
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/ribbon.scss";

.preview {
  max-height: 920px;
}

h1 {
  @media screen and ($desktop: $desktop-width) {
    font-size: 6em;
    text-shadow: 5px 5px 10px rgba(0,0,0,0.7);
  }
  @media screen and ($device: $device-width) {
    font-size: 2.5em;
    text-shadow: 3px 3px 6px rgba(0,0,0,0.7);
  }
  margin-top: 0 !important;
  margin-bottom: 0.4em;
  padding-top: 0.2em;
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

.protector {
  margin-top: 7em;
}

@media screen and ($desktop: $desktop-width) {
  .ribbon {
    font-size: 40px;
  }
  .super-ribbon {
    margin: -3em auto 3em;
  }
}
@media screen and ($device: $device-width) {
  .ribbon {
    display: none;
  }
  .super-ribbon {
    display: none;
  }
}

.list-item {
  display: inline-block;
}
.list-enter-active, .list-leave-active {
  transition: all 1s ease-out;
  img {
    transition: all 1s ease-out;
  }
}
.list-enter, .list-leave-to {
  opacity: 0;
  width: 0px;
  height: 0px;
  img {
    width: 0px !important;
    height: 0px !important;
  }
}

.players {
  text-align: center;
  width: 80%;
  margin: 4em auto;

  @media screen and ($desktop: $desktop-width) {
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

    .protector {
      font-size: 3em !important;
      .ribbon {
      }
    }
  }

  @media screen and ($device: $device-width) {
    .player {
      display: inline-block;
      width: 5em;
      margin-top: -30px;

      img {
        object-fit: cover;
        border-radius: 100%;
        width:  6em;
        height: 6em;

        box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
        background-color: rgba(10,12,14,0.3);
      }
      .ribbon {
        width: 88%;
      }
    }
  }
}
</style>
