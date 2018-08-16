<template>
  <div>
    <headful title="Settings - DrunkenFall"></headful>
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

          <input type="image" @click="character" :class="colorClass('green')" id="green" :src="archerImg('green')">
          <input type="image" @click="character" :class="colorClass('blue')" id="blue" :src="archerImg('blue')">
          <input type="image" @click="character" :class="colorClass('pink')" id="pink" :src="archerImg('pink')">
          <input type="image" @click="character" :class="colorClass('orange')" id="orange" :src="archerImg('orange')">
          <input type="image" @click="character" :class="colorClass('white')" id="white" :src="archerImg('white')">
          <input type="image" @click="character" :class="colorClass('yellow')" id="yellow" :src="archerImg('yellow')">
          <input type="image" @click="character" :class="colorClass('cyan')" id="cyan" :src="archerImg('cyan')">
          <input type="image" @click="character" :class="colorClass('purple')" id="purple" :src="archerImg('purple')">
          <input type="image" @click="character" :class="colorClass('red')" id="red" :src="archerImg('red')">
        </div>
      </div>

      <div class="section selector">
        <div class="split">
          <input id="normal" name="normal" v-model="archer_type" type="radio" value="0" />
          <label for="normal">Normal archers</label>
        </div>

        <div class="split">
          <input id="alternate" name="normal" v-model="archer_type" type="radio" value="1" />
          <label for="alternate">Alternate archers</label>
        </div>
      </div>

      <div class="links">
        <a @click="submit" :class="{disabled: !canSave}">
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
      archer_type: "0",
    }
  },

  methods: {
    colorClass (color) {
      if (color === this.color) {
        return "selected"
      }
      return ""
    },
    archerImg (color) {
      if (this.archer_type === "1") {
        return `/static/img/${color}-alt.png`
      }
      return `/static/img/${color}-selected.png`
    },
    character (event) {
      event.preventDefault()
      let color = event.target.id
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
        color: this.color,
        archer_type: this.archer_type === "1" ? 1 : 0,
      }

      this.$http.post('/api/user/settings/', payload).then((res) => {
        var j = res.json()
        $vue.$store.commit('setUser', Person.fromObject(j.person, $vue.$cookie))
        this.$router.push({name: "tournaments"})
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
      if (!q.nick && !u.nick && this.name) {
        this.$set(this.$data, "nick", this.name.split(" ")[0])
      } else {
        this.$set(this.$data, "nick", q.nick ? q.nick : u.nick)
      }

      this.$set(this.$data, "color", q.color ? q.color : u.color)
      this.$set(this.$data, "archer_type", u.archer_type.toString())
    }
  },

  computed: {
    canSave () {
      let u = this.user
      return u.color !== this.color || u.name !== this.name || u.nick !== this.nick || u.archer_type.toString() !== this.archer_type
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
@import "../css/colors.scss";

form {
  .section {
    width: 80%;
    margin: 10px auto;

    &.selector {
      height: 2em;
      display: flex;
      align-items: center;
      justify-content: space-around;
      margin-bottom: 2em;

      .split {
        margin: 0 2em;

        label {
          display: inline;
          margin: 0;
        }
      }
    }
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
