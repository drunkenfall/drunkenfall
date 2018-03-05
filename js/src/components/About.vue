<template>
<div>
  <headful title="DrunkenFall"></headful>
  <div id="start" :class="{loaded: !loading}">
    <div class="main">
      <div class="logo">
        <img alt="DrunkenFall" src="/static/img/oem-text.svg"/>
      </div>

      <div class="about">
        <h1>Wait, what's this?</h1>
        <p>
          DrunkenFall is a monthly video game tournament showdown of
          archery skills. We play
          <a href="http://www.towerfall-game.com" target="_blank">TowerFall</a>
          with hardcore tournament rules, and we top that off with
          punishments whenever you lose a point or otherwise embarrass
          yourself.
        </p>

        <p>
          Our events are streamed on
          <a href="https://twitch.tv/drunkenfallofficial">our Twitch
            channel</a>, and you are most welcome to come shout at the
          players and pressure them to make more mistakes.
        </p>

        <h1>Why would you do this?</h1>
        <p>
          It is awesome fun! We've been at it for years; starting at
          Christmas Day 2013 and escalating ever since! We started out
          with just holding a semblance of a tournament with pen and
          paper, and now we have a streaming platform with an app that
          helps us keep track of who fights who!
        </p>

        <h1>Can I join?</h1>
        <p>
          Oh, we thought you'd never ask! Our motto has always been '<i>the more the
          merrier</i>', and that will always stay true!
        </p>
      </div>
    </div>
    <div class="schedule" v-if="nextTournament">
      <router-link :to="{ name: 'tournament', params: {'tournament': nextTournament.id}}" class="next">
        <div class="title">Next event</div>
        <img class="tournament" :src="nextTournament.cover" />
        <div class="subtitle">{{nextTournament.subtitle}}</div>
        <div class="date">{{nextTournament.scheduled.format("MMMM Do")}}</div>
      </router-link>
      <!-- <div class="upcoming">
      <h2>Future events</h2>
      </div> -->
    </div>
  </div>
</div>
</template>

<script>
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'About',
  mixins: [DrunkenFallMixin],
  computed: {
    loading () {
      return _.isEmpty(this.tournaments)
    },
  }
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

#start {
  display: flex;
  @media screen and ($device: $device-width) {
    flex-direction: column;
  }

  transition: 1.5s;
  opacity: 0;

  &.loaded {
    opacity: 1;
  }

  .main {
    padding: 3% 5%;
    .logo {
      display: flex;

      .name {
        font-size: 9em;
        display: flex;
        align-items: center;
        div {
          margin-left: 0.2em;
        }
      }
    }

    .about {
      h1 {
        margin-top: 1em;
        text-align: left;
      }
      p {
        margin: 1em 0;
        font-family: "Lato";
        line-height: 1.5em;
        /* font-size: 1.4em; */
        a {
          box-shadow: none;
          font-weight: bold;
          color: $accent;
        }
      }
    }
  }

  .schedule {
    padding: 3% 5%;
    text-align: center;
    @media screen and ($desktop: $desktop-width) {
      width: 40%;
    }

    .next {
      .tournament {
        margin: 0 auto;
        background-color: rgba(10,0,0,0.3);
        width: 100%;
        box-shadow: 2px 2px 3px rgba(0,0,0,0.7);
      }

      .title {
        margin: 1em 0 0.5em;
        font-size: 3em;
        color: $fg-disabled;
      }

      .subtitle {
        margin-top: 0.5em;
        font-size: 3em;
      }

      .date {
        font-size: 2em;
      }
    }
    .upcoming {
      h2 {
        margin-top: 1em;
      }
    }
  }
}

</style>
