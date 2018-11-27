<template>
<div class="sidebar-buttons" v-if="tournament && user && user.isJudge && showSidebar">
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

    <button-link v-if="user.isProducer"
      :cls="{disabled: tournament.isEnded}"
      :to="{ name: 'participants', params: { tournament: tournament.id }}"
      :icon="'users'" :iconClass="{ warning: tournament.isStarted }" :label="'Players'" />

    <button-link v-if="user.isProducer && tournament.isEnded"
      :to="{ name: 'credits', params: { tournament: tournament.id }}"
      :iconClass="'positive'"
      :icon="'film'"
      :label="'Roll credits'" />

    <button-link v-if="user.isProducer"
      :to="{ name: 'casters', params: { tournament: tournament.id }}"
      :icon="'microphone'"
      :label="'Set casters'" />

      <button-link v-if="user.isProducer && tournament.isTest && tournament.canStart"
        :func="usurp"
        :cls="{ disabled: !tournament.isUsurpable}"
        :iconClass="'warning'"
        :icon="'user-plus'"
        :label="'Add testing players'"
        :tooltip="'Tournament is full.'" />

      <button-link v-if="user.isProducer && tournament.isTest && tournament.isRunning"
        :func="autoplay"
        :iconClass="'warning'"
        :icon="'forward'"
        :label="autoplayLabel" />

      <button-link v-if="user.isJudge && !tournament.isEnded"
        :to="{ name: 'judge', params: { tournament: tournament.id }, query: {fullscreen: 'youhavelostthegame'}}"
        :iconClass="'positive'"
        :icon="'beer'" :label="'Judge'" />
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
      this.api.usurp({ id: this.tournament.id }).then((res) => {
        console.log("usurp response", res)
      }, (err) => {
        this.$alert("Usurp failed. See console.")
        console.error(err)
      })
    },
    autoplay () {
      this.api.autoplay({ id: this.tournament.id }).then((res) => {
        console.log("autoplay response", res)
      }, (err) => {
        this.$alert("Autoplay failed. See console.")
        console.error(err)
      })
    },
    start () {
      this.api.startTournament({ id: this.tournament.id }).then((res) => {
        console.debug("start response:", res)
        this.$router.push({'name': 'tournament', params: {'tournament': this.tournament.id}})
      }, (err) => {
        this.$alert("Start failed. See console.")
        console.error(err)
      })
    },
    next () {
      this.$router.push({name: "match", params: {
        "match": this.tournament.next
      }})
    },
  },

  computed: {
    autoplayLabel () {
      return `Autoplay ${this.tournament.nextMatch.kind}s`
    }
  },

  created () {
    let root = "/api/tournaments/{id}"

    this.api = this.$resource("/api/", {}, {
      startTournament: { method: "GET", url: `${root}/start/` },
      next: { method: "GET", url: `${root}/next/` },
      usurp: { method: "GET", url: `${root}/usurp/` },
      autoplay: { method: "GET", url: `${root}/autoplay/` },
    })
  },
}
</script>
