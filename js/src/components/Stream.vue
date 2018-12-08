<template>
<div id="live" v-if="tournament">
  <headful :title="tournament.subtitle + ' / Stream'"></headful>
  <div class="stream-sidebar">
    <div class="title subtitle-logo">
      <img alt="" src="/static/img/oem.svg"/>
      <div class="text">
        <p class="header">DrunkenFall</p>
        <p class="subtitle" :class="tournament.color">{{tournament.subtitle}}</p>
      </div>
    </div>
    <div class="casters cam">
      <div></div>
    </div>

    <div class="casternames" v-if="tournament.casters.length === 2">
      <div class="caster first">
        {{tournament.casters[1].displayName}}
      </div>
      <div class="amp">&amp;</div>
      <div class="caster">
        {{tournament.casters[0].displayName}}
      </div>
    </div>

    <div class="casternames" v-if="tournament.casters.length === 1">
      <div class="sole caster">
        {{tournament.casters[0].displayName}}
      </div>
    </div>

    <div class="players cam">
      <div></div>
    </div>
    <div class="match">
      <div class="name">{{match.title}}</div>
      <div class="level">
        {{match.levelTitle}}
        -
        Match to {{match.endScore}} points
      </div>
      <div class="round">Round {{round}}</div>
    </div>
  </div>

  <div class="stream-main">
    <div class="game">
      <div>
      </div>
    </div>

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
          <div class="person">
            <div class="stats">
              <div class="kills">
                <div class="emoji">
                  <img alt="coin" src="/static/img/coin.png"/>
                </div>
                <div class="number">{{p.kills}}</div>
              </div>
              <div class="shots">
                <div class="emoji">
                  <img alt="bullet" src="/static/img/bullet.png"/>
                </div>
                <div class="number">{{p.shots}}</div>
              </div>
            </div>
            <p class="name" :class="p.color">
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

export default {
  name: 'Stream',
  mixins: [DrunkenFallMixin],
  computed: {
    tournament () {
      if (this.$route.params.tournament) {
        return this.tournaments[this.$route.params.tournament]
      }
      return this.trackingTournament
    },
    match () {
      return this.tournament.currentMatch
    }
  },
  methods: {
    startSimulation () {
      console.log("starting")

      this.api.start({ id: this.tournament.id }).then((res) => {
        console.log("start response:", res)
      }, (err) => {
        console.error(err)
      })
    },

    stopSimulation () {
      console.log("stopping")

      this.api.stop({ id: this.tournament.id }).then((res) => {
        console.log("stop response:", res)
      }, (err) => {
        console.error(err)
      })
    },
  },
  created () {
    this.loadAll()

    document.getElementsByTagName("body")[0].className = "scroll-less sidebar-less"
    this.api = this.$resource("/api", {}, {
      start: { method: "GET", url: "/api/simulator/start/{id}" },
      stop: { method: "GET", url: "/api/simulator/stop/{id}" },
    })
  },
}
</script>


<style lang="scss">
@import "../css/colors.scss";
$stream-sidebar: 580px;
$bottom: 75px;
$chromakey: #0f0;
/* $chromakey: #212; */

#live {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row !important;

  .stream-sidebar {
    display: flex;
    width: $stream-sidebar;
    flex-direction: column;
    /* background-color: rgba(0,0,10,0.1); */

    >div {
      width: 100%;
    }

    .cam {
      height: 326px;
      background-color: $chromakey;
      display: flex;
      align-items: center;
      justify-content: center;

      div {
        color: rgba(255,255,255,0.3);
        font-size: 2em;
      }
    }

    .casternames {
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 1em;

      .caster {
        font-size: 2.5em;
        padding: 0.5em;
        width: 45%;
        &.first {
          text-align: right;
        }
        &.sole {
          text-align: center;
        }
      }
      .amp {
        color: $fg-secondary;
        font-size: 1.5em;
        text-align: center;
        display: inline-block;
        margin: 0 0.5em;
        width: 1%
      }
    }

    .match {
      height: $bottom;
      text-align: center;

      .name {
        margin-top: 0.5em;
        @include display2();
      }
      .level {
        @include title();
        color: $fg-secondary;
        margin: 0.3em;
      }
      .round {
        @include display1();
      }
    }
  }

  .stream-main {
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    flex-grow: 1;

    .game {
      width: 100%;
      flex-grow: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: $chromakey;

      div {
        color: rgba(255,255,255,0.3);
        /* font-size: 2em; */
      }
    }

    .statusbar {
      height: $bottom;
      display: flex;

      .player {
        display: flex;
        height: 100%;
        width: 25%;

        .avatar {
          height: $bottom;
          width: $bottom;
          overflow: hidden;
          transition: 0.3s;
          position: relative;

          img {
            height: $bottom;
            width: $bottom;
            object-fit: cover;
            z-index: 10;

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

            font-size: 30px;
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
          flex-direction: column;
          padding: 0.2em 0.1em;

          .status {
            min-height: 50%;
            display: flex;
            flex-direction: row;
            padding-left: 0.1em;

            .orbs, .shield {
              width: 22px;
              img {
                height: 26px;
              }
            }
            .arrows {
              img {
                height: 26px;
                margin-right: -15px;
              }
            }

            .reason {
              width: 100%;
              text-align: center;
              font-size: 20px;
              display: flex;
              justify-content: center;
              align-items: center;
            }

          }
          .person {
            display: flex;
            justify-content: space-between;
            line-height: 100%;
            padding-left: 0.2em;
            padding-bottom: 15px;

            .stats {
              order: 2;
              /* width: 75px; */
              display: flex;

              .kills, .shots {
                display: flex;
                padding: 0 3px;

                .emoji {
                  font-size: 0.5em;
                  img {
                    margin-top: 5px;
                    height: 20px;
                    /* width: 20px; */
                  }
                }
                .number {
                  width: 25px;
                  text-align: center;
                  /* margin-right: 0.5em; */
                }
              }
            }

            .name {
              font-size: 28px;
              order: 1
            }
          }
        }
      }
    }
  }
}
</style>
