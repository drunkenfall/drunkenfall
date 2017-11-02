<template>
  <div>
    <div id="sidebar">
      <div id="hamburger">
        <div class="icon positive">
          <icon name="bars"></icon>
        </div>
      </div>

      <router-link
        :to="{name: 'dispatch'}">
        <img class="logo" alt="" src="/static/img/drunkenfall.png"/>
      </router-link>

      <div class="content main">

        <!-- <div class="tournaments">
        <sidebar-tournament
          class="tournament"
          :tournament="tournament"
          v-for="tournament in tournaments"></sidebar-tournament>
        </div>
        -->

        <div class="blocks">
          <router-link class="action"
            :to="{ name: 'tournaments'}">
            <div class="icon">
              <icon name="fire"></icon>
            </div>
            <p>Tournaments</p>
            <div class="clear"></div>
          </router-link>

          <router-link class="action"
            :to="{ name: 'archers'}">

            <div class="icon">
              🏹
            </div>
            <p>Archers</p>
            <div class="clear"></div>
          </router-link>

          <router-link v-if="user.isProducer" class="action"
            :to="{ name: 'admin'}">
            <div class="icon">
              <icon name="superpowers"></icon>
            </div>
            <p>Admin</p>
            <div class="clear"></div>
          </router-link>
        </div>

        <div v-if="user.isProducer && viewing(['tournaments']) && false">
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
        <img @click="profile" :alt="user.firstName" :src="user.avatar"/>

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
            <p>Sign in</p>
            <div class="clear"></div>
          </a>
        </div>
      </div>

      <div class="clear"></div>
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
      return this.$http.get('/api/tournament/clear/').then(function (res) {
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
      this.$router.push({name: "settings"})
    },
    profile () {
      this.$router.push(`/profile/${this.user.id}`)
    },
  },
  created () {
    document.onkeydown = (e) => {
      if (e.keyCode === 37) {
        let b = document.getElementsByTagName("body")[0]
        if (b.className === "sidebar-less") {
          b.className = ""
        } else {
          b.className = "sidebar-less"
        }
      }
    }
  },
}
</script>

<style lang="scss">
@import "../variables.scss";

#sidebar {
  background-color: $bg-default;
  transition: 0.5s ease-in-out;
  box-shadow: 4px 0px 4px rgba(0,0,0,0.3);
  user-select: none;
  z-index: 100;
  font-size: 1.3em;

  .blocks .action {
    transition: 0.2s;
    text-shadow: 2px 2px 3px rgba(0,0,0,0.3);
    color: #888;
    display: block;
    width: 100%;
    text-align: center;
    /* background-color: rgba(0,0,0,0.2); */
    /* border-bottom: 1px solid rgba(0,0,0,0.1); */
    /* border-top: 1px solid rgba(255,255,255,0.1); */
    padding: 1.3rem 0;
    box-shadow: none;

    &.router-link-active {
      background-color: rgba(0,0,0,0.2);
      color: $fg-default;
      .icon {
        color: $accent;
      }
    }

    .icon {
      font-size: 3em;
      text-shadow: 3px 3px 5px rgba(0,0,0,0.5);
    }
    p {
      margin-top: 1em;
      font-weight: 100;
      font-size: 0.8em;
      text-transform: uppercase;
    }
  }

  /* Real devices */
  @media screen and (min-device-width: 770px) {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 100;
    font-size: 1.3em;
    width: $sidebar-width;

    .logo {
      width: 100%;
    }

    #hamburger {
      display: none;
    }

    >.content {
      /* margin: 1.5rem; */
    }
  }

  /* Small devices */
  @media screen and (max-device-width: 769px) {
    position: relative;
    width: 100%;
    height: 5em;

    .logo {
      /* width: 40%; */
      height: 100%;
      float: left;
      display: block;
      margin: 0px auto;
    }

    >.content.main {
      display: none;
    }

    #hamburger {
      display: none;
      cursor: pointer;
      float: left;
      font-size: 1em;
      margin: 1em;
    }
  }

  h1 {
    font-size: 1em;
  }

  >a {
    box-shadow: none;
  }

  @media screen and (min-device-width: 770px) {
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
        cursor: pointer;
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
          color: $fg-default;
        }

        &.settings {
          left: 30px;
        }
        &.logout {
          right: 30px;
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


  @media screen and (max-device-width: 769px) {
    .user {
      float: right;
      width: 30%;

      img {
        position: absolute;
        right: 5%;
        top: 50%;
        transform: translateY(-50%);
        cursor: pointer;

        height: 60%;
        object-fit: cover;
        border-radius: 100%;
        box-shadow: $shadow-default;
        background-color: rgba(10,12,14,0.3);
        display: block;
      }
      h1 {
        display: none;
      }
      div {
        display: none;
      }
    }

    .facebook {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);

      >.links {
        margin: 1rem;
      }
    }
  }
}
</style>
