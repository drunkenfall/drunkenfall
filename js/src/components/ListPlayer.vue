<template>

  <div class="player" v-if="person">
    <div class="avatar">
      <img :alt="player.nick" :src="person.avatar"/>
    </div>

    <div class="name">
      <span :class="color">{{person.nick}}</span>
    </div>
  </div>

</template>

<script>
export default {
  name: 'ListPlayer',

  props: {
    player: Object,
    match: {},
    index: 0
  },

  computed: {
    avatar () {
      return this.player.avatar
    },
    display_name () {
      return this.person.nick
    },
    color () {
      return this.player.color || this.person.preferred_color
    },
    person () {
      return this.$store.getters.getPerson(this.player.person_id)
    }
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
  padding: 0.5em 0em;

  .avatar {
    display: flex;
    justify-content: center;
    align-items: center;

    padding: 0em 1em;

    img {
      display: inline-block;
      height: 3vh;
      width:  3vh;
      object-fit: cover;
      border-radius: 100%;
    }
  }
  .name {
    display: flex;
    align-items: center;
    font-size: 4vh;
    flex-grow: 1;
  }
}

.player:nth-child(even) {
  background-color: $bg-default;
  img {
  }
}

</style>
