<template>
  <div>
    <form v-on:submit="create">
      <h1>Create tournament! |o/</h1>
      <label for="name">Name</label>
      <input class="text name" v-model="name" name="name" type="text" value="" placeholder="Name"/>

      <label for="id">ID</label>
      <input class="text id" v-model="id" name="id" type="text" value="" placeholder="ID"/>

      <label for="scheduled">Scheduled to start at</label>
      <input class="text date" v-model="scheduled"/>

      <div class="fake-checkbox">
        <label for="checkbox">Fake</label>
        <input type="checkbox" id="checkbox" v-model="fake">
      </div>

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

export default {
  name: 'New',

  data () {
    return {
      name: '',
      id: '',
      scheduled: null,
      fake: false,
    }
  },

  created () {
    let $vue = this
    this.scheduled = moment().hour(20).minute(0).second(0).format()

    // Come on baby, fake it for me <3
    return this.$http.get('/api/towerfall/fake/name/').then(function (res) {
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
        fake: this.fake,
      }

      this.$http.post('/api/towerfall/new/', payload).then((res) => {
        var j = res.json()
        this.$router.push('/towerfall' + j.redirect)
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
h1 {
  font-size: 3em;
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
