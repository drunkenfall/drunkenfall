<template>
  <div>
    <header>
      <div class="content">
        <div class="title">Drunken TowerFall / New tournament</div>
      </div>
      <div class="clear"></div>
    </header>

    <div class="clear"></div>

    <form v-on:submit="create">
      <input v-model="name" name="name" type="text" value="" placeholder="Name"/>
      <input v-model="id" name="id" type="text" value="" placeholder="ID"/>

      <input type="submit"/>
    </form>
  </div>
</template>

<script>
export default {
  name: 'New',

  data () {
    return {
      name: '',
      id: ''
    }
  },

  methods: {
    create (event) {
      event.preventDefault()

      var payload = {
        name: this.name,
        id: this.id
      }

      this.$http.post('/api/towerfall/new/', payload).then((res) => {
        // Success callback
        console.log('win')
        console.log(res)
        var j = res.json()
        console.log(j)
        this.$route.router.go('/towerfall' + j.redirect)
      }, (res) => {
        console.log('fail')
        console.log(res)
      })
    }
  }
}
</script>

<style lang="scss">
</style>
