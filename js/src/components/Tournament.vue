<template>
  <div>
    <header>
      <div class="content">
        <div class="title">{{tournament.name}}</div>
      </div>
      <div v-if="can_join" class="links">
        <a v-link="{path: 'join/'}">Join</a>
      </div>
      <div class="clear"></div>
    </header>

    <div class="category tryouts">
      <h3>Tryouts</h3>
      <div class="matches">
        <template v-for="m in tournament.tryouts">
          <match :match="m" id="{{m.kind}}-{{m.index}}" class="match {{m.kind}}">
        </template>
      </div>
      <div class="clear"></div>
    </div>

    <div class="category semis">
      <h3>Semi-finals</h3>
      <div class="matches">
        <template v-for="m in tournament.semis">
          <match :match="m" id="{{m.kind}}-{{m.index}}" class="match {{m.kind}}">
        </template>
      </div>
      <!--
      <h3>Runnerups</h3>
      <div class="runnerups">
        <div v-for="m in tournament.GetRunnerups">
          <div class="runnerup">
            <p class="name">{{.Name}}</p>
            <p class="score">{{.Score}}p / {{.Matches}}m</p>
            <div class="clear"></div>
          </div>
        </div>
      </div>
      -->

    </div>
    <div class="category final">
      <h3>Final</h3>
      <div class="matches">
      <!--
        <match :match="tournament.final"
               id="{{tournament.final.kind}}-{{tournament.final.index}}"
               class="match {{tournament.final.kind}}">
       -->
      </div>
    </div>
  </div>
</template>

<script>
import Match from './Match.vue'

export default {
  name: 'Tournament',

  components: {
    Match
  },

  data () {
    return {
      tournament: {},
      can_join: false
    }
  },

  route: {
    data ({ to }) {
      this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
        this.$set('tournament', res.data.Tournament)
        this.$set('can_join', res.data.CanJoin)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    }
  }
}
</script>

<style lang="scss">
.tournament {
  position: relative;
}

.tryouts, .semis, .final {
  width: 29%;
  float: left;
  margin-left: 3%;
  position: relative;
}

.category h3 {
  text-align: center;
  font-size: 200%;
  margin: 4%;
}

.match {
  width: 100%;
  font-size: 150%;
  display: block;

  .player {
    float: left;
    width: 50%;
    line-height: 120%;
    text-align: center;
    overflow: hidden;
    white-space: nowrap;
    padding: 0.2em 0;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.7);

    &.prefill:nth-child(1), &.prefill:nth-child(4) {
      background-color: #333;
    }
    &.prefill:nth-child(2), &.prefill:nth-child(3) {
      background-color: #383838;
    }
    &.prefill {
      color: #555;
    }

    &.green  { background-color: #4E9110; }
    &.blue   { background-color: #4C7CBA; }
    &.pink   { background-color: #E39BB5; }
    &.orange { background-color: #CF9648; }
    &.white  { background-color: #dbdbdb; }
    &.yellow { background-color: #D1BD66; }
    &.cyan   { background-color: #59C2C1; }
    &.purple { background-color: #762c7a; }
  }
}

.matches .match .player {
  height: 40%;
}

.runnerups {
  .runnerup {
    padding: 0.1em 0.3em;
    font-size: 24px;
    color: #aaa;

    p {
      margin: 1px;
      .name {
        float: left;
        font-weight: bold;
      }
      .score {
        float: right;
      }
    }

  }
  .runnerup:nth-child(odd) {
    background-color: #333;
  }
  .runnerup:nth-child(even) {
    background-color: #272727;
  }

}

</style>
