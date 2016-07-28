<template>
  <router-view
    keep-alive
    transition>
  </router-view>
</template>

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
p {
  margin: 0;
}

h1 {
  text-align: center;
}
a, a:visited {
  color: inherit;
  text-decoration: none;
}

header {
  background-color: #333339;
  top: 0px;
  left: 0px;
  right: 0px;

  .content {
    margin: 0 20px;
  }
  .title {
    float: left;
    height: 2em;
    font-size: 2em;
    line-height: 2em;
    font-weight: bold;
  }
  .links {
    float: right;
    a {
      margin: 10px !important;
      font-size: 24px;
      float: right;
      background-color: #405060;
      color: #dbdbdb;
      display: block;
      font-weight: bold;
      padding: 7px 30px;
      text-align: center;
      text-decoration: none;
      margin: 10px auto;
      min-width: 100px;
    }

  }
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
