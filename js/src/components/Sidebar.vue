<template>
  <div>
    <div id="sidebar">
      <router-link
        :to="{name: 'start'}">
        <img class="logo" alt="" src="/static/img/drunkenfall.png"/>
      </router-link>


      <div class="content">
        <div class="tournaments">
          <sidebar-tournament
            class="tournament"
            :tournament="tournament"
            v-for="tournament in tournaments"></sidebar-tournament>
        </div>

        <div v-if="user.isProducer && viewing(['start'])">
          <h1>Actions</h1>
          <div class="actions links">
            <router-link class="action"
              :to="{ name: 'new'}">
              <div class="icon positive">
                <icon name="plus"></icon>
              </div>
              <p>New tournament</p>
              <div class="clear"></div>
            </router-link>

            <a href="/api/facebook/login">
              <div class="icon warning">
                <icon name="facebook"></icon>
              </div>
              <p>Re-facebook</p>
              <div class="clear"></div>
            </a>

            <a @click="clear">
              <div class="icon danger">
                <icon name="trash"></icon>
              </div>
              <p>Clear test tournaments</p>
              <div class="clear"></div>
            </a>
          </div>
        </div>
      </div>

      <div v-if="userLoaded && user.authenticated" class="user">
        <div @click="settings" class="settings">
          <icon name="cog"></icon>
        </div>
        <img :alt="user.firstName" :src="user.avatar"/>

        <div @click="logout" class="logout">
          <icon name="sign-out"></icon>
        </div>

        <h1 id="sidebar-username" :class="user.color">{{user.nick}}</h1>
      </div>

      <div v-if="!viewing(['facebook']) && userLoaded && !user.authenticated" class="content facebook">
        <div class="links">
          <a href="/api/facebook/login">
            <div class="icon">
              <icon name="facebook"></icon>
            </div>
            <p>Facebook login</p>
            <div class="clear"></div>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import _ from "lodash"
import Person from '../models/Person.js'
import SidebarTournament from './SidebarTournament'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Sidebar',
  mixins: [DrunkenFallMixin],
  components: {
    Person,
    SidebarTournament,
  },

  methods: {
    viewing (names) {
      // Returns true if currently viewing any of the route names.
      return _.includes(names, this.$route.name)
    },
    clear (event) {
      let $vue = this
      event.preventDefault()
      return this.$http.get('/api/towerfall/tournament/clear/').then(function (res) {
        console.log(res)
      }, function (res) {
        $vue.$alert("Clearing failed. See console.")
        console.error(res)
        return { tournaments: [] }
      })
    },
    logout () {
      this.user.logout(this)
    },
    settings () {
      this.$router.push("/towerfall/settings")
    },
  },
}
</script>

<style lang="scss">
@import "../variables.scss";

#sidebar {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 100;
  font-size: 1.3em;

  width: $sidebar-width;
  background-color: #333339;
  transition: 0.5s ease-in-out;

  box-shadow: 4px 0px 4px rgba(0,0,0,0.3);

  user-select: none;

  h1 {
    font-size: 1em;
  }

  >a {
    box-shadow: none;
  }

  .logo {
    width: 100%;
  }

  >.content {
    margin: 1.5rem;
  }

  .user {
    height: 140px;
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;

    img {
      height: 75px;
      width:  75px;
      object-fit: cover;
      border-radius: 100%;
      box-shadow: 2px 2px 3px rgba(0,0,0,0.5);
      background-color: rgba(10,12,14,0.3);
      margin: 0 auto;
      display: block;
    }
    h1 {
      transition: 0.5s;
      font-size: 1.5rem;
    }
    div {
      position: absolute;
      top: 30px;
      height: 1.3em;

      color: #999;
      filter: drop-shadow(2px 2px 3px rgba(0,0,0,0.5));
      transition: 0.3s;
      font-size: 1.2em;

      &:hover {
        cursor: pointer;
        color: #dbdbdb;
      }

      &.settings {
        left: 40px;
      }
      &.logout {
        right: 40px;
      }
    }
  }

  .facebook {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;

    >.links {
      margin: 1.5em;
    }
  }
}

</style>
