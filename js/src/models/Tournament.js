import store from '../core/store'
import { isGoZeroDateOrFalsy } from '../util/date.js'
import moment from 'moment'
import Event from './Event.js'
import Player from './Player.js'
import Person from './Person.js'
import _ from 'lodash'

export default class Tournament {
  static fromObject (obj) {
    let t = new Tournament()
    t.raw = obj
    Object.assign(t, obj)

    t.opened = moment(t.opened)
    t.scheduled = moment(t.scheduled)
    t.started = moment(t.started)
    t.ended = moment(t.ended)
    t.qualifyingEnded = moment(t.qualifying_end)

    t.players = _.map(t.players, Player.fromObject)
    t.casters = _.map(t.casters, Person.fromObject)

    let events = t.events
    _.each(t.matches, (m) => {
      events = _.concat(events, m.events)
    })

    events = _.omitBy(events, _.isNil) // TODO(thiderman): Why are there nil items? hwat
    events = _.sortBy(events, [(o) => { return o.time }])
    events = _.reverse(events)
    events = _.map(events, Event.fromObject)
    t.events = events

    return t
  }

  get next () {
    return this.currentMatch.index
  }

  playerJoined (person) {
    let p = store.getters.getPlayerSummary(this.id, person.id)
    return p !== undefined
  }

  get numeral () {
    let n = this.name.split(" ")[1]
    return n.substring(0, n.length - 1) // Remove the colon
  }

  get numeralColor () {
    if (this.color) {
      return this.color
    }

    return "default-numeral"
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

  get matches () {
    return store.getters.getMatches(this.id)
  }

  get betweenMatches () {
    if (this.isEnded) {
      return false
    } else if (this.isStarted && !this.matches[0].isStarted) {
      return true
    }

    // If there are no matches that are currently running, we are inbetween matches
    return _.filter(this.matches, 'isRunning').length === 0
  }

  get endedRecently () {
    return moment().isBetween(this.ended, moment(this.ended).add(6, 'hours'))
  }

  get canStart () {
    return !this.isStarted
  }

  get isUpcoming () {
    return moment().isBefore(this.scheduled) && this.canStart
  }

  get isNext () {
    if (store.getters.upcoming.length === 0) {
      return false
    }
    return store.getters.upcoming[0].id === this.id
  }

  get isToday () {
    return moment().isSame(this.scheduled, 'day')
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }

  get isUsurpable () {
    return true
    // return this.players.length < 32
  }

  get currentMatch () {
    let c = _.first(_.filter(this.matches, (m) => !m.isEnded))
    return c
  }

  // This is supposed to be used when you need the next match before
  // it is started. It returns the upcoming match if the current one
  // is ended.
  get nextMatch () {
    let n = _.first(_.filter(this.matches, 'canStart'))
    return n
  }

  get nextNextMatch () {
    let n = _.filter(this.matches, 'canStart')[1]
    return n
  }

  get runnerups () {
    return store.getters.runnerups(this.id)
  }

  get playoffs () {
    // throw new Error("call to non-ported tournament.playoffs()")
    return _.slice(this.matches, 0, this.matches.length - 5)
  }

  get semis () {
    throw new Error("call to non-ported tournament.semis()")
    // let l = this.matches.length
    // return _.slice(this.matches, l - 3, l - 1)
  }

  get final () {
    return this.matches[this.matches.length - 1]
  }

  get qualifyingOpen () {
    return isGoZeroDateOrFalsy(this.qualifyingEnded)
  }
}
