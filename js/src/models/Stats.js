import moment from 'moment'
import _ from 'lodash'

class PlayerSnapshot {
  static fromObject (obj, cookie) {
    let ps = new PlayerSnapshot()
    Object.assign(ps, obj)
    ps.playtime = moment.duration(ps.playtime / 1000 / 1000)
    return ps
  }
}

export default class Stats {
  static fromObject (obj, cookie) {
    let s = new Stats()
    Object.assign(s, obj)

    _.each(s, (ps) => {
      ps.total = PlayerSnapshot.fromObject(ps.total)

      _.each(_.keys(ps.tournaments), (key) => {
        ps.tournaments[key] = PlayerSnapshot.fromObject(ps.tournaments[key])
      })
      ps.participated = _.keys(ps.tournaments).length
    })

    return s
  }
}
