<template>
<div v-if="tournament">
  <headful :title="tournament.subtitle + ' - DrunkenFall'"></headful>
  <tournament-controls />


</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"
import _ from 'lodash'

export default {
  name: 'TournamentOverview',
  mixins: [DrunkenFallMixin],

  components: {
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



</style>
