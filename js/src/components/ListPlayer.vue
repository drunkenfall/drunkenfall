<template>

  <div class="player" v-if="player.person">
    <div class="avatar">
      <img :alt="player.nick" :src="player.person.avatar"/>
    </div>

    <div class="name">
      <span :class="color">{{player.person.nick}}</span>
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
      return this.player.person.nick
    },
    color () {
      return this.player.color || this.player.person.preferred_color
    }
  },
  created () {
    this.player.person = this.$store.getters.getPerson(this.player.person_id)
    console.log(this)
  }
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
