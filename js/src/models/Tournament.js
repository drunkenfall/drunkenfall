import { isGoZeroDateOrFalsy } from '../util/date.js'
import moment from 'moment'

export default class Tournament {
  static fromObject (obj) {
    let t = new Tournament()
    Object.assign(t, obj)

    t.opened = moment(t.opened)
    t.scheduled = moment(t.scheduled)
    t.started = moment(t.started)
    t.ended = moment(t.ended)

    return t
  }

  get isStarted () {
    // tournament is started if 'started' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.started)
  }

  get isEnded () {
    // tournament is ended if 'ended' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.ended)
  }

  get canStart () {
    return !this.isStarted
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }
}
