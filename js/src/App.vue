<script>
/* eslint-env browser */
import _ from 'lodash'
import DrunkenFallMixin from "./mixin.js"

export default {
  data () {
    return {
      // The main websocket object
      ws: null,
      reconnections: 0,
    }
  },
  mixins: [DrunkenFallMixin],
  methods: {
    reconnect: function () {
      this.ws = null

      if (this.reconnections < 30) {
        console.warn(`closed uncleanly, reconnecting (try number ${this.reconnections})`)
        setTimeout(() => {
          this.reconnections = this.reconnections + 1
          this.connect()
        }, 500)
      } else {
        console.warn("Tried too many times, stopping.")
        this.reconnections = 0
      }
    },

    // (Re-)Connect the websocket.
    // Is safe to run when the connection is already up - then it will be a noop.
    connect: function () {
      let $vue = this
      if (!this.ws) {
        let proto = window.location.protocol === "https:" ? 'wss://' : "ws://"
        console.log(`Setting up new ${proto} websocket`)
        this.$set(
          this.$data,
          'ws',
          new WebSocket(proto + window.location.host + '/api/auto-updater'),
        )

        let timeoutId

        this.ws.onmessage = (event) => {
          let res
          try {
            res = JSON.parse(event.data)
          } catch (e) {
            $vue.$alert("Data could not be parsed.")
            console.error("Failed to parse message data:", e)
            return
          }

          // If p is set, this is a ping message that only serves to keep the connection open.
          // Break immediately.
          if (res.p !== undefined) {
            if (timeoutId) {
              clearTimeout(timeoutId)
            }

            timeoutId = setTimeout(() => {
              if (this.ws && this.ws.readyState === 0) {
                // we're already trying to reconnect, don't try more
              } else {
                // the server probably died, try to connect again
                $vue.$warn("Connection lost. Trying to reconnect a few times.")
                this.reconnect()
              }
            }, 10000)
            console.debug("Ping")
            return
          }

          if (res.data) {
            if (res.data.tournaments) {
              // The main bulk update. This contains the latest state.
              $vue.$store.commit('updateAll', {
                '$vue': $vue,
                'tournaments': res.data.tournaments,
              })
              return
            }

            console.log('no tournaments received, did not update anything')
            return
          }

          console.log('Unknown websocket update:', res)
        }

        this.ws.onopen = () => {
          this.reconnections = 0
          console.debug("websocket connected:", this.ws)
        }
        this.ws.onerror = (errorEvent) => { console.error("websocket error:", errorEvent) }
        this.ws.onclose = (closeEvent) => {
          $vue.$warn("Connection closed.")
          console.debug("websocket closed", closeEvent)
          if (!closeEvent.wasClean) {
            this.reconnect()
          }
        }
      }
    },

    get: function (tid) {
      return _.find(this.tournaments, { id: tid })
    }
  }
}
</script>

<style lang="scss">
@import "./css/main.scss";
</style>
