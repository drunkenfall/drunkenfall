<template>
  <div v-if="people">
    <h2>Combatants</h2>
    <div class="players">
      <router-link :to="{name: 'profile', params: {id: person.id}}" v-for="person in people" class="person">
        <img :id="person.id" :alt="person.nick" :src="person.avatar"/>
        <p>{{person.displayName}}</p>
      </router-link>
    </div>
  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'People',
  mixins: [DrunkenFallMixin],

  computed: {
    people () {
      return this.$store.state.people
    }
  },
}
</script>

<style lang="scss" scoped>

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

  a {
    box-shadow: none;
  }

  .person {
    display: inline-block;
    width: 100px;
    margin-top: 0px;
    cursor: pointer;

    img {
      object-fit: cover;
      border-radius: 100%;
      width:  100px;
      height: 100px;
      box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
      background-color: rgba(10,12,14,0.3);
      margin-bottom: -30px;
    }
    p {
      width: 80%;
      text-align: center;
      padding: 0.2em 0.3em;
      margin: 0.5em auto;
      display: inline-block;
      box-shadow: 1px 1px 6px rgba(0,0,0,0.5);
      font-weight: bold;
      background-color: #333339;
    }
  }
}
.joined p {
  background-color: #406040;
}
.not-joined p {
  background-color: #604040;
}

</style>
