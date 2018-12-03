<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' / Casters - DrunkenFall'"></headful>
    <tournament-controls />

    <h2>Casters ({{inCast.length}} / 2)</h2>
    <div class="players joined">
      <div v-for="person in inCast" class="player">
        <img @click="remove" :id="person.id" :alt="person.nick" :src="person.avatar"/>
        <p>{{person.firstName}}</p>
      </div>
    </div>

    <h2>Sheeple</h2>
    <div class="players not-joined">
      <div v-for="person in targets" class="player">
        <img @click="add" :id="person.id" :alt="person.nick" :src="person.avatar"/>
        <p>{{person.firstName}}</p>
      </div>
    </div>

    <div class="links standalone">
      <a @click="commit" :class="{ disabled: !canCommit}">
        <div class="icon positive">
          <icon name="cloud-upload"></icon>
        </div>
        <p>Set</p>
        <div class="clear"></div>
      </a>
    </div>
  </div>
</template>

<script>
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"
import TournamentControls from "./buttons/TournamentControls"

export default {
  name: 'Casters',
  mixins: [DrunkenFallMixin],
  components: {
    TournamentControls,
  },

  data () {
    return {
      selected: [],
    }
  },

  methods: {
    add (e) {
      if (this.tournament.casters.length < 2) {
        let p = _.find(this.targets, (p) => p.id === e.target.id)
        this.tournament.casters.push(p)
      }
    },
    remove (e) {
      this.selected = _.remove(this.tournament.casters, (p) => p.id === e.target.id)
    },
    commit () {
      let $vue = this
      let b = _.map(this.inCast, (p) => p.id)
      console.log(b)

      this.api.casters({id: this.tournament.id}, b).then((res) => {
        console.debug("casters response:", res)
        this.$router.push({name: "tournament", params: {id: this.tournament.id}})
      }, (err) => {
        $vue.$alert("Setting casters failed. See console.")
        console.error(err)
      })
    }
  },

  computed: {
    targets () {
      let $vue = this
      let r = _.filter(this.people, function (p) {
        let pl = _.find($vue.inCast, function (o) {
          return p.id === o.id
        })
        return pl === undefined
      })
      return _.sortBy(r, ['name'])
    },
    toCast () {
      let $vue = this
      return _.filter(this.targets, function (p) {
        let pl = _.find($vue.selected, function (o) {
          return p.id === o
        })
        return pl !== undefined
      })
    },
    inCast () {
      return this.tournament.casters
    },
    canCommit () {
      return this.selected.length === 2
    },
  },

  created () {
    this.api = this.$resource("/api", {}, {
      casters: { method: "POST", url: "/api/tournaments/{id}/casters/" },
    })
  }
}
</script>

<style lang="scss" scoped>

* {
  text-align: center;
}

.players {
  text-align: center;
  width: 80%;
  margin: 10px auto;
  min-height: 120px;

  .player {
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
