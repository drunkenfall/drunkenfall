<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' / Log - DrunkenFall'"></headful>
    <basic-event v-for="event in tournament.events" :event="event"></basic-event>
    <div class="clear"></div>
  </div>
</template>

<script>
import BasicEvent from '../components/events/BasicEvent.vue'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Log',
  mixins: [DrunkenFallMixin],

  components: {
    BasicEvent,
  },
}
</script>

<style lang="scss">
@import "../css/colors.scss";

</style>
