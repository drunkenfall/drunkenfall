<template>
  <div v-if="stats">
    <h2>Combatants</h2>
    <div class="players">
      <template v-for="c in combatants" ref="combatants">
        <list-player :person="c.person"></list-player>
      </template>
    </div>

    <h2>Unfought</h2>
    <div class="players unfought">
      <template v-for="c in unfought" ref="unfought">
        <list-player :person="c.person"></list-player>
      </template>
    </div>

  </div>
</template>

<script>
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"
import ListPlayer from './players/ListPlayer.vue'

export default {
  name: 'Archers',
  mixins: [DrunkenFallMixin],
  components: {
    ListPlayer,
  },
  computed: {
    stats () {
      return _.filter(this.$store.state.stats, (p) => {
        return !p.person.disabled
      })
    },
    people () {
      return this.$store.state.people
    },
    combatants () {
      return _.sortBy(_.filter(this.stats, (p) => {
        return p.total.score > 0
      }), 'rank')
    },
    unfought () {
      return _.sortBy(_.filter(this.stats, (p) => {
        return p.total.score === 0
      }), 'person.displayName')
    },
  },
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";

.sidebared-content {
  text-align: center;
}

h2 {
  font-size: 3em;
  margin-top: 0.4em;
  text-align: center;
}

.players {
  text-align: center;
  width: 80%;
  margin: 10px auto;
}

</style>
