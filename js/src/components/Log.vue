<template>
  <div v-if="tournament">
    <basic-event v-for="event in events" :event="event"></basic-event>
    <div class="clear"></div>
  </div>
</template>

<script>
import BasicEvent from '../components/events/BasicEvent.vue'
import Event from '../models/Event.js'
import _ from 'lodash'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Log',
  mixins: [DrunkenFallMixin],

  components: {
    BasicEvent,
  },

  computed: {
    events () {
      let events = this.tournament.events
      _.each(this.tournament.tryouts, (tryout) => {
        events = _.concat(events, tryout.events)
      })
      _.each(this.tournament.semis, (semi) => {
        events = _.concat(events, semi.events)
      })
      events = _.concat(events, this.tournament.final.events)

      events = _.omitBy(events, _.isNil) // TODO(thiderman): Why are there nil items? hwat
      events = _.sortBy(events, [(o) => { return o.time }])
      events = _.reverse(events)
      events = _.map(events, Event.fromObject)

      return events
    }
  },
}
</script>

<style lang="scss">
@import "../variables.scss";

</style>
