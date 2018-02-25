<template>
<div>
  <headful title="DrunkenFall"></headful>
  <div id="start">

    <div class="main">
      <div class="logo">
        <img alt="DrunkenFall" src="/static/img/oem-text.svg"/>
      </div>

      <div class="about">
        <h1>Wait, what's this?</h1>
        <p>
          DrunkenFall is a monthly video game tournament showdown of
          archery skills. We play
          <a href="http://towerfall-game.com" target="_blank">TowerFall</a>
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
          with just pen and paper, and now we have a streaming platform
          with an app that helps us keep track of who fights who!
        </p>

        <h1>Can I join?</h1>
        <p>
          Oh, I thought you'd never ask! Our motto has always been '<i>the more the
          merrier</i>', and that will always stay true!
        </p>
      </div>
    </div>
    <div class="schedule">
      <div class="next">
        <div class="title">Next event</div>
        <div class="tournament"></div>
        <div class="subtitle">Fifty Shades of Arrows</div>
        <div class="date">February 17th</div>
      </div>
      <!-- <div class="upcoming">
      <h2>Future events</h2>
      </div> -->
    </div>
  </div>
</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'About',
  mixins: [DrunkenFallMixin],
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

#start {
  display: flex;

  .main {
    width: 66.6%;

    >div {
      margin: 2em 5em;
    }

    .logo {
      display: flex;

      img {
        height: 17em;
      }
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
        font-size: 1.4em;
        a {
          box-shadow: none;
          font-weight: bold;
          color: $accent;
        }
      }
    }
  }

  .schedule {
    width: 33.3%;
    text-align: center;

    .next {
      .tournament {
        margin: 1em auto;
        background-color: rgba(10,0,0,0.3);
        width: 14em;
        height: 20em;
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
