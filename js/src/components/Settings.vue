<template>
  <div>
    <form @submit="submit">
      <div class="section">
        <label for="nick">Display name</label>
        <input class="text" name="nick" v-model="nick" type="text" value=""/>
      </div>

      <div class="section">
        <label for="name">Full name</label>
        <input class="text smaller" name="name" v-model="name" type="text" value=""/>
      </div>

      <div class="section">
        <div id="join-images" class="images">
          <label for="images">Preferred archer</label>

          <input type="image" @click="character" :class="colorClass('green')" id="green" src="/static/img/green-selected.png">
          <input type="image" @click="character" :class="colorClass('blue')" id="blue" src="/static/img/blue-selected.png">
          <input type="image" @click="character" :class="colorClass('pink')" id="pink" src="/static/img/pink-selected.png">
          <input type="image" @click="character" :class="colorClass('orange')" id="orange" src="/static/img/orange-selected.png">
          <input type="image" @click="character" :class="colorClass('white')" id="white" src="/static/img/white-selected.png">
          <input type="image" @click="character" :class="colorClass('yellow')" id="yellow" src="/static/img/yellow-selected.png">
          <input type="image" @click="character" :class="colorClass('cyan')" id="cyan" src="/static/img/cyan-selected.png">
          <input type="image" @click="character" :class="colorClass('purple')" id="purple" src="/static/img/purple-selected.png">
          <input type="image" @click="character" :class="colorClass('red')" id="red" src="/static/img/red-selected.png">
        </div>
      </div>

      <div class="links standalone">
        <a @click="submit" :class="{ disabled: !canSave}">
          <div class="icon positive">
            <icon name="floppy-o"></icon>
          </div>
          <p>Save</p>
          <div class="clear"></div>
        </a>
      </div>
    </form>

  </div>
  </template>

<script>
import Person from "../models/Person.js"
import DrunkenFallMixin from "../mixin.js"

export default {
  name: 'Settings',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      name: "",
      nick: "",
      color: "",
    }
  },

  methods: {
    colorClass (color) {
      if (color === this.color) {
        return "selected"
      }
      return ""
    },
    character (event) {
      event.preventDefault()
      let color = event.srcElement.id
      this.$data.color = color
      let name = document.getElementById("sidebar-username")
      if (name) {
        name.className = color // For funz <3
      }
    },
    submit (event) {
      let $vue = this
      event.preventDefault()
      if (this.canSave === false) {
        console.log('No changes to save.')
        return
      }

      var payload = {
        name: this.name,
        nick: this.nick,
        color: this.color
      }

      this.$http.post('/api/user/settings/', payload).then((res) => {
        var j = res.json()
        $vue.$store.commit('setUser', Person.fromObject(j.person, $vue.$cookie))
        $vue.$router.push('/towerfall' + j.redirect)
      }, (res) => {
        $vue.$alert('Post failed. See console.')
        console.error(res)
      })
    },
    updateData () {
      let q = this.$route.query
      let u = this.user

      // Set either the query parameter, or fall back to the session object.
      this.$set(this.$data, "name", q.name ? q.name : u.name)

      // If no nick is set, just suggest it from the full name
      if (!q.nick && !u.nick) {
        this.$set(this.$data, "nick", this.name.split(" ")[0])
      } else {
        this.$set(this.$data, "nick", q.nick ? q.nick : u.nick)
      }

      this.$set(this.$data, "color", q.color ? q.color : u.color)
    }
  },

  computed: {
    canSave () {
      let u = this.user
      return u.color !== this.color || u.name !== this.name || u.nick !== this.nick
    },
  },

  mounted () {
    this.updateData()
  },

  watch: {
    // This is so that a hard page reload will still be able to catch the changes.
    user (val) {
      this.updateData()
    },
    nick (val) {
      let name = document.getElementById("sidebar-username")
      if (val && name) {
        name.innerText = val
      }
    }
  },
}
</script>

<style lang="scss" scoped>
@import "../variables.scss";

form {
  .section {
    width: 80%;
    margin: 10px auto;
  }

  label {
    display: block;
    text-align: center;
    font-size: 1.5em;
    margin: 1.5em 0 0.5em 0;
  }

  input.text {
    width: 100%;
    margin: 0 auto;
    display: inherit;
    text-align: center;
    box-shadow: inset 5px 5px 4px rgba(0,0,0,0.1);

    &:focus {
      outline: 0;
    }

    &.smaller {
      font-size: 1em;
      width: 16em;
    }
  }

  h2 {
    text-align: center;
  }

  .images {
    margin-top: 3em;
    transition: 0.3s;

    > input {
      box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
      border: 1px solid #111117;
    }

    #green {background-image: url(/static/img/green-selected.png);}
    #blue {background-image: url(/static/img/blue-selected.png);}
    #pink {background-image: url(/static/img/pink-selected.png);}
    #orange {background-image: url(/static/img/orange-selected.png);}
    #white {background-image: url(/static/img/white-selected.png);}
    #yellow {background-image: url(/static/img/yellow-selected.png);}
    #cyan {background-image: url(/static/img/cyan-selected.png);}
    #purple {background-image: url(/static/img/purple-selected.png);}

    .selected {
      opacity: 1;
      box-shadow: 3px 3px 5px rgba(0,0,0,0.5);
    }

    text-align: center;

    input {
      image-rendering: pixelated;
      display: inline-block;
      opacity: 0.3;
      transition: 0.3s;
      cursor: pointer;
      outline: none;

      @media screen and ($desktop: $desktop-width) {
        width: 9.5%;
        margin: 0.3%;
      }
      @media screen and ($device: $device-width) {
        width: 14%;
        margin: 2%;
      }
    }
  }
}
</style>
