<template>
<div v-if="profile && stats">
  <div class="profile">
    <div class="intro">
      <div class="picture">
        <img class="avatar" :src="profile.avatar" :alt="profile.nick" />
        <p>
          {{ordinal(stats.rank)}}
          <span v-for="n in stats.total.wins">🏆</span>
        </p>
      </div>


      <div class="header">
        <h1 :class="profile.color">
          {{profile.displayName}}
        </h1>
        <h3>
          {{stats.total.kills}} <span>kills</span> /
          {{stats.total.rounds}} <span>rounds</span>
        </h3>
      </div>
    </div>

    <div class="achievements">
      <div class="achievement">
        <div class="title">🥃 Drunkard</div>
        <div class="commentary">Had more than 10 shots</div>
      </div>

      <div class="achievement">
        <div class="title">🐕 Underdog</div>
        <div class="commentary">Lost playoffs, still won tourmament</div>
      </div>

      <div class="achievement">
        <div class="title">👴 Veteran</div>
        <div class="commentary">Made it to the semis at least three times</div>
      </div>
    </div>

    <div class="stats" :class="profile.color">
      <div class="stat">
        <div class="data">
          <span class="main">{{stats.participated}}</span>
        </div>
        <p>Tournaments</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.matches}}</span>
        </div>
        <p>Matches</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.rounds}}</span>
        </div>
        <p>Rounds</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.kills}}</span>
        </div>
        <p>Kills</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.self}}</span>
        </div>
        <p>Suicides</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.sweeps}}</span>
        </div>
        <p>Sweeps</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.shots}}</span>
        </div>
        <p>Shots</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.score}}</span>
        </div>
        <p>Score</p>
      </div>

      <div class="stat">
        <div class="data">
          <span class="main">{{stats.total.playtime.minutes()}}</span>
          <span>m</span>
        </div>
        <p>Game time</p>
      </div>

      <div class="clear"></div>
    </div>

  </div>
</div>
</template>

<script>
import Person from '../models/Person.js'
import DrunkenFallMixin from "../mixin"

export default {
  name: "Profile",
  mixins: [DrunkenFallMixin],
  props: {
    profile: new Person(),
    index: Number,
  },
  computed: {
    stats () {
      if (!this.$store.state.stats) {
        return undefined
      }
      return this.$store.state.stats[this.profile.id]
    },
  },
  methods: {
    getOrdinal (n) {
      var s = ["th", "st", "nd", "rd"]
      var v = n % 100
      return s[(v - 20) % 10] || s[v] || s[0]
    }
  },
}

</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/fonts.scss";

.profile {
  width: 60%;
  margin: 3em auto;
  text-align: center;
  position: fixed;
  top: 0;
  display: flex;
  flex-direction: column;

  .intro {
    display: flex;
    align-content: start;

    .picture {
      width: 50%;
      display: flex;
      flex-direction: column;

      p {
        font-size: 5em;
        margin-top: 0.25em;
      }
    }

    .header {
      width: 50%;
      display: flex;
      flex-direction: column;

      h1 {
        display: block;
        @include display4();
        text-align: left;
      }
      h3 {
        display: block;
        font-size: 6em;
        margin-top: 0;
        text-align: left;

        span {
          font-size: 0.5em;
          color: $fg-disabled;
        }
      }
    }
  }

  .achievements {
    display: none !important;
    margin-top: 3em;
    display: flex;
    justify-content: space-around;

    .achievement {
      box-shadow: $shadow-default;
      width: 25%;
      background-color: $bg-default;
      padding: 1.2em;

      .title {
        @include display2();
      }
      .commentary {
        margin-top: 0.5em;
        color: $fg-secondary;
        @include title();
      }
    }
  }

  .stats {
    display: flex;
    background-color: $bg-default;

    /* &.green  {background-color: $green-bg-darkest;} */
    /* &.blue   {background-color: $blue-bg-darkest;} */
    /* &.pink   {background-color: $pink-bg-darkest;} */
    /* &.orange {background-color: $orange-bg-darkest;} */
    /* &.white  {background-color: $white-bg-darkest;} */
    /* &.yellow {background-color: $yellow-bg-darkest;} */
    /* &.cyan   {background-color: $cyan-bg-darkest;} */
    /* &.purple {background-color: $purple-bg-darkest;} */
    /* &.red    {background-color: $red-bg-darkest;} */

    width: 80%;
    margin: 3rem auto;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

    .stat {
      width: 20%;
      padding-bottom: 1em;

      &:nth-child(odd) {
        background-color: rgba(0,0,0,0.1);
      }

      .data {
        margin-top: 0.5rem;
        height: 5em;
        span.main {
          font-size: 5em;
        }
      }
      p {
        color: $fg-secondary;
      }
    }
  }

  .graphs {
    background-color: rgba(0,0,0,0.3);
    padding: 2em 1em;
    font-size: 5em;
    color: rgba(255,255,255,0.1);
  }

  img.avatar {
    object-fit: cover;
    border-radius: 100%;
    margin: 0 auto;
    width:  300px;
    height: 300px;
    min-height: 200px;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    background-color: rgba(10,12,14,0.3);
  }


}
</style>
