<template>
  <div v-if="profile">
     <div class="profile">
       <img class="avatar" :src="profile.avatar" :alt="profile.nick" />

       <h1 :class="profile.color">
         {{profile.nick}}
       </h1>

       <div class="stats" :class="profile.color">
         <div class="stat">
           <div class="data">
             <span class="main">2</span>
             <span class="ordinal">{{getOrdinal(2)}}</span>
           </div>
           <p>Rank</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">6</span>
           </div>
           <p>Tournaments</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">97</span>
           </div>
           <p>Matches</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">841</span>
           </div>
           <p>Rounds</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">638</span>
           </div>
           <p>Kills</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">89</span>
           </div>
           <p>Selfs</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">21</span>
           </div>
           <p>Sweeps</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">121</span>
           </div>
           <p>Shots</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">3298</span>
           </div>
           <p>Score</p>
         </div>

         <div class="stat">
           <div class="data">
             <span class="main">1</span>
             <span>h</span>
             <span class="main">37</span>
             <span>m</span>
           </div>
           <p>Play time</p>
         </div>

         <div class="clear"></div>
       </div>
     </div>
  </div>
</template>

<script>
import DrunkenFallMixin from "../mixin"

export default {
  name: "Profile",
  mixins: [DrunkenFallMixin],
  computed: {
    profile () {
      return this.user
    }
  },
  methods: {
    getOrdinal (n) {
      var s = ["th", "st", "nd", "rd"]
      var v = n % 100
      return s[(v - 20) % 10] || s[v] || s[0]
    }
  },
}

</script>

<style lang="scss" scoped>
@import "../variables.scss";

.profile {
  width: 80%;
  margin: 3em auto;
  text-align: center;

  img.avatar {
    display: inline-block;
    object-fit: cover;
    border-radius: 100%;
    margin: 0 auto;
    width:  25%;
    box-shadow: -1px -1px 6px rgba(0,0,0,0.5);
    background-color: rgba(10,12,14,0.3);
    margin-bottom: -30px;
  }

  .stats {

    &.green  {background-color: $green-bg-darkest;}
    &.blue   {background-color: $blue-bg-darkest;}
    &.pink   {background-color: $pink-bg-darkest;}
    &.orange {background-color: $orange-bg-darkest;}
    &.white  {background-color: $white-bg-darkest;}
    &.yellow {background-color: $yellow-bg-darkest;}
    &.cyan   {background-color: $cyan-bg-darkest;}
    &.purple {background-color: $purple-bg-darkest;}
    &.red    {background-color: $red-bg-darkest;}

    width: 80%;
    margin: 3rem auto;
    box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

    .stat {
      width: 20%;
      height: 100%;
      float: left;
      padding-bottom: 1em;

      &:nth-child(odd) {
        background-color: rgba(255,255,255,0.04);
      }

      .data {
        margin-top: 0.5rem;
        height: 5em;
        span.main {
          font-size: 5em;
        }
      }
      p {
        color: #888;
      }
    }
  }
}
</style>
