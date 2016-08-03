<template>
  <div>
    <header>
      <div class="content">
      <div class="title">{{tournament.name}} / Join</div>
      </div>
      <div class="clear"></div>
    </header>
    <form id="join" @submit="submit">
      <input v-model="name" class="text" type="text" maxlength="20" placeholder="Your player tag"/>

      <h2>Select your preferred archer</h2>

      <div>
        <img @click="character" id="green" src="/static/img/green-unselected.png">
        <img @click="character" id="blue" src="/static/img/blue-unselected.png">
        <img @click="character" id="pink" src="/static/img/pink-unselected.png">
        <img @click="character" id="orange" src="/static/img/orange-unselected.png">
        <img @click="character" id="white" src="/static/img/white-unselected.png">
        <img @click="character" id="yellow" src="/static/img/yellow-unselected.png">
        <img @click="character" id="cyan" src="/static/img/cyan-unselected.png">
        <img @click="character" id="purple" src="/static/img/purple-unselected.png">
      </div>

      <input type="submit" class="submit" value="Go go go!"
        v-bind:class="{'ready': ready}"/>
    </form>
  </div>
</template>

<script>
export default {
  name: 'Join',

  data () {
    return {
      tournament: {},
      can_join: false,
      name: '',
      color: ''
    }
  },

  methods: {
    character (event) {
      this.clear()
      var img = event.srcElement
      img.src = '/static/img/' + img.id + '-selected.png'
      img.style = 'opacity: 1'
      this.$data.color = img.id
    },
    clear () {
      var elem = document.getElementById('join').getElementsByTagName('img')
      for (var i = 0; i < elem.length; i++) {
        var item = elem[i]
        item.src = '/static/img/' + item.id + '-unselected.png'
        item.style = ''
      }
    },
    submit (event) {
      event.preventDefault()
      if (this.$data.ready === false) {
        console.log('no')
        return
      }

      var payload = {
        name: this.name,
        color: this.color
      }

      // Reset the form so that it is clean if the page is reused
      this.clear()
      this.name = ''
      this.color = ''

      this.$http.post('/api/towerfall/' + this.$data.tournament.id + '/join/', payload).then((res) => {
        // Success callback
        console.log(res)
        var j = res.json()
        console.log(j)
        this.$route.router.go('/towerfall' + j.redirect)
      }, (res) => {
        console.log('fail')
        console.log(res)
      })
    }
  },

  computed: {
    ready: function () {
      return this.$data.color !== '' && this.$data.name !== ''
    }
  },

  route: {
    data ({ to }) {
      this.$http.get('/api/towerfall/tournament/' + to.params.tournament + '/').then(function (res) {
        console.log(res.data)
        this.$set('tournament', res.data.Tournament)
        this.$set('can_join', res.data.CanJoin)
      }, function (res) {
        console.log('error when getting tournament')
        console.log(res)
      })
    }
  }
}
</script>

<style lang='scss' scoped>
#join {
  div {
    width: 80%;
    margin: 10px auto;
  }

  input.text {
    width: 11em;
    margin: 0 auto;
    display: inherit;
    text-align: center;
  }

  input.submit {
    margin: 20px auto;
    width: 250px;
    border: none;
    font-size: 2em;
    transition: 1.0s;
    background-color: #333333;

    &.ready {
      background-color: #405060;
    }
  }
  h2 {
    text-align: center;
  }

  img {
    image-rendering: pixelated;
    width: 22%;
    margin: 1%;
    opacity: 0.3;
    transition: 1.2s;
  }
}
</style>
