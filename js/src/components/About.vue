<template>
<div>
  <headful title="DrunkenFall"></headful>
  <div id="start" :class="{loaded: !loading}">
    <div class="main">
      <div class="hero">
        <div class="image"></div>
        <div class="logo">
          <img alt="DrunkenFall" src="/static/img/oem-text.svg"/>
        </div>
      </div>

      <div class="about">
        <div class="text">
          <h1>TowerFall tournaments</h1>
          <p>
            DrunkenFall is a video game tournament showdown of archery skills! We play
            <a href="http://www.towerfall-game.com" target="_blank">TowerFall</a> with hardcore
            tournament rules, and we top that off with punishments whenever you lose a point or
            otherwise embarrass yourself.
          </p>

          <p>
            Our events are streamed on <a href="https://twitch.tv/drunkenfallofficial">our Twitch
            channel</a>, and you are most welcome to come shout at the players and pressure them to
            make more mistakes.
          </p>
        </div>

        <div class="video">
          <h1>Trailer video</h1>
          <iframe
            src="https://www.youtube.com/embed/-VZN6F0eo8c"
            frameborder="0"
            allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen></iframe>
        </div>
      </div>
    </div>

    <div class="schedule" v-if="nextTournament">
      <router-link :to="{ name: 'tournament', params: {'tournament': nextTournament.id}}" class="next">
        <div class="title">Next event</div>
        <img class="tournament" :src="nextTournament.cover" />
        <div class="subtitle">{{nextTournament.subtitle}}</div>
        <div class="date">{{nextTournament.scheduled.format("MMMM Do")}}</div>
      </router-link>
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

@keyframes logo-rotate {
  0% {   transform: rotate(0deg);}
  11% {  transform: rotate(1deg);}
  33% {  transform: rotate(0deg);}
  50% {  transform: rotate(-3deg);}
  75% {  transform: rotate(2deg);}
  100% { transform: rotate(0deg);}
}

#start {
  display: flex;
  @media screen and ($device: $device-width) {
    flex-direction: column;
  }

  opacity: 0;

  &.loaded {
    opacity: 1;
  }

  .main {
    width: 100%;

    .hero {
      height: 33vh;
      position: relative;
      box-shadow: 10px 5px 10px rgba(0,0,0,0.3);

      .image {
        height: 100%;
        width: 100%;
        background-image: url("/static/img/hero.jpg");
        filter: grayscale(66%) brightness(50%);
        background-position: center;
        background-repeat: no-repeat;
        background-size: cover;
      }

      .logo {
        position: absolute;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;

        display: flex;
        justify-content: center;
        align-items: center;

        img {
          animation: logo-rotate 133.4s infinite;
          filter: drop-shadow(3px 3px 3px rgba(0,0,0,0.9));
          width: 75%;

          @media screen and ($desktop: $desktop-width) {
            max-width: 50vw;
          }
        }
      }

    }

    .about {
      display: flex;
      padding: 1% 3%;

      @media screen and ($desktop: $desktop-width) {
        >div {
          width: 50%;
          padding: 1%;
        }
      }

      @media screen and ($device: $device-width) {
        flex-direction: column;
      }

      iframe {
        width: 100%;
        height: 30vh;
      }
      h1 {
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
