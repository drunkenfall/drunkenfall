import _ from 'lodash'

let DEATH = ['ded', 'rip', 'pwnt', 'rekt']

export default class PlayerState {
  static fromObject (obj) {
    let ps = new PlayerState()
    Object.assign(ps, obj)
    return ps
  }

  get death () {
    return _.sample(DEATH)
  }
}
