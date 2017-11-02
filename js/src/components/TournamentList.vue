<template>
  <div>
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
import Tournament from "../models/Tournament.js"
import DrunkenFallMixin from "../mixin"

export default {
  name: 'TournamentList',
  mixins: [DrunkenFallMixin],

  methods: {
    clear (event) {
      let $vue = this
      event.preventDefault()
      return this.$http.get('/api/tournament/clear/').then(function (res) {
        this.$set('tournaments', _.map(res.data.tournaments, Tournament.fromObject))
      }, function (res) {
        $vue.$alert("Clearing failed. See console.")
        console.error(res)
        return { tournaments: [] }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";

.tournaments {
  transition: 0.3s ease-in-out;
  margin: 5em 0;
}

.tournaments a {
  text-shadow: $shadow-default;
  background-color: $bg-disabled;
  color: $fg-default;
  display: block;
  font-size: 2.5em;
  font-weight: bold;
  padding: 1% 3%;
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
    background-color: #353535;
    font-size: 1.2em;
    width: 30%;
  }

  &.current {
    background-color: #405060;
    width: 60%;
    font-size: 3.5em;
  }
}
</style>
