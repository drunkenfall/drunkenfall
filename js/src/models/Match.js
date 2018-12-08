import { isGoZeroDateOrFalsy } from '../util/date.js'
import moment from 'moment'
import _ from 'lodash'
import Player from './Player.js'
// import Tournament from './Tournament.js'

import store from '../core/store.js'

export default class Match {
  static fromObject (obj, t) {
    let m = new Match()
    Object.assign(m, obj)

    m.started = moment(m.started)
    m.ended = moment(m.ended)
    m.scheduled = moment(m.scheduled)

    m.endScore = m.length
    m.commits = []

    // TODO(thiderman): There used to be this weird bug where some
    // matches are created without a length. This circumvents this in
    // a semi-ugly way. The bug is fixed, but workarounds are forever. <3
    if (m.endScore === 0) {
      m.endScore = m.kind === "final" ? 20 : 10
    }

    if (t !== undefined) {
      m.tournament = t
    }

    return m
  }

  get matchID () {
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
    if (this.kind === "special") {
      return "YaLo Winnah"
    }
    return _.capitalize(this.kind) // + " " + this.relativeIndex
  }

  get players () {
    if (this._players && this._players.length !== 0) {
      return this._players
    }

    let x = store.getters.matchPlayers(this.tournament_id, this.index)
    return x
  }

  set players (v) {
    this._players = _.map(_.sortBy(v, "id"), Player.fromObject)
  }

  get relativeIndex () {
    if (this.kind === "playoff") {
      // return this.index - (this.tournament.matches.length -)
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

  get canReset () {
    return this.isRunning && this.commits.length > 0
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }

  get levelTitle () {
    if (this.level === "twilight") {
      return "Twilight Spire"
    } else if (this.level === "kingscourt") {
      return "King's Court"
    } else if (this.level === "frostfang") {
      return "Frostfang Keep"
    } else if (this.level === "sunken") {
      return "Sunken City"
    } else if (this.level === "amaranth") {
      return "The Amaranth"
    }
    return this.level.charAt(0).toUpperCase() + this.level.slice(1)
  }
};
