<template>
  <div v-if="stats" class="archers">
    <headful :title="active.nick + ' - DrunkenFall'"></headful>
    <div class="players">
      <template v-for="(c, i) in combatants" ref="combatants">
        <league-player :person="c.person" :index="i" :ref="c.person.id"></league-player>
      </template>
    </div>

    <profile :profile="active" class="selected"></profile>
  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import Profile from './Profile.vue'
import LeaguePlayer from './players/LeaguePlayer.vue'

export default {
  name: 'Archers',
  mixins: [DrunkenFallMixin],
  components: {
    LeaguePlayer,
    Profile,
  },
  computed: {
    active () {
      if (this.combatants.length !== 0 && this.$route.params.id === undefined) {
        return this.combatants[0].person
      }

      return this.$store.getters.getPerson(
        this.$route.params.id
      )
    },
    people () {
      return this.$store.state.people
    },
  },
  beforeRouteUpdate (to, from, next) {
    this.$refs[to.params.id][0].$el.scrollIntoView({
      behavior: "smooth",
      block: "center",
      inline: "nearest",
    })
    next()
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.sidebared-content {
  text-align: center;
}

h2 {
  font-size: 3em;
  margin-top: 0.4em;
  text-align: center;
}

.archers {
  display: flex;
  align-content: flex-end;
  flex-direction: row-reverse;
}

.selected {
  width: 71%;
}

.players {
  float: right;
  width: 25%;
  background-color: $bg-default;
  box-shadow: -4px 0px 4px rgba(0, 0, 0, 0.3);
}

</style>
