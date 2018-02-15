<template>
  <div>
    <headful title="Tournaments - DrunkenFall"></headful>
    <div class="sidebar-buttons" v-if="user && user.isProducer && showSidebar">
      <div class="links">
        <button-link :to="{ name: 'new'}"
          :icon="'plus'" :iconClass="'positive'" :label="'New'" />

        <button-link :func="clear" :cls="{ disabled: !canClear}"
          :icon="'trash'" :iconClass="'danger'" :label="'Clear tests'" />
      </div>
    </div>

    <div class="tournaments" :class="{ loading: !tournaments }">
      <div v-for="tournament in tournaments"
        :tournament="tournament.id" track-by="id">

        <router-link :to="{ name: 'tournament', params: { tournament: tournament.id }}"
          :class="{ test: tournament.isTest, current: !tournament.isStest && !tournament.isStarted}">
          {{tournament.name}}
        </router-link>
      </div>
    </div>

    <h1 v-if="tournaments.length === 0">
      Loading... &lt;3
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

.tournaments {
  transition: 0.3s ease-in-out;
  margin: 5em 0;
}

.tournaments a {
  @include display3();

  text-shadow: $shadow-default;
  background-color: $bg-default;
  color: $fg-default;
  display: block;
  padding: 1.3% 3.4% 1.1%;
  text-align: center;
  text-decoration: none;
  width: 50%;
  margin: 0.2em auto;
  border-left: 5px solid $accent;
  transition: 0.3s;

  &:hover {
    background-color: $bg-default-hover;
  }

  &.test {
    @include display1();
    background-color: $bg-disabled-secondary;
    width: 30%;
  }

  &.current {
    background-color: $secondary;
    width: 60%;
    font-size: 3.5em;
  }
}
</style>
