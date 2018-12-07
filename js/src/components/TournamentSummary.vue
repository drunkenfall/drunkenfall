<template>
<div v-if="tournament && winner">
  <headful :title="tournament.subtitle + ' - Summary'"></headful>

  <div class="players">
    <div class="winner">
      <h1>A winner is you!</h1>
      <img :alt="winner.name" :src="winner.avatar"/>
      <h1>✨ Congration you done it! ✨</h1>
    </div>
  </div>
</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
// import Person from '../models/Person.js'
import _ from 'lodash'

export default {
  name: 'TournamentSummary',
  mixins: [DrunkenFallMixin],

  computed: {
    tournament () {
      return this.tournaments[this.$route.params.tournament]
    },
    winner () {
      return _.sortBy(this.tournament.final.players, 'kills')[0]
    }
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";
@import "../css/fonts.scss";
@import "../css/ribbon.scss";

.players {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;

  .winner {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;

    h1 {
      @include display4();
    }

    img {
      height: 650px;
      width:  650px;
      object-fit: cover;
      border-radius: 100%;
    }

  }
}

</style>
