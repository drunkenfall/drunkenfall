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

  <div class="tournaments" :class="{ loading: !tournaments }">
    <div
      v-for="(t, i) in currentLeague"
      :tournament="t.id"
      track-by="id"
      @click="gotoTournament(t)"
      :class="{clickable: t.isEnded || t.isNext || t.isRunning}">
      <img class="cover" :src="t.cover" />
      <div class="text">
        <div v-if="t.isEnded || t.isRunning || t.isNext" class="title">
          {{t.subtitle}}
        </div>
        <div v-else-if="t.isUpcoming && i < 8" class="title dark">
          DrunkenFall 2018: {{i+1}}
        </div>
        <div v-else class="title">
          2018 Grand Finale
        </div>
        <div class="date">{{t.scheduled.format("MMMM Do")}}</div>
      </div>
    </div>
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
    gotoTournament (t) {
      if (t.isEnded || t.isNext || t.isRunning) {
        return this.$router.push({name: "tournament", params: {tournament: t.id}})
      }
    },
    clear (event) {
      if (!this.canClear) {
        return
      }

      let $vue = this
      event.preventDefault()
      return this.$http.delete('/api/tournament/').then(function (res) {
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

  >div {
    background-color: $bg-default;
    display: block;
    text-align: center;
    margin: 1%;

    @media screen and ($desktop: $desktop-width) {
      width: 30%;
    }

    &.clickable {
      cursor: pointer;
    }

    .cover {
      width: 100%;
      min-height: 232px;
      background-color: rgba(0,0,0,0.5);
    }

    .text {
      padding: 0.75em;
      .title {
        @include display2();
        &.dark {
          color: $fg-disabled;
        }
      }
      .date {
        @include title();
        color: $fg-secondary;
      }
    }
  }
}

</style>
