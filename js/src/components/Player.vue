<template>

  <div class="player" v-if="summary">
    <div class="avatar">
      <img :alt="player.nick" :src="person.avatar"/>
    </div>

    <div class="name">
      <span :class="player.color">{{player.nick}}</span>
    </div>

    <div class="points">
      {{summary.score}} pts, {{summary.matches}}m
    </div>
  </div>

</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Player',
  mixins: [DrunkenFallMixin],

  props: {
    player: Object,
    index: 0
  },

  computed: {
    avatar () {
      return this.player.avatar
    },
    person () {
      return this.$store.getters.getPerson(this.player.person_id)
    },
    summary () {
      return this.$store.getters.getPlayerSummary(
        this.tournament.id,
        this.player.person_id,
      )
    },
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.player {
  display: flex;
  flex-grow: 1;
  flex-basis: 0;
  background-color: $bg-default-alt;

  .avatar {
    display: flex;
    justify-content: center;
    align-items: center;

    padding: 0em 1em;

    img {
      display: inline-block;
      height: 11vh;
      width:  11vh;
      object-fit: cover;
      border-radius: 100%;
    }
  }
  .name {
    display: flex;
    align-items: center;
    font-size: 6.5vh;
    flex-grow: 1;
  }
  .points {
    display: flex;
    align-items: center;
    font-size: 5vh;
    padding-right: 1em;
  }
}

.player:nth-child(even) {
  background-color: $bg-default;
  img {
  }
}

</style>
