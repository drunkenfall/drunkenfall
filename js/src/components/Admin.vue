<template>
<div v-if="userLoaded && user.isProducer">
  <h1>Superpowers</h1>

  <div class="section">
    <h2>Tournaments</h2>
    <div class="links">
      <button-link
        :to="{ name: 'new'}"
        :icon="'plus'"
        :iconClass="'positive'"
        :label="'New tournament'" />

      <button-link
        :to="{ name: 'new'}"
        :cls="{ disabled: !canClear}"
        :icon="'trash'"
        :iconClass="'danger'"
        :label="'Clear test tournaments'" />

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
      return this.$http.get('/api/tournament/clear/').then(function (res) {
        console.log(res)
      }, function (res) {
        $vue.$alert("Clearing failed. See console.")
        console.error(res)
        return { tournaments: [] }
      })
    },
    reset () {
      return this.runningMatch.reset()
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
      if (m.isEnded) {
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
  }
}
</script>

<style lang="scss">
@import "../variables.scss";

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
