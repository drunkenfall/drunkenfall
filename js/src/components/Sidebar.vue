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
}
</script>

<style lang="scss" scoped>

#sidebar {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;

  width: 280px;
  background-color: #333339;

  box-shadow: 4px 0px 4px rgba(0,0,0,0.3);

  >a {
    box-shadow: none;
  }

  .logo {
    width: 100%;
  }

  .tournaments {
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
