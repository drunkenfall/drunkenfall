import Vue from 'vue'
import Router from 'vue-router'

import About from '../components/About.vue'
import Admin from '../components/Admin.vue'
import Casters from '../components/Caster.vue'
import Credits from '../components/Credits.vue'
import Disable from '../components/Disable.vue'
import Edit from '../components/Edit.vue'
import Join from '../components/Join.vue'
import JudgeInterface from '../components/Judge.vue'
import Log from '../components/Log.vue'
import Match from '../components/Match.vue'
import New from '../components/New.vue'
import NextScreen from '../components/NextScreen.vue'
import Participants from '../components/Participants.vue'
import Archers from '../components/Archers.vue'
import PostMatch from '../components/PostMatch.vue'
import Runnerups from '../components/Runnerups.vue'
import ScoreScreen from '../components/ScoreScreen.vue'
import Settings from '../components/Settings.vue'
import Statusbar from '../components/Statusbar.vue'
import HUD from '../components/Hud.vue'
import Stream from '../components/Stream'
import TournamentList from '../components/TournamentList.vue'
import TournamentView from '../components/Tournament.vue'

import DrunkenFallNew from '../components/new/DrunkenFall.vue'
import GroupNew from '../components/new/Group.vue'

Vue.use(Router)

var router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'about',
      component: About
    },
    {
      path: '/facebook/finalize',
      name: 'facebook',
      component: Settings
    },
    {
      path: '/tournaments/',
      name: 'tournaments',
      component: TournamentList
    },
    {
      path: '/tournaments/new/',
      name: 'new',
      component: New,
    },
    {
      path: '/tournaments/new/drunkenfall/',
      name: 'newDrunkenfall',
      component: DrunkenFallNew,
    },
    {
      path: '/tournaments/new/group/',
      name: 'newGroup',
      component: GroupNew,
    },
    {
      path: '/settings/',
      name: 'settings',
      component: Settings,
    },
    {
      path: '/archers/',
      name: 'archers',
      component: Archers,
    },
    {
      path: '/archers/:id',
      name: 'archer',
      component: Archers,
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin,
    },
    {
      path: '/admin/disable',
      name: 'disable',
      component: Disable,
    },
    {
      path: '/tournaments/:tournament/',
      name: 'tournament',
      component: TournamentView
    },
    {
      path: '/tournaments/:tournament/join/',
      name: 'join',
      component: Join
    },
    {
      path: '/tournaments/:tournament/participants/',
      name: 'participants',
      component: Participants
    },
    {
      path: '/tournaments/:tournament/runnerups/',
      name: 'runnerups',
      component: Runnerups
    },
    {
      path: '/tournaments/:tournament/casters/',
      name: 'casters',
      component: Casters
    },
    {
      path: '/tournaments/:tournament/edit/',
      name: 'edit',
      component: Edit
    },
    {
      path: '/tournaments/:tournament/scores/',
      name: 'scores',
      component: ScoreScreen
    },
    {
      path: '/tournaments/:tournament/next/',
      name: 'next',
      component: NextScreen
    },
    {
      path: '/tournaments/:tournament/judge/',
      name: 'judge',
      component: JudgeInterface
    },
    {
      path: '/tournaments/:tournament/charts/',
      name: 'charts',
      component: PostMatch
    },
    {
      path: '/tournaments/:tournament/log/',
      name: 'log',
      component: Log
    },
    {
      path: '/live/',
      name: 'live',
      component: Stream
    },
    {
      path: '/status/',
      name: 'status',
      component: Statusbar
    },
    {
      path: '/hud/',
      name: 'hud',
      component: HUD
    },
    {
      path: '/tournaments/:tournament/credits/',
      name: 'credits',
      component: Credits
    },
    {
      path: '/tournaments/:tournament/:match/',
      name: 'match',
      component: Match
    },
    {
      path: '/tournaments/:tournament/*',
      redirect: '/tournaments/:tournament/',
    },
  ],
})

export default router
