import moment from 'moment'
import _ from 'lodash'

export const INTERVAL = 1000
const PLACEHOLDER = 'Soon™'

export class Countdown {
  constructor () {
    this.intervalID = 0
    this.time = PLACEHOLDER
  }

  pad (n, width) {
    n = n + ''
    return n.length >= width ? n : new Array(width - n.length + 1).join("0") + n
  }

  start (event) {
    // If there is already a clock ticking, kill it.
    if (this.intervalID !== 0) {
      clearInterval(this.intervalID)
    }

    this.intervalID = setInterval(() => {
      var eventTime = event.unix()
      var currentTime = moment().unix()
      var diffTime = eventTime - currentTime
      var d = moment.duration(diffTime, 'seconds') // duration

      d = moment.duration(d - INTERVAL, 'milliseconds')

      // If we're ever at a negative interval, stop immediately.
      // Technically we probably only really need the seconds here, but
      // if we use all of them any future cases will be fixed immediately.
      if (_.some([d.hours(), d.minutes(), d.seconds()], (n) => n < 0)) {
        console.log("Closing interval.")
        clearInterval(this.intervalID)
        this.time = PLACEHOLDER
        return
      }

      // Add hours left, but only if there are hours left.
      let hours = ""
      if (d.hours() > 0) {
        hours = this.pad(d.hours(), 2) + ":"
      }

      if (d.hours() === 0 && d.minutes() === 0 && d.seconds() === 0) {
        this.time = PLACEHOLDER
      } else {
        this.time = hours + this.pad(d.minutes(), 2) + ":" + this.pad(d.seconds(), 2)
      }
    }, INTERVAL)
  }
}

export class Clock {
  constructor () {
    this.time = this._getTime()
  }

  start () {
    setInterval(() => {
      this.time = this._getTime()
    }, 1000)
  }

  _getTime () {
    return moment().format("HH:mm (Z)")
  }
}
