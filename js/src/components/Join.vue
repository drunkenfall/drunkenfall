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

      <div id="join-images" class="images">
        <input type="image" @click="character" id="green" src="/static/img/green-selected.png">
        <input type="image" @click="character" id="blue" src="/static/img/blue-selected.png">
        <input type="image" @click="character" id="pink" src="/static/img/pink-selected.png">
        <input type="image" @click="character" id="orange" src="/static/img/orange-selected.png">
        <input type="image" @click="character" id="white" src="/static/img/white-selected.png">
        <input type="image" @click="character" id="yellow" src="/static/img/yellow-selected.png">
        <input type="image" @click="character" id="cyan" src="/static/img/cyan-selected.png">
        <input type="image" @click="character" id="purple" src="/static/img/purple-selected.png">
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
      event.preventDefault()
      this.clear()
      var img = event.srcElement
      this.$data.color = img.id
      img.className = 'selected'
    },
    clear () {
      var elem = document.getElementById('join-images').getElementsByTagName('input')
      for (var i = 0; i < elem.length; i++) {
        var item = elem[i]
        item.className = ''
      }
    },
    submit (event) {
      event.preventDefault()
      if (this.ready === false) {
        console.log('Did not fill in all details')
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

  .images {
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
    }

    input {
      image-rendering: pixelated;
      width: 22%;
      margin: 1%;
      opacity: 0.2;
      transition: 0.3s;
      cursor: pointer;
    }
  }
}
</style>
