import { isGoZeroDateOrFalsy } from '../util/date.js'
import moment from 'moment'
import _ from 'lodash'
import Player from './Player.js'

export default class Match {
  static fromObject (obj, $vue, t) {
    let m = new Match()
    Object.assign(m, obj)
    m.$vue = $vue

    m.started = moment(m.started)
    m.ended = moment(m.ended)
    m.scheduled = moment(m.scheduled)
    m.players = _.map(m.players, Player.fromObject)

    m.endScore = m.length

    // TODO(thiderman): There is this weird bug where some matches are
    // created without a length. This circumvents this in a semi-ugly
    // way.
    if (m.endScore === 0) {
      m.endScore = m.kind === "final" ? 20 : 10
    }

    let root = "/api/towerfall/tournament{/id}{/index}"
    m.api = $vue.$resource("/api/towerfall", {}, {
      start: { method: "GET", url: `${root}/start/` },
      commit: { method: "POST", url: `${root}/commit/` },
      end: { method: "GET", url: `${root}/end/` },
      reset: { method: "GET", url: `${root}/reset/` },
    })

    if (t !== undefined) {
      m.tournament_id = t.id
    }

    return m
  }

  start () {
    let $vue = this
    if (this.tournament.shouldBackfill) {
      console.log("Not starting match; backfill required.")
      return
    }

    console.log("Starting match...")
    this.api.start(this.id).then((res) => {
      console.log("Match started.", res)
    }, (res) => {
      $vue.$alert("Starting failed. See console.")
      console.error(res)
    })
  }

  end () {
    let $vue = this
    console.log("Ending match...")
    this.api.end(this.id).then((res) => {
      console.log("Match ended.", res)
      this.$vue.$router.push(`/towerfall/${this.tournament_id}/`)
    }, (res) => {
      $vue.$alert("Ending failed. See console.")
      console.error(res)
    })
  }

  reset () {
    let $vue = this
    console.log("Resetting match...")
    this.api.reset(this.id).then((res) => {
      console.log("Match reset.", res)
    }, (res) => {
      $vue.$alert("Reset failed. See console.")
      console.error(res)
    })
  }

  // TODO(thiderman): This could somehow not be migrated to here from
  // Match.vue. When moved, the request turns from a POST into a GET
  // and the backend rightfully denies it. A thing for later, I guess.
  // commit ($control, payload) {
  //   this.api.commit(this.id, payload).then((res) => {
  //     console.log("Round committed.")
  //     _.each($control.players, (p) => { p.reset() })
  //   }, (res) => {
  //     console.error('error when setting score', res)
  //   })
  // }

  get id () {
    return {
      id: this.tournament_id,
      index: this.index,
    }
  }

  get endScore () { return this._end }
  set endScore (value) { this._end = value }

  get title () {
    if (this.kind === "final") {
      return "Final"
    }
    return _.capitalize(this.kind) + " " + this.relativeIndex
  }

  get relativeIndex () {
    if (this.kind === "semi") {
      return this.index - this.tournament.playoffs.length + 1
    }
    return this.index + 1
  }

  get isStarted () {
    // match is started if 'started' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.started)
  }

  get isEnded () {
    // match is ended if 'ended' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.ended)
  }

  get isScheduled () {
    return !isGoZeroDateOrFalsy(this.scheduled)
  }

  get canStart () {
    return !this.isStarted
  }

  get canEnd () {
    // can't end if already ended
    if (this.isEnded) {
      return false
    }

    // can end if at least one player has enough kills (ie >= end)
    return _.some(this.players, (player) => { return player.kills >= this.endScore })
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }

  get chartData () {
    var out = []
    for (var i = 0; i < this.players.length; i++) {
      out.push([0])
      _.forEach(this.commits, function (commit) {
        let pastScore = _.last(out[i])
        let roundScore = _.sum(commit.kills[i])
        out[i].push(pastScore + roundScore)
      })
    }
    return out
  }

  get tournament () {
    return this.$vue.$store.getters.getTournament(this.tournament_id)
  }
};
