import Vue from 'vue'
import Router from 'vue-router'
import Resource from 'vue-resource'
import App from './App'

import TournamentList from './components/TournamentList.vue'
import Tournament from './components/Tournament.vue'
import New from './components/New.vue'
import Join from './components/Join.vue'
import Match from './components/Match.vue'
import Facebook from './components/Facebook.vue'
import FacebookFinalize from './components/FacebookFinalize.vue'

// install router
Vue.use(Router)
Vue.use(Resource)

// routing
var router = new Router({
  'hashbang': false,
  'history': true
})

router.map({
  '/facebook/': {
    component: Facebook,
    name: 'facebook',
  },
  '/facebook/finalize': {
    component: FacebookFinalize
  },
  '/towerfall/': {
    component: TournamentList
  },
  '/towerfall/new/': {
    name: 'new',
    component: New
  },
  '/towerfall/:tournament/': {
    name: 'tournament',
    component: Tournament
  },
  '/towerfall/:tournament/join/': {
    name: 'join',
    component: Join
  },
  '/towerfall/:tournament/:kind/:match/': {
    name: 'match',
    component: Match
  }
})

router.beforeEach(function () {
  window.scrollTo(0, 0)

  router.app.connect()
})

// As long as we only have Drunken TowerFall on drunkenfall.com, we should
// always redirect to the towerfall app right away.
router.redirect({
  '/': '/towerfall/'
})

router.start(App, 'app')
