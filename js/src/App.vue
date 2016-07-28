<template>
</template>

<script>
import Hello from './components/Hello'

export default {
  components: {
    Hello
  }
}
</script>

<style lang="scss">
body, html {
  background-color: #22222a;
  color: #dbdbdb;
  font-family: 'Yanone Kaffeesatz', sans-serif;
  margin: 0;
}
.clear {
  clear: both;
}

h1 {
  text-align: center;
}
a, a:visited {
  color: inherit;
  text-decoration: none;
}


form {
  width: 80%;
  margin: 30px auto;
}
form input.text {
  font-size: 2em;
  /* border-radius: 10px; */
  border: none;
  padding: 10px 30px;
  width: 100%;
  background-color: #333339;
  color: white;
  box-sizing: border-box;
}
form input.submit {
  margin-top: 20px;
  width: 250px;
  border: none;
  font-size: 2em;
  opacity: 0;
  transition: 1.0s;
}
form h2 {
  text-align: center;
}

</style>
