import { isGoZeroDateOrFalsy } from '../util/date.js'
import moment from 'moment'
import Player from './Player.js'
import Person from './Person.js'
import Match from './Match.js'
import _ from 'lodash'

export default class Tournament {
  static fromObject (obj, $vue) {
    let t = new Tournament()
    Object.assign(t, obj)

    t.opened = moment(t.opened)
    t.scheduled = moment(t.scheduled)
    t.started = moment(t.started)
    t.ended = moment(t.ended)

    t.tryouts = _.map(t.tryouts, (m) => { return Match.fromObject(m, $vue) })
    t.semis = _.map(t.semis, (m) => { return Match.fromObject(m, $vue) })
    t.final = Match.fromObject(t.final, $vue)

    t.players = _.map(t.players, Player.fromObject)
    t.runnerups = _.map(t.runnerups, Person.fromObject)

    let root = "/api/towerfall/{/id}"

    t.api = $vue.$resource("/api/towerfall/", {}, {
      startTournament: { method: "GET", url: `${root}/start/` },
      next: { method: "GET", url: `${root}/next/` },
      reshuffle: { method: "GET", url: `${root}/reshuffle/` },
      usurp: { method: "GET", url: `${root}/usurp/` },
    })

    return t
  }

  next ($vue) {
    console.log("Going to next match...")
    this.api.next({ id: this.id }).then((res) => {
      console.debug("next response:", res)
      $vue.$router.push('/towerfall' + res.json().redirect)
    }, (err) => {
      console.error(`next for ${this.tournament} failed`, err)
    })
  }

  usurp () {
    this.api.usurp({ id: this.id }).then((res) => {
      console.log("usurp response", res)
    }, (err) => {
      console.error(`usurp for ${this.tournament} failed`, err)
    })
  }

  get numeral () {
    let n = this.name.split(" ")[1]
    return n.substring(0, n.length - 1) // Remove the colon
  }

  get subtitle () {
    return this.name.replace(/.*: /, "")
  }

  get isStarted () {
    return !isGoZeroDateOrFalsy(this.started)
  }

  get isEnded () {
    return !isGoZeroDateOrFalsy(this.ended)
  }

  get isTest () {
    return !this.name.startsWith('DrunkenFall')
  }

  get canStart () {
    return !this.isStarted
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }

  get canShuffle () {
    // We can only shuffle after the tournament has started (otherwise
    // technically no matches exists, so nothing can be shuffled
    // into), and before the first match has been started.
    let match = Match.fromObject(this.tryouts[0])
    return this.isStarted && !match.isStarted
  }

  get isUsurpable () {
    return this.players.length < 32
  }
}
