<template>
  <div>
    <div id="sidebar">
      <router-link
        :to="{name: 'start'}">
        <img class="logo" alt="" src="/static/img/drunkenfall.png"/>
      </router-link>


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

      <div v-if="user" class="user">
        <icon class="settings" name="cog"></icon>
        <img :alt="user.firstName" :src="user.avatar"/>
        <icon class="logout" name="sign-out"></icon>

        <h1 :class="user.color">{{user.nick}}</h1>
      </div>
    </div>
  </div>
</template>

<script>
import _ from "lodash"
import Person from '../models/Person.js'
import SidebarTournament from './SidebarTournament'

export default {
  name: 'Sidebar',
  components: {
    Person,
    SidebarTournament,
  },

  computed: {
    user () {
      return this.$store.state.user
    },
    tournaments () {
      return this.$store.state.tournaments
    },
  },
  methods: {
    viewing (names) {
      // Returns true if currently viewing any of the route names.
      return _.includes(names, this.$route.name)
    },
    clear (event) {
      event.preventDefault()
      return this.$http.get('/api/towerfall/tournament/clear/').then(function (res) {
        console.log(res)
      }, function (res) {
        console.error('error when clearing tournaments', res)
        return { tournaments: [] }
      })
    }
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

  .tournaments, .links{
    margin: 1.5rem;
  }

  .links {
    margin-left: 1em;
    font-size: 1.2em;

    a {
      position: relative;
      background-color: $button-bg;
      text-shadow: 2px 2px 2px rgba(0,0,0,0.3);
      width: 100%;
      display: block;
      margin-bottom: 0.5em;
      transition: 0.3s;

      &:hover {
        background-color: $button-hover-bg;
      }

      .tooltip {
        visibility: hidden;
        opacity: 0;
        position: absolute;
        top: 0.4em;
        left: 105%;
        z-index: 1;

        transition: opacity .15s, visibility .15s;

        color: #dbdbdb;
        background-color: #504040;
        box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

        white-space: nowrap;
        font-size: 0.8em;
      }

      &.router-link-active p {
        color: #708090;
      }

      &.disabled {
        color: #555;
        background-color: #373739;

        .icon {
          background-color: #404040 !important;
        }

        &:hover .tooltip {
          visibility: visible;
          opacity: 1;
        }
      }

      .icon, p {
        padding: 0.5em 0.5em 0.2em;
      }

      .icon {
        float: left;
        width: 1em;
        text-align: center;
        background-color: #405060;

        &.positive {background-color: #406040;}
        &.warning {background-color: #774e2e;}
        &.danger {background-color: #772e2e;}

        .fa-icon {
          filter: drop-shadow(1px 1px 1px rgba(0,0,0,0.3));
        }
      }

      p {
        display: inline-block;

      }
    }
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
      font-size: 1rem;
    }
    .fa-icon {
      position: absolute;
      top: 30px;
      height: 1.3em;

      color: #999;
      filter: drop-shadow(2px 2px 3px rgba(0,0,0,0.5));
      transition: 0.3s;

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
}

</style>
