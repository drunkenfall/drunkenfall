<script>
/* eslint-env browser */

import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import store from './core/store'

const proto = window.location.protocol === "https:" ? 'wss://' : "ws://"
const wsURL = `${proto}${window.location.host}/api/auto-updater`
Vue.use(VueNativeSock, wsURL, {
  store: store,
  format: 'json',
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
  reconnectionAttempts: 5, // (Number) number of reconnection attempts before giving up (Infinity),
  reconnectionDelay: 1000, // (Number) how long to initially wait before attempting a new (1000)
})
</script>

<style lang="scss">
@import "./css/main.scss";
</style>
