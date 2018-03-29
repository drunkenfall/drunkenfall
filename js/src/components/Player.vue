<template>
  <div :class="'player ' + classes">
    <img :src="avatar" :alt="display_name"/>
    <p>{{display_name}}</p>
  </div>
</template>

<script>
export default {
  name: 'Player',

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
    classes () {
      if (this.match.isEnded) {
        // TODO(thiderman): This makes the old tournaments work
        // again. They lack color, but hey.
        if (!this.match.kill_order) {
          return this.player.color
        }

        if (this.index === this.match.kill_order[0]) {
          return 'gold'
        } else if (this.index === this.match.kill_order[1]) {
          return 'silver'
        } else if (this.index === this.match.kill_order[2] && this.match.kind === 'final') {
          return 'bronze'
        }

        return 'out'
      }

      return this.player.color
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.player {
  &.gold, &.silver, &.bronze {
    p {
      color: #fff !important;
    }
  }

  img {
    display: inline-block;
    height: 38px;
    width: 38px;
    object-fit: cover;
    border-radius: 100%;
    float: left;
    margin-left: 5px;
  }
}

.player:nth-child(even) {
  img {
    float: right;
    margin-right: 5px;
  }
}

</style>
