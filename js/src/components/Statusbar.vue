<template>
<div v-if="tournament">
  <div class="statusbar">
    <div v-for="(p, idx) in match.players" class="player">
      <div class="avatar">
        <img :alt="p.person.nick" :src="p.person.avatar" :class="{dead: !p.state.alive}"/>
        <div :class="{dead: !p.state.alive}">
          <p>{{p.state.death}}</p>
        </div>
      </div>

      <div class="data">
        <div class="status">
          <div class="gamestats">
            <div class="orbs" v-if="p.state.alive && p.state.lava">
              <img alt="" :src="lavaOrbImage()"/>
            </div>
            <div class="shield" v-if="p.state.alive && p.state.shield">
              <img alt="" :src="shieldImage()"/>
            </div>
            <div class="arrows" v-if="p.state.alive">
              <img v-for="a in p.state.arrows" alt="" :src="arrowImage(a)"/>
            </div>
            <div class="reason" v-else>
              <p>
                Killed by
                <span v-if="p.state.killer === -1">
                  the level, lol
                </span>
                <span v-else-if="p.state.killer !== idx" :class="match.players[p.state.killer].color">
                  {{match.players[p.state.killer].displayName}}
                </span>
                <span v-else>
                  suicide :'(
                </span>
              </p>
            </div>

          </div>

          <div class="points">
            <div v-for="n in match.endScore"
              class="scores"
              :class="bulletClass(p, idx, n)">
              <p>{{bulletDisplay(p, idx, n)}}</p>
            </div>
          </div>

        </div>
        <div class="person">
          <div class="stats">
            <div class="kills">
              <div class="emoji">💰</div>
              <div class="number">{{p.kills}}</div>
            </div>
            <div class="shots">
              <div class="emoji">🥃</div>
              <div class="number">{{p.shots}}</div>
            </div>
          </div>
          <div class="name">
            <p :class="p.color">
              {{p.person.displayName}}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import DrunkenFallMixin from "../mixin"
import NextScreen from './NextScreen'

export default {
  name: 'Statusbar',
  mixins: [DrunkenFallMixin],
  components: {
    NextScreen,
  },
  computed: {
    tournament () {
      return this.runningTournament
    },
    match () {
      return this.tournament.matches[this.tournament.current]
    }
  },
  created () {
    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"
  },
  methods: {
    bulletClass (player, playerIndex, n) {
      let c = this.match.endScore === 20 ? "final " : ""

      if (n === this.match.endScore && player.kills > this.match.endScore) {
        return c + 'overkill'
      } else if (n > player.kills && n <= player.kills) {
        return c + 'up'
      } else if (n === player.kills && this.downs === -1) {
        return c + 'down'
      } else if (n <= player.kills) {
        return c + 'kill'
      }
      return c
    },

    bulletDisplay (player, playerIndex, n) {
      // Change the display to +1 for every overkill the player makes
      // at the end of the match.
      if (n === this.match.endScore && player.kills > this.match.endScore) {
        return "+" + (player.kills - this.match.endScore)
      }
      return n
    },
  },

}

</script>

<style lang="scss" scoped>
@import "../css/colors.scss";

.statusbar {

  height: 100%;
  display: flex;
  flex-direction: column;

  .player {
    width: 100%;
    height: 25%;

    &:nth-child(odd) {
      background-color: rgba(0,0,0,0.1)
    }

    box-sizing: border-box;
    padding: 20px 40px;

    display: flex;
    flex-direction: row;

    .avatar {
      height: 230px;
      width: 230px;
      overflow: hidden;
      transition: 0.3s;
      position: relative;

      img {
        border-radius: 100%;
        height: 230px;
        width: 230px;
        object-fit: cover;

        &.dead {
          filter: brightness(33%) blur(1px) grayscale(75%);
        }
      }

      div {
        opacity: 0;
        transition: 0.2s;
        z-index: 100;

        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;

        display: flex;
        justify-content: center;
        align-items: center;

        font-size: 3em;
        text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
        color: #aaa;

        &.dead {
          opacity: 1;
        }
      }
    }

    .data {
      display: flex;
      flex-grow: 1;
      flex-direction: row;
      padding: 0.2em 0.1em;

      .status {
        margin-left: 20px;
        display: flex;
        flex-direction: column;
        flex-grow: 1;

        .gamestats {
          width: 100%;
          min-height: 50%;
          display: flex;
          flex-direction: row;
          padding-left: 0.1em;

          img {
            image-rendering: pixelated;
          }

          .orbs, .shield {
            img {
              height: 100px;
            }
          }
          .arrows {
            flex-grow: 1;
            img {
              height: 100px;
              margin-right: -20px;
            }
          }

          .reason {
            width: 100%;
            text-align: center;
            font-size: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 2em;
          }
        }
      }

      .points {
        width: 100%;
        height: 50%;
        display: flex;
        flex-direction: row;
        justify-content: space-around;
        align-items: center;

        .scores {
          display: flex;
          justify-content: center;
          align-items: center;
          text-align: center;

          background-color: $bg-default;
          border-radius: 100%;
          color: $fg-disabled;

          font-size: 1.5em;
          height: 100px;
          width: 100px;

          &.final {
            width: 50px;
            border-radius: 30px;
          }

          &.kill {
            background-color: $fg-default;
            p {color: $bg-bottom;}
          }
          &.overkill {
            background-color: #daa520;
            text-shadow: 2px 2px 2px rgba(0,0,0,0.5);
            p {color: #fff;}
          }
          &.up {
            background-color: #508850;
            p {color: #fff;}
          }
          &.down {
            background-color: #885050;
            p {color: #fff;}
          }
        }
      }

      .person {
        display: flex;
        flex-direction: column;
        width: 300px;
        justify-content: space-between;
        line-height: 100%;
        padding-left: 0.2em;
        padding-bottom: 15px;

        .stats {
          order: 2;
          height: 25%;

          display: flex;
          justify-content: center;
          align-items: center;

          .kills, .shots {
            display: flex;
            /* padding: 0 0.5em; */
            font-size: 2em;

            .emoji {
              font-size: 0.5em;
            }
            .number {
              width: 1em;
              text-align: center;
              font-size: 1.2em;
            }
          }
        }

        .name {
          order: 1;
          height: 75%;

          display: flex;
          justify-content: center;
          align-items: center;
          text-align: center;
          line-height: 200%;

          p {
            font-size: 2.5em;
          }
        }
      }
    }
  }
}
</style>
