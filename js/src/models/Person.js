export const PRODUCER = 100
export const COMMENTATOR = 50
export const JUDGE = 30
export const PLAYER = 10

export default class Person {
  static fromObject (obj, cookie) {
    let p = new Person()
    Object.assign(p, obj)
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
      return ""
    }
    return this.color_preference[0]
  }

  get displayName () {
    return this.nick || this.firstName
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

  get authenticated () {
    if (this.userlevel === undefined) {
      return false
    }
    return this.userlevel > 0
  }

  logout ($vue) {
    console.log("Logging out...")
    let $this = this

    $vue.$http.get('/api/towerfall/user/logout/').then(function (res) {
      console.log("Logged out: ", res)
      $this.userlevel = 0
      $vue.$store.commit("setUser", $this)
    }, function (res) {
      $vue.$alert("Logging out failed. See console.")
      console.error(res)
    })
  }
}
