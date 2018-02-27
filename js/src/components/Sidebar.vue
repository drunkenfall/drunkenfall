<template>
<div v-if="showSidebar">
  <div id="sidebar" :class="{visible: hamburgerActive}">
    <router-link :to="dispatchLink" class="logo">
      <div>
        <img alt="" src="/static/img/oem.svg"/>
      </div>
    </router-link>

    <div @click="toggle" class="hamburger">
      <icon name="bars" />
    </div>

    <div class="blocks" @click="toggle" :class="{visible: hamburgerActive}">
      <router-link class="action" :to="{ name: 'tournaments'}">
        <div class="icon">🔥</div>
        <p>Tournaments</p>
      </router-link>

      <router-link class="action" :to="{ name: 'archers'}">
        <div class="icon">🏹</div>
        <p>Archers</p>
      </router-link>

      <router-link v-if="user.isProducer" class="action" :to="{ name: 'admin'}">
        <div class="icon">💪</div>
        <p>Superpowers</p>
      </router-link>
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

    <div class="clear"></div>
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

  data () {
    return {
      hamburgerActive: false,
    }
  },

  methods: {
    toggle () {
      this.hamburgerActive = !this.hamburgerActive
    },
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
  justify-content: space-between;
  align-items: space-between;

  /* Real devices */
  @media screen and ($desktop: $desktop-width) {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 100;
    font-size: 1.3em;
    width: $sidebar-width;

    .hamburger {
      display: none;
    }

    .logo {
      padding: 1em;
      width: 50%;
      margin: 0px auto;
    }

    .blocks {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      flex-grow: 1;

      .action {
        @include headline();
        transition: 0.2s;
        text-shadow: 2px 2px 3px rgba(0,0,0,0.3);
        color: #888;

        width: 100%;
        text-align: center;
        padding: 2rem 0;
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
    }
  }

  /* Small devices */
  @media screen and ($device: $device-width) {
    position: relative;
    width: 100%;
    min-height: 3em;
    transform: 0.3s;
    flex-direction: row;
    flex-wrap: wrap;

    .visible {
      min-height: 5em;
    }

    .logo, .hamburger, .user {
      width: 33%;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0.5em 0;
    }

    .hamburger {
      font-size: 1.5em;
      order: 1;
      cursor: pointer;
    }

    .logo {
      order: 2;

      img {
        margin: 0 auto;
        height: 2em;
      }
    }

    .user {
      order: 3;
    }

    .blocks {
      order: 4;
      transform: 0.3s;
      display: none;
      opacity: 0;

      &.visible {
        opacity: 1;
        display: flex;
        width: 100%;
        justify-content: center;
        align-items: center;

        .action {
          flex-grow: 1;
          text-align: center;
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
