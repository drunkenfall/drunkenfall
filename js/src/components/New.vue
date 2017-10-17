<template>
  <div>
    <h1>New tournament</h1>
    <h3>Choose a kind</h3>

    <div class="kind">
      <div class="icon">
        <h1>DrunkenFall</h1>
        <p>Business as usual</p>
      </div>
      <img alt="To The Moon" src="/static/img/moon.jpg"/>
      <div class="clear"></div>
    </div>

    <div class="kind group">
      <div class="icon">
        <h1>Group</h1>
        <p>On-the-fly tournaments for a group of people</p>
      </div>
      <img alt="To The Unomaly" src="/static/img/group.jpg"/>
      <div class="clear"></div>
    </div>


    <div class="kind test">
      <div class="icon">
        <h1>Test</h1>
        <p>Faky fake test tournaments for testy testing</p>
      </div>
      <img alt="To The Test" src="/static/img/test.jpg"/>
      <div class="clear"></div>
    </div>

  </div>
</template>

<script>
export default {
  name: 'New',
}
</script>

<style lang="scss">
@import "../variables.scss";
h1 {
  font-size: 3em;
}

.kind {
  width: 80%;
  margin: 3em auto;
  background-color: #333339;
  box-shadow: 2px 2px 3px rgba(0,0,0,0.3);
  cursor: pointer;
  overflow: hidden;
  height: 12em;

  &.test {
    .icon {
      background-color: $danger;
    }
  }

  &.group {
    .icon {
      background-color: $green-bg;
    }
  }

  .icon {
    text-align: center;
    background-color: #405060;
    width: 20%;
    float: left;
    height: 100%;

    h1 {
      font-size: 2em;
    }

    p {
      margin: 1em 3em 2em 3em;
    }
  }

  img {
    position: relative;
    background-blend-mode: multiply;
    top: 0;
    bottom: 0;
    height: 100%;
    width: 80%;
    float: left;
    object-fit: cover;
    -webkit-filter: grayscale(100%); /* Safari 6.0 - 9.0 */
    filter: grayscale(100%);
  }
}

</style>
