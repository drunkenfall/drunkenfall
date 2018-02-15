<template>
  <div v-if="tournament">
    <headful :title="tournament.subtitle + ' / Edit - DrunkenFall'"></headful>
    <form v-on:submit="edit">
      <textarea v-model="data"
        autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false"
        cols="30" id="" name="" rows="50"></textarea>

      <input id="doit" type="submit" value="Edit tournament"/>
    </form>
  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Edit',
  mixins: [DrunkenFallMixin],

  data () {
    return {
      data: "",
    }
  },

  methods: {
    edit (e) {
      let $vue = this
      e.preventDefault()
      this.api.edit({ id: this.tournament.id }, this.data).then((res) => {
        console.log("edit response:", res)
        this.$router.push({name: "tournament", params: {id: this.tournament.id}})
      }, (err) => {
        $vue.$alert(`editing tournament failed`)
        console.error(err)
      })
    },
  },

  created () {
    this.api = this.$resource("/", {}, {
      edit: { method: "POST", url: "/api/{id}/edit/" },
    })
  },
  mounted () {
    this.$set(this.$data, 'data', JSON.stringify(this.tournament.raw, null, 2))
  },
  watch: {
    tournament (val) {
      this.$set(this.$data, 'data', JSON.stringify(val.raw, null, 2))
    }
  },
}
</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.sidebared-content {
  text-align: center;
}

h1 {
  font-size: 4em;
}

p {
  font-size: 1.5em;
  width: 25em;
  margin: 1em auto;
}

label {
  font-size: 2.5em;
}

#doit {
  margin: 1em auto;
  width: 350px;
  border: none;
  font-size: 2em;
  transition: 1.0s;
  color: $fg-default;
  background-color: $secondary;
  padding: 0.4em 0.8em;
}

textarea {
  text-align: left;
  color: $fg-default;
  background-color: $bg-default;
  width: 90%;
  height: 90%;
  box-shadow: inset 2px 2px 2px 0px rgba(0,0,0,0.5);
  outline: none;
  border: none;
  padding: 1em;
}

</style>
