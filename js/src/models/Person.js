export const PRODUCER = 100
export const COMMENTATOR = 50
export const JUDGE = 30
export const PLAYER = 10

export default class Person {
  static fromObject (obj) {
    let p = new Person()
    Object.assign(p, obj)
    p.authenticated = p.session && true || false
    return p
  }

  get avatar () {
    if (this.avatar_url) {
      return this.avatar_url
    }
    return "https://graph.facebook.com/" + this.facebook_id + "/picture?width=9999"
  }

  get color () {
    if (!this.color_preference) {
      return "white"
    }
    return this.color_preference[0]
  }

  get displayName () {
    return this.nick
  }

  get firstName () {
    if (!this.name) {
      return ""
    }
    return this.name.split(" ")[0]
  }

  get isProducer () {
    return this.userlevel >= PRODUCER
  }
  get isCommentator () {
    return this.userlevel >= COMMENTATOR
  }
  get isJudge () {
    return this.userlevel >= JUDGE
  }
  get isPlayer () {
    return this.userlevel >= PLAYER
  }
}
