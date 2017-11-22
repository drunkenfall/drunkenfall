<template>
  <router-link :to="{name: 'archer', params: {id: person.id}}" class="person">
    <img :id="person.id" :alt="person.nick" :src="person.avatar"/>
    <p :class="person.color">{{person.displayName}}</p>
  </router-link>
</template>

<script>
import Person from '../../models/Person.js'

export default {
  name: "ListPlayer",
  props: {
    person: new Person(),
  },
}

</script>

<style lang="scss" scoped>
@import "../../css/colors.scss";

a {
  box-shadow: none;
}

.person {
  display: inline-block;
  width: 135px;
  margin-top: 0px;
  cursor: pointer;

  img {
    object-fit: cover;
    border-radius: 100%;
    width:  125px;
    height: 125px;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    background-color: rgba(10,12,14,0.3);
    margin-bottom: -30px;
  }
  p {
    width: 80%;
    text-align: center;
    padding: 0.3em 0.3em 0.2em;
    margin: 0.5em auto;
    display: inline-block;
    box-shadow: 1px 1px 6px rgba(0,0,0,0.5);
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);

    font-size: 1.2em;
    background-color: $bg-default;
    z-index: 100;
  }
}

</style>
