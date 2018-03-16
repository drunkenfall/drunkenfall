<template>
<div v-if="tournament">
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
import Person from '../models/Person.js'

export default {
  name: 'TournamentSummary',
  mixins: [DrunkenFallMixin],

  computed: {
    tournament () {
      return this.trackingTournament
    },
    winner () {
      return Person.fromObject(this.tournament.winners[0].person)
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
