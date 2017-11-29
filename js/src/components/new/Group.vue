<template>
  <div>
    <form v-on:submit="create">
      <h1 v-html="name" :class="color"></h1>

      <div id="colors">
        <div class="green" @click="setColor('green')"></div>
	      <div class="blue" @click="setColor('blue')"></div>
	      <div class="pink" @click="setColor('pink')"></div>
	      <div class="orange" @click="setColor('orange')"></div>
	      <div class="white" @click="setColor('white')"></div>
	      <div class="yellow" @click="setColor('yellow')"></div>
	      <div class="cyan" @click="setColor('cyan')"></div>
	      <div class="purple" @click="setColor('purple')"></div>
	      <div class="red" @click="setColor('red')"></div>
      </div>

      <label for="name" v-model="name">Name</label>
      <input class="text name" v-model="name" name="name" type="text" value="" placeholder="Name"/>

      <label for="id">ID</label>
      <input class="text id" v-model="id" name="id" type="text" value="" placeholder="ID"/>

      <label for="scheduled">Scheduled to start at</label>
      <input class="text date" v-model="scheduled"/>

      <div class="links standalone">
        <a @click="create" :class="{ disabled: !isReady}">
          <div class="icon positive">
            <icon name="plus"></icon>
          </div>
          <p>Create</p>
          <div class="clear"></div>
        </a>
      </div>
    </form>
  </div>
</template>

<script>
import moment from "moment"
import _ from 'lodash'

export default {
  name: 'DrunkenFallNew',

  data () {
    return {
      name: '',
      id: '',
      color: '',
      scheduled: null,
    }
  },

  created () {
    let $vue = this
    this.scheduled = moment().hour(20).minute(0).second(0).format()

    this.$set(
      this.$data,
      "color",
      _.sample(["green", "blue", "pink", "orange", "white", "yellow", "cyan", "purple", "red"])
    )

    // Come on baby, fake it for me <3
    return this.$http.get('/api/fake/name/').then(function (res) {
      console.log(res)
      let j = JSON.parse(res.data)
      this.$set(this.$data, "name", j.name)
      this.$set(this.$data, "id", j.numeral)
    }, function (res) {
      $vue.$alert("Faking failed. See console.")
      console.error(res)
    })
  },

  methods: {
    setColor (color) {
      this.$set(this.$data, "color", color)
    },
    create (event) {
      if (!this.isReady) {
        console.log("Not ready to post.")
        return
      }
      let $vue = this
      event.preventDefault()

      var payload = {
        name: this.name,
        id: this.id,
        scheduled: moment(this.scheduled),
        color: this.color,
      }

      this.$http.post('/api/new/', payload).then((res) => {
      }, (res) => {
        $vue.$alert("Creating failed. See console.")
        console.error(res)
      })
    },
  },

  computed: {
    isReady () {
      return this.name !== '' && this.id !== ''
    }
  },
}
</script>

<style lang="scss">
@import "../../css/colors.scss";
h1 {
  font-size: 3em;
  transition: 0.3s;
}

#colors {
  width: 70%;
  height: 3em;

  div {
    width: 9%;
    margin: 1%;
    float: left;
    height: 100%;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
    cursor: pointer;

    &.green  {background-color: $green-bg;}
    &.blue   {background-color: $blue-bg;}
    &.pink   {background-color: $pink-bg;}
    &.orange {background-color: $orange-bg;}
    &.white  {background-color: $white-bg;}
    &.yellow {background-color: $yellow-bg;}
    &.cyan   {background-color: $cyan-bg;}
    &.purple {background-color: $purple-bg;}
    &.red    {background-color: $red-bg;}
  }
}

form {
  display: flex;
  flex-direction: column;
  align-items: center;

  label {
    font-size: 1.5em;
    margin: 1.5em 0 0.5em 0;
  }

  div.fake-checkbox {
    margin: 1.5em 0 0.5em 0;
  }

  input.text {
    width: 11em;

    &.name {
      width: 28em;
    }
    &.date {
      width: 22em;
      font-size: 1em;
    }
  }
}
</style>
