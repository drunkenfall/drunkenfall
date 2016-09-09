<template>
  <div>
    <header>
      <div class="content">
        <div class="title">Drunken TowerFall</div>
      </div>
      <div class="links">
        <a v-link="{path: 'new/'}">New Tournament</a>
        <a v-link="{path: '/facebook/'}">Facebook</a>
        </div>
      <div class="clear"></div>
    </header>

    <div class="tournaments" :class="{ loading: !tournament }">
      <div v-for="tournament in tournaments" :tournament="tournament.id" track-by="id">
        <a v-link="{ name: 'tournament', params: { tournament: tournament.id }}">{{tournament.name}}</a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TournamentList',

  // components: {
  //   Series
  // },

  data () {
    return {
      tournaments: []
    }
  },

  route: {
    data ({ to }) {
      return this.$http.get('/api/towerfall/tournament/').then(function (res) {
        return {
          tournaments: res.data
        }
      }, function (res) {
        console.warn('error when getting tournaments')
        console.error(res)
        return { tournaments: [] }
      })
    }
  }
}
</script>

<style lang="scss">
.tournaments a {
  background-color: #405060;
  color: #dbdbdb;
  display: block;
  font-size: 300%;
  font-weight: bold;
  padding: 1% 3%;
  text-align: center;
  text-decoration: none;
  width: 600px;
  margin: 10px auto;
}
</style>
