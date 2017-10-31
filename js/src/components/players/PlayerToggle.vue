<template>
  <div>
    <h2>{{onLabel}} ({{on.length}})</h2>
    <div class="people on">
      <div v-for="person in on" class="person">
        <img @click="toggle" :id="person.id" :alt="person.nick" :src="person.avatar"/>
        <p>{{person.name.split(" ")[0]}}</p>
      </div>
    </div>

    <h2>{{offLabel}} ({{off.length}})</h2>
    <div class="people off">
      <div v-for="person in off" class="person">
        <img @click="toggle" :id="person.id" :alt="person.nick" :src="person.avatar"/>
        <p>{{person.name.split(" ")[0]}}</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PlayerToggle',

  props: [
    "on",
    "onLabel",
    "off",
    "offLabel",
    "toggle",
  ]
}
</script>

<style lang="scss" scoped>

.sidebared-content {
  text-align: center;
}

h2 {
  margin: 2rem;
  font-size: 2em;
  text-align: center;
}

.people {
  text-align: center;
  width: 80%;
  margin: 10px auto;

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
    }
  }
}
.on p {
  background-color: #406040;
}
.off p {
  background-color: #604040;
}

</style>
