<template>
  <div class="sidebar-buttons" v-if="user && user.isJudge">
    <div class="links">
      <button-link v-if="tournament.canStart && user.isCommentator"
        :func="start"
        :iconClass="'positive'"
        :icon="'play'"
        :label="'Start tournament'" />

      <button-link v-if="user.isJudge &&tournament.isRunning && !tournament.shouldBackfill"
        :func="next"
        :iconClass="'positive'"
        :icon="'play'"
        :label="'Next match'" />

      <button-link v-if="user.isJudge"
        :to="{ name: 'log', params: { tournament: tournament.id }}"
        :icon="'book'" :label="'Log'" />

      <button-link v-if="user.isProducer"
        :cls="{disabled: tournament.isEnded}"
        :to="{ name: 'participants', params: { tournament: tournament.id }}"
        :icon="'users'" :iconClass="{ warning: tournament.isStarted }" :label="'Participants'" />

      <button-link v-if="user.isProducer"
        :to="{ name: 'edit', params: { tournament: tournament.id }}"
        :icon="'pencil'" :iconClass="'danger'" :label="'Edit'" />

      <button-link v-if="user.isProducer && tournament.isEnded"
        :to="{ name: 'credits', params: { tournament: tournament.id }}"
        :iconClass="'positive'"
        :icon="'film'"
        :label="'Roll credits'" />

      <button-link v-if="user.isCommentator && tournament.shouldBackfill"
        :to="{ name: 'runnerups', params: { tournament: tournament.id }}"
        :iconClass="'positive'"
        :icon="'cloud-upload'"
        :label="'Backfill semis'" />

      <button-link v-if="user.isProducer && tournament.canShuffle"
        :func="tournament.reshuffle"
        :iconClass="'warning'"
        :icon="'random'"
        :label="'Reshuffle'" />

      <button-link v-if="user.isProducer && tournament.isTest && tournament.canStart"
        :func="usurp"
        :cls="{ disabled: !tournament.isUsurpable}"
        :iconClass="'warning'"
        :icon="'user-plus'"
        :label="'Add testing players'"
        :tooltip="'Tournament is full.'" />

      <div class="maybe-clear"></div>
    </div>
  </div>
</template>

<script>
import DrunkenFallMixin from "../../mixin"
import ButtonLink from "./ButtonLink"

export default {
  name: "TournamentControls",
  mixins: [DrunkenFallMixin],
  components: {
    ButtonLink,
  },
  methods: {
    usurp () {
      return this.tournament.usurp()
    },
    start () {
      return this.tournament.start()
    },
    next () {
      this.$router.push({name: "match", params: {
        "match": this.tournament.next()
      }})
    },
  }
}
</script>
