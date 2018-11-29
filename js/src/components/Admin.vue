<template>
<div v-if="userLoaded && user.isProducer && tournaments">
  <headful title="Superpowers - DrunkenFall"></headful>
  <div class="section">
    <h2>Tournaments</h2>
    <div class="links">
      <button-link
        :to="{ name: 'new'}"
        :icon="'plus'"
        :iconClass="'positive'"
        :label="'New tournament'" />

      <button-link
        :func="clear"
        :cls="{ disabled: !canClear}"
        :icon="'trash'"
        :iconClass="'danger'"
        :label="'Clear tests'" />

    </div>

    <div class="tournament" v-for="t in tournaments">
      <router-link :to="{ name: 'tournament', params: {'tournament': t.id}}">
        {{t.name}}
      </router-link>
    </div>
  </div>

  <div class="section">
    <h2>Users</h2>
    <div class="links">
      <button-link
        :to="{ name: 'disable'}"
        :icon="'ban'"
        :label="'Disable users'" />
    </div>
  </div>

  <div class="section">
    <h2>Matches</h2>
    <div class="links">
      <button-link
        :func="reset"
        :cls="{ disabled: !canReset}"
        :iconClass="'danger'"
        :icon="'recycle'"
        :label="matchLabel" />

    </div>
  </div>

</div>
</template>

<script>
import _ from "lodash"
import DrunkenFallMixin from "../mixin"
import ButtonLink from "./buttons/ButtonLink"

export default {
  name: 'Admin',
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
      return this.$http.delete('/api/tournaments/').then(function (res) {
        console.log(res)
      }, function (res) {
        $vue.$alert("Clearing failed. See console.")
        console.error(res)
        return { tournaments: [] }
      })
    },

    reset () {
      let $vue = this
      console.log("Resetting match...", this.runningMatch)
      this.api.reset(this.runningMatch.matchID, {}).then((res) => {
        console.log("Match reset.", res)
      }, (res) => {
        $vue.$alert("Reset failed. See console.")
        console.error(res)
      })
    },
  },
  computed: {
    canClear () {
      return _.some(this.tournaments, 'isTest')
    },
    canReset () {
      if (!this.runningMatch) {
        return false
      }
      return this.runningMatch.canReset
    },
    runningMatch () {
      if (!this.runningTournament) {
        return
      }
      let m = this.runningTournament.currentMatch
      if (m && m.isEnded) {
        return
      }
      return m
    },
    matchLabel () {
      if (!this.runningMatch) {
        return "Reset match"
      }

      return `Reset ${this.runningMatch.title}`
    },
  },
  created () {
    this.api = this.$resource("/api", {}, {
      reset: { method: "POST", url: "/api/tournaments/{id}/match/{index}/reset/" },
    })
  },
}
</script>

<style lang="scss">
@import "../css/colors.scss";

.tournament {
  margin: 0.5em;
  padding: 0.5em 1em;
  background-color: $bg-default;
  a {
  }
}

.section {
  @media screen and ($desktop: $desktop-width) {
    width: 28%;
    float: left;
  }

  margin: 2.5%;

  h2 {
    font-size: 2em;
    text-align: center;
  }

  .links {
    padding: 2.5%;

    background-color: $bg-default;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

    > div {
      &:last-child a {
        margin-bottom: 0;
      }
    }
  }
}

</style>
