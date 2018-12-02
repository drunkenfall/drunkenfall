<template>

<div>
  <h1>End qualifying</h1>
  <div class="time">
    <input name="time" type="time" v-model="time"/>
  </div>
  <button @click="submit" id="doit">Go go</button>
</div>

</template>

<script>
import DrunkenFallMixin from "../mixin"
import moment from 'moment'

export default {
  name: 'EndQualifying',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      time: "",
    }
  },

  methods: {
    submit (e) {
      let $vue = this
      e.preventDefault()

      console.log(this.time)
      let spl = this.time.split(":")
      let target = moment().hour(spl[0]).minute(spl[1]).second(0)

      this.api.submit({ id: this.tournament.id }, {time: target.format()}).then((res) => {
        console.log("endq response:", res)
        this.$router.push({name: "tournament", params: {id: this.tournament.id}})
      }, (err) => {
        $vue.$alert("Endq failed. See console.")
        console.error(err)
      })
    },
  },

  created () {
    this.api = this.$resource("/api", {}, {
      submit: { method: "POST", url: "/api/tournaments/{id}/endqualifying/" },
    })
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

* {
  text-align: center;
}

h1 {
  font-size: 4em;

}

.time {
  flex-grow: 4;
  display: flex;
  align-items: center;
  justify-content: center;
  input {
    width: 80%;
    font-size: 3em;
    padding: 1em;
    background-color: $bg-default;
    color: $fg-default;
    border: none;
  }
}


#doit {
  flex-grow: 2;
  margin: 1em auto;
  width: 80%;
  border: none;
  font-size: 2em;
  transition: 1.0s;
  color: $fg-default;
  background-color: $secondary;
  padding: 0.4em 0.8em;
}

</style>
