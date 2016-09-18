export default class Player {
  static fromObject (obj) {
    let p = new Player()
    Object.assign(p, obj)
    return p
  }

  // TODO(thiderman): Fix so that every avatar is set per tournament
  get avatar () {
    return "https://graph.facebook.com/" + this.person.facebook_id + "/picture?width=9999"
  }

  get color () {
    return this.person.color_preference[0] || ""
  }
}
