<template>
<div v-if="tournament">
  <tournament-controls />

  <div class="subheader" v-if="user.isCommentator && nextMatch && !tournament.isEnded">
    <div v-if="!nextMatch.isScheduled">
      <p>
        Pause until
          <span>{{tournament.currentMatch.title}}</span>
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

    <div class="category playoffs">
      <h2>Playoffs</h2>
      <div class="matches">
        <template v-for="m in tournament.playoffs">
          <match-overview :match="m" :class="'match ' + m.kind"></match-overview>
        </template>
      </div>
      <div class="clear"></div>
    </div>

    <div class="category semis">
      <h2>Semi-finals</h2>
      <div class="matches">
        <template v-for="m in tournament.semis">
          <match-overview :match="m" :class="'match ' + m.kind"></match-overview>
        </template>
      </div>
      <div class="clear"></div>
    </div>
    <div class="category final">
      <h2>Final</h2>
      <div class="matches">
        <match-overview :match="tournament.final" class="match final"></match-overview>
      </div>
    </div>
  </div>
</template>

<script>
import MatchOverview from './MatchOverview'
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"
import _ from 'lodash'

export default {
  name: 'TournamentOverview',
  mixins: [DrunkenFallMixin],

  components: {
    MatchOverview,
    TournamentControls,
  },

  computed: {
    canParticipants () {
      return this.tournament.currentMatch.kind === 'playoff'
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

    usurp () {
      this.tournament.usurp()
    },
  },

  created () {
    this.api = this.$resource("/api", {}, {
      setTime: { method: "GET", url: "/api/{/id}/time{/time}" },
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.tournament {
  position: relative;
}

@media screen and ($desktop: $desktop-width) {
  .playoffs, .semis, .final {
    width: 29%;
    float: left;
    margin-left: 3%;
    position: relative;
  }
}

@media screen and ($device: $device-width) {
  .playoffs, .semis, .final {
    /* width: 90%; */
    margin: 0 auto;
  }
}

.category h3 {
  text-align: center;
  font-size: 200%;
  margin: 4%;
}

.match {
  width: 100%;
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
    background-color: $bg-default;
  }
  .runnerup:nth-child(even) {
    background-color: #272727;
  }
}

.selected-runnerups {
  .runnerup:nth-child(odd) {
    background-color: $bg-default;
  }
  .runnerup:nth-child(even) {
    background-color: #394939;
  }

  .button {
    width: 50px;
    margin: 10px auto;
    padding: 0.3em 0.5em;

    cursor: pointer;
    text-shadow: 1px 1px 1px rgba(0,0,0,0.4);
    text-align: center;
  }
}

.subheader {
  @include subheading();
  width: 80%;

  @media screen and ($desktop: $desktop-width) {
    p {
      float: left;
    }
    .links {
      float: right;
      a {
        float: right;
      }
    }
  }
  @media screen and ($device: $device-width) {
    & {
      text-align: center;
      padding: 0.5em;
    }
    .links {
      a:last-child {
        margin-bottom: 1em;
      }
    }
  }

  margin: 30px auto;
  background-color: $bg-default;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
  text-shadow: 2px 2px 3px rgba(0,0,0,0.5);

  p {
    font-size: 2em;
    padding: 0.3em 0.5em;

    span {
      text-transform: capitalize;
    }
  }

  .links {
    a, .action {
      @include button();
      margin: 15px  !important;
      background-color: $button-bg;
      border-left: 3px solid $accent;
      color: $fg-default;
      display: block;
      font-weight: bold;
      padding: 7px 30px;
      text-align: center;
      text-decoration: none;
      margin: 10px auto;
      min-width: 60px;

      padding: 0.5em 0.7em;
    }
  }
}

</style>
