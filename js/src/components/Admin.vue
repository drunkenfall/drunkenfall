<template>
<div v-if="userLoaded && user.isProducer">
  <h1>Superpowers</h1>

  <div class="section">
    <h2>Tournaments</h2>
    <div class="links">
      <router-link class="action"
        :to="{ name: 'new'}">
        <div class="icon positive">
          <icon name="plus"></icon>
        </div>
        <p>New tournament</p>
        <div class="clear"></div>
      </router-link>

      <a @click="clear"
        :class="{ disabled: !canClear}">
        <div class="icon danger">
          <icon name="trash"></icon>
        </div>
        <p>Clear test tournaments</p>
        <div class="clear"></div>
      </a>

    </div>
  </div>

  <div class="section">
    <h2>Users</h2>
    <div class="links">
      <router-link class="action" v-if="user.isJudge" :to="{ name: 'disable'}">
        <div class="icon">
          <icon name="ban"></icon>
        </div>
        <p>Disable users</p>
        <div class="clear"></div>
      </router-link>
    </div>
  </div>

</div>
</template>

<script>
import _ from "lodash"
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Admin',
  mixins: [DrunkenFallMixin],
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
  }
}
</script>

<style lang="scss" scoped>

.section {
  width: 45%;
  float: left;
  margin: 2.5%;

  h2 {
    font-size: 2em;
    text-align: center;
  }

  .links {
    padding: 2.5%;

    background-color: #333339;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

    a {
      font-size: 1.5em;
      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}

</style>
