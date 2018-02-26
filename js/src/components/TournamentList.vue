<template>
<div class="flexmaster">
  <headful title="Tournaments - DrunkenFall"></headful>

  <div class="sidebar-buttons" v-if="user && user.isProducer && showSidebar">
    <div class="links">
      <button-link :to="{ name: 'new'}"
        :icon="'plus'" :iconClass="'positive'" :label="'New'" />

      <button-link :func="clear" :cls="{ disabled: !canClear}"
        :icon="'trash'" :iconClass="'danger'" :label="'Clear tests'" />
    </div>
  </div>

  <h1>DrunkenFall 2018</h1>

  <!--
  <p>Pellentesque dapibus suscipit ligula.  Donec posuere augue in
  quam.  Etiam vel tortor sodales tellus ultricies commodo.
  Suspendisse potenti.  Aenean in sem ac leo mollis blandit.  Donec
  neque quam, dignissim in, mollis nec, sagittis eu, wisi.  Phasellus
  lacus.  Etiam laoreet quam sed arcu.  Phasellus at dui in ligula
  mollis ultricies.  Integer placerat tristique nisl.  Praesent augue.
  Fusce commodo.  Vestibulum convallis, lorem a tempus semper, dui dui
  euismod elit, vitae placerat urna tortor vitae lacus.  Nullam libero
  mauris, consequat quis, varius et, dictum id, arcu.  Mauris mollis
  tincidunt felis.  Aliquam feugiat tellus ut neque.  Nulla facilisis,
  risus a rhoncus fermentum, tellus tellus lacinia purus, et dictum
  nunc justo sit amet elit. </p>
 -->

  <div class="tournaments" :class="{ loading: !tournaments }">
    <router-link
      v-for="t in currentLeague"
      :tournament="t.id"
      track-by="id"
      :to="{ name: 'tournament', params: { tournament: t.id }}"
      class="tournament">
      <img class="cover" v-if="t.cover" :src="t.cover" />
      <div class="cover" v-else />
      <div class="text">
        <div class="title">{{t.subtitle}}</div>
        <div class="date">{{t.scheduled.format("MMMM Do")}}</div>
      </div>
    </router-link>
  </div>

  <h1 v-if="tournaments.length === 0">
    Loading... 💜
  </h1>
</div>
</template>

<script>
import _ from "lodash"
import DrunkenFallMixin from "../mixin"
import ButtonLink from "./buttons/ButtonLink"

export default {
  name: 'TournamentList',
  mixins: [DrunkenFallMixin],
  components: {
    ButtonLink,
  },

  methods: {
    clear (event) {
      if (!this.canClear) {
        return
      }

      let $vue = this
      event.preventDefault()
      return this.$http.get('/api/tournament/clear/').then(function (res) {
        console.log(res)
      }, function (res) {
        $vue.$alert("Clearing failed. See console.")
        console.error(res)
        return { tournaments: [] }
      })
    },
  },
  computed: {
    canClear () {
      return _.some(this.tournaments, 'isTest')
    },
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/fonts.scss";

h1 {
  @media screen and ($desktop: $desktop-width) {
    @include display4();
  }
  @media screen and ($device: $device-width) {
    @include display2();
  }
  margin: 0.3em;
}

p {
  @include body2();
  width: 50%;
  margin: 1em auto;
}

.flexmaster {
  display: flex;
  flex-direction: column;
  min-height: 100%;
}

.tournaments {
  margin: 3% 0%;
  padding: 0% 5%;
  transition: 0.3s ease-in-out;
  flex-grow: 1;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  @media screen and ($device: $device-width) {
    height: 100%;
  }
  @media screen and ($desktop: $desktop-width) {
    align-content: space-around;
  }

  .tournament {
    background-color: $bg-default;
    display: block;
    text-align: center;
    margin: 1%;

    @media screen and ($desktop: $desktop-width) {
      width: 30%;
    }

    .cover {
      width: 100%;
      background-color: rgba(0,0,0,0.5);
    }

    .text {
      padding: 0.75em;
      .title {
        @include display2();
      }
      .date {
        @include title();
        color: $fg-secondary;
      }
    }
  }
}

</style>
