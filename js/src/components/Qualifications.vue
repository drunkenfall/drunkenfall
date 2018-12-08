<template>

<div v-if="tournament" class="main">
  <headful :title="tournament.subtitle + ' - DrunkenFall'"></headful>
  <div class="main">
    <div v-for="(p, idx) in leaderboard" class="player">
      <div class="index">
        {{idx+1}}.
      </div>
      <div class="avatar">
        <img :alt="p.nick" :src="p.avatar"/>
      </div>

      <p class="nick" :class="p.person.color">
        {{p.displayName}}
      </p>

      <div class="scores">
        {{p.skill_score}}pts,
        {{p.matches}}m
      </div>
    </div>
  </div>
</div>

</template>

<script>
import DrunkenFallMixin from "../mixin"
// import Person from "../models/Person"
import _ from 'lodash'

export default {
  name: 'Qualifications',
  mixins: [DrunkenFallMixin],

  computed: {
    leaderboard () {
      return _.reverse(_.sortBy(this.playerSummaries, 'skill_score').slice(0, 16))
    },
  },

  created () {
    this.loadAll()
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.main {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow: hidden;
  height: 100vh;

  .player {
    width: 100%;
    height: 6.25vh;

    &:nth-child(odd) {
      background-color: rgba(0,0,0,0.1)
    }

    display: flex;
    flex-direction: row;

    .index {
      font-size: 3vh;
      padding: 1.8vh;
      opacity: 0.6;
      width: 4.5vh;
    }

    .avatar {
      max-height: 100%;
      padding: 0.625vh;

      img {
        height: 5vh;
        width:  5vh;
        object-fit: cover;
        border-radius: 100%;
      }
    }

    .nick {
      font-size: 4vh;
      padding: 1.2vh;
      flex-grow: 1;
    }

    .scores {
      padding: 1.6vh;
      font-size: 3vh;
    }
  }
}

</style>
