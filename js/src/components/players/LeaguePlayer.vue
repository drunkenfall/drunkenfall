<template>
<router-link :to="{name: 'archer', params: {id: person.id}}" :class="cls">
  <img :id="person.id" :alt="person.nick" :src="person.avatar"/>
  <div class="data">
    <div class="position">
      <div class="ordinal">{{ordinal(index+1)}}</div>
      <div class="trophies" v-if="stats.total.wins > 0">
        <span v-for="n in stats.total.wins">🏆</span>
      </div>
    </div>
    <p class="name" :class="person.color">{{person.displayName}}</p>
  </div>
</router-link>
</template>

<script>
import Person from '../../models/Person.js'
import DrunkenFallMixin from "../../mixin"

export default {
  name: "LeaguePlayer",
  mixins: [DrunkenFallMixin],
  props: {
    person: new Person(),
    index: Number,
  },
  computed: {
    active () {
      if (this.$route.params.id === undefined) {
        return this.combatants[0].person
      }
      return this.$store.getters.getPerson(
        this.$route.params.id
      )
    },
    cls () {
      if (this.person.id === this.active.id) {
        return "person active"
      }

      return "person"
    },
    stats () {
      if (!this.$store.state.stats) {
        return undefined
      }
      return this.$store.state.stats[this.person.id]
    },
  },
}

</script>

<style lang="scss" scoped>
@import "../../css/colors.scss";

.person {
  display: flex;
  align-items: center;
  padding: 0.9em 0.9em;
  box-shadow: none;

  &:nth-child(odd) {
    background-color: $bg-default-alt;
  }

  &.active {
    background-color: $bg-default-dark;
  }

  &:hover {
    background-color: $bg-default-hover;
  }

  .position {
    display: flex;
    justify-content: space-between;
    font-size: 1.5em;
    color: $fg-disabled;
  }

  .data {
    margin-left: 1em;
    flex-grow: 1;
  }

  img {
    object-fit: cover;
    border-radius: 100%;
    width:  75px;
    height: 75px;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    background-color: rgba(10,12,14,0.3);
  }

  .name {
    vertical-align: middle;
    /* margin: 0.4em 0.3em 0; */
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);
    font-size: 2.2em;
  }
}

</style>
