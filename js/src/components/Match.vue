<template>
  <a v-link="{path: match.kind + '/' + match.index + '/'}">
    <div v-for="player in match.players" id="player-{{player.index}}" class="player {{player.classes}}">
      <p class="name">{{player.name}}</p>
    </div>
  </a>
</template>

<script>
export default {
  name: 'Match',
  props: {
    match: Object
  }
}
</script>

<style lang="scss" scoped>
.match {
  float: left;
  margin: 10px;
  cursor: pointer;
  width: 290px;
  text-align: center;

  img {
    width: 100%;
  }

}
</style>
