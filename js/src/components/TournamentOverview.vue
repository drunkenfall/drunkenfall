<template>
  <div>
    <div class="subheader" v-if="user.isCommentator && nextMatch && !tournament.isEnded">
      <div v-if="!nextMatch.isScheduled">
        <p>
          Pause until
          <span>{{tournament.currentMatch.title}}</span>:
        </p>
        <div class="links">
          <a @click="setTime(10)">10 min</a>
          <a @click="setTime(7)">7 min</a>
          <a @click="setTime(5)">5 min</a>
          <a @click="setTime(3)">3 min</a>
        </div>
        <div class="clear"></div>
      </div>
      <div v-if="nextMatch.isScheduled">
        <p class="center">
          <span>{{tournament.currentMatch.title}}</span> scheduled at
          {{nextMatch.scheduled.format("HH:mm")}}
        </p>
        <div class="clear"></div>
      </div>
    </div>

    <div class="category tryouts">
      <h3>Tryouts</h3>
      <div class="matches">
        <template v-for="m in tournament.tryouts">
          <match-overview :match="m" :class="'match ' + m.kind"></match-overview>
        </template>
      </div>
      <div class="clear"></div>
    </div>

    <div class="category semis">
      <h3>Semi-finals</h3>
      <div class="matches">
        <template v-for="m in tournament.semis">
          <match-overview :match="m" :class="'match ' + m.kind"></match-overview>
        </template>
      </div>
      <div class="clear"></div>
    </div>
    <div class="category final">
      <h3>Final</h3>
      <div class="matches">
        <match-overview :match="tournament.final" class="match final"></match-overview>
      </div>
    </div>
  </div>
</template>

<script>
import MatchOverview from './MatchOverview'
import DrunkenFallMixin from "../mixin"
import _ from 'lodash'

export default {
  name: 'TournamentOverview',
  mixins: [DrunkenFallMixin],

  components: {
    MatchOverview,
  },

  computed: {
    canParticipants () {
      return this.tournament.currentMatch.kind === 'tryout'
    },
    runnerups () {
      let t = this.tournament

      if (!t.runnerups) {
        return []
      }

      return _.map(t.runnerups, (runnerup) => {
        return _.find(t.players, function (p) {
          return p.person.id === runnerup.id
        })
      })
    },
  },

  methods: {
    setTime (x) {
      let $vue = this
      this.api.setTime({ id: this.tournament.id, time: x }).then((res) => {
        console.debug("settime response:", res)
        console.log('Redirect to /towerfall' + res.json().redirect)
        // this.$route.router.push('/towerfall' + j.redirect)
      }, (err) => {
        $vue.$alert("Setting time failed. See console.")
        console.error(err)
      })
    },
    selectRunnerup (p) {
      if (this.isSelected(p)) {
        // TODO(thiderman): Doesn't work. Fuck this.
        console.log("selected, to remove", this.selected)
        this.selected = _.remove(this.selected, function (o) {
          console.log(o.person.id, p.person.id)
          return o.person.id === p.person.id
        })
        return
      }

      this.selected.push(p)
    },
    isSelected (p) {
      return _.find(this.selected, p) !== undefined
    },
  },

  created () {
    this.api = this.$resource("/api/towerfall", {}, {
      setTime: { method: "GET", url: "/api/towerfall{/id}/time{/time}" },
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";

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
  position: relative;

}

.runnerups, .selected-runnerups {
  width: 100%;
  margin: 10px;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

  .runnerup {
    padding: 0.1em 0.3em;
    font-size: 24px;
    color: #aaa;

    p {
      margin: 1px;
      &.name {
        float: left;
        // font-weight: bold;
      }
      &.score {
        float: right;
      }
      b {
        text-shadow: 1px 1px 1px rgba(0,0,0,0.4);
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

.selected-runnerups {
  .runnerup:nth-child(odd) {
    background-color: #405040;
  }
  .runnerup:nth-child(even) {
    background-color: #394939;
  }

  .button {
    width: 50px;
    margin: 10px auto;
    padding: 0.3em 0.5em;
    background-color: #405060;
    cursor: pointer;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.4);
    text-align: center;
  }
}

.subheader {
  width: 80%;
  margin: 30px auto;
  background-color: #333339;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
  text-shadow: 2px 2px 3px rgba(0,0,0,0.5);

  p {
    font-size: 2em;
    padding: 0.3em 0.5em;
    float: left;

    span {
      text-transform: capitalize;
    }
  }

  .links {
    float: right;
    a, .action {
      margin: 15px  !important;
      font-size: 24px;
      float: right;
      background-color: #405060;
      color: #dbdbdb;
      display: block;
      font-weight: bold;
      padding: 7px 30px;
      text-align: center;
      text-decoration: none;
      margin: 10px auto;
      min-width: 60px;
    }
  }
}

</style>
