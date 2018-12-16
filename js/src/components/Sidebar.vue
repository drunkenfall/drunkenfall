<template>

<div id="sidebar" :class="{visible: hamburgerActive}" v-if="showSidebar">
  <router-link :to="dispatchLink" class="logo">
    <div>
      <img :class="{ded: !isConnected}" alt="One-Eye" src="/static/img/oem.svg"/>
    </div>
  </router-link>

  <div class="blocks">
    <router-link class="action" :to="{ name: 'tournaments'}">
      <div class="icon">
        <icon name="gamepad"></icon>
      </div>
      <p>Tournament</p>
    </router-link>

    <router-link class="action" :to="{ name: 'archers'}">
      <div class="icon">
        <icon name="trophy"></icon>
      </div>
      <p>Leaderboard</p>
    </router-link>

    <router-link class="action" :to="{ name: 'rules'}">
      <div class="icon">
        <icon name="list"></icon>
      </div>
      <p>Rules</p>
    </router-link>

    <router-link class="action" :to="{ name: 'about'}">
      <div class="icon">
        <icon name="question"></icon>
      </div>
      <p>About</p>
    </router-link>

    <!-- <router-link v-if="user.isProducer" class="action" :to="{ name: 'admin'}">
    <div class="icon">
      <icon name="balance-scale"></icon>
    </div>
    <p>Superpowers</p>
    </router-link> -->

  </div>

  <div v-if="userLoaded && user.authenticated" class="user">
    <div class="controls">
      <div @click="settings" class="settings">
        <icon name="cog"></icon>
      </div>
      <img @click="settings" :alt="user.firstName" :src="user.avatar"/>

      <div @click="logout" class="logout">
        <icon name="sign-out"></icon>
      </div>
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
</div>

</template>

<script>
import _ from "lodash"
import Person from '../models/Person.js'
import DrunkenFallMixin from "../mixin"

export default {
  name: 'Sidebar',
  mixins: [DrunkenFallMixin],
  components: {
    Person,
  },

  methods: {
    viewing (names) {
      // Returns true if currently viewing any of the route names.
      return _.includes(names, this.$route.name)
    },
    clear (event) {
      let $vue = this
      event.preventDefault()
      return this.$http.delete('/api/tournaments/').then(function (res) {
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
  computed: {
    dispatchLink () {
      return {name: "about"}
    }
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
  mounted () {
    if (!this.showSidebar) {
      document.getElementsByTagName("body")[0].className += " sidebar-less"
    }
  }
}
</script>

<style lang="scss">
@import "../css/colors.scss";
@import "../css/main.scss";

#sidebar {
  background-color: $bg-default;
  transition: 0.5s ease-in-out;
  box-shadow: 4px 0px 4px rgba(0,0,0,0.3);
  user-select: none;
  z-index: 100;
  font-size: 1.3em;

  display: flex;
  flex-direction: column;

  /* Real devices */
  @media screen and ($desktop: $desktop-width) {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 100;
    font-size: 1.3em;
    width: $sidebar-width;

    justify-content: space-between;
    align-items: space-between;

    .hamburger {
      display: none;
    }

    .logo {
      padding: 1em;
      width: 50%;
      margin: 0px auto;

      img {
        transition: 1s;

        &.ded {
          -webkit-filter: grayscale(100%) brightness(75%);
          filter: grayscale(100%) brightness(75%);
        }
      }
    }

    .blocks {
      display: flex;
      flex-direction: column;
      justify-content: space-around;
      flex-grow: 1;

      .action {
        @include headline();
        transition: 0.2s;
        text-shadow: 2px 2px 3px rgba(0,0,0,0.3);
        color: #888;
        flex-grow: 1;
        min-height: 100px;

        width: 100%;
        text-align: center;
        box-shadow: none;

        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;

        &.router-link-exact-active {
          background-color: rgba(0,0,0,0.2);
          color: $fg-default;
          .icon {
            color: $accent;
          }
        }

        .icon {
          font-size: 2em;
          text-shadow: 3px 3px 5px rgba(0,0,0,0.5);
        }
        p {
          font-weight: 100;
          font-size: 0.8em;
          text-transform: uppercase;
        }
      }
    }
  }

  /* Small devices */
  @media screen and ($device: $device-width) {
    position: relative;
    width: 100%;
    height: 130px;
    min-height: 130px;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: space-between;

    .logo, .user {
      width: 50%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 0.5em;
    }

    .logo {
      order: 1;

      img {
        margin: 0 auto;
        height: 2em;
        transform: rotate(90deg);
        transition: 0.7s ease-in-out;
      }
    }

    .user {
      order: 2;
      justify-content: flex-end;
    }

    .blocks {
      order: 4;
      transform: 0.3s;

      background-color: rgba(0,0,0,0.1);
      box-shadow: 2px 2px 3px rgba(0,0,0,0.3);

      font-family: "Lato";
      font-size: 10px;

      display: flex;
      width: 100%;
      justify-content: space-around;
      align-items: center;
      text-align: center;

      border-top: 1px solid rgba(255,255,255,0.1);

      .action {
        flex-grow: 1;
        flex-basis: 0;
        opacity: 0.6;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;

        &.router-link-exact-active {
          opacity: 1;
          background-color: rgba(0,0,0,0.2);
          color: $accent;
          .icon {
            color: $accent;
          }
        }
      }
    }
  }

  h1 {
    font-size: 1em;
  }

  >a {
    box-shadow: none;
  }

  @media screen and ($desktop: $desktop-width) {
    .user {
      margin-top: 1em;
      height: 140px;
      display: flex;
      flex-direction: column;

      .controls {
        display: flex;
        justify-content: center;
        align-items: center;
      }

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
        width: 100%;
      }

      .settings, .logout {
        height: 1.3em;
        display: flex;
        flex-grow: 1;
        flex-direction: column;
        justify-content: center;
        align-items: center;

        color: #999;
        filter: drop-shadow(2px 2px 3px rgba(0,0,0,0.5));
        transition: 0.3s;
        font-size: 1.2em;

        &:hover {
          cursor: pointer;
          color: $fg-default;
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


  @media screen and ($device: $device-width) {
    .user {
      img {
        cursor: pointer;
        height: 2em;
        object-fit: cover;
        border-radius: 100%;
        box-shadow: $shadow-default;
        background-color: rgba(10,12,14,0.3);
      }
      h1 {
        display: none;
      }
      .settings, .logout {
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
