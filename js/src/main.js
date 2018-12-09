/* eslint-env browser */

import Vue from 'vue'
import Resource from 'vue-resource'
import Cookie from 'vue-cookie'
import Vue2Filters from 'vue2-filters'
import Toast from 'vue-easy-toast'
import vueHeadful from 'vue-headful'
import * as Icon from 'vue-awesome'
import App from './DrunkenFall'
import {version} from './version'

// Print version info, but only in prod
if (window.location.host === 'drunkenfall.com') {
  console.log(`DrunkenFall ${version}`)
}

import Sidebar from './components/Sidebar.vue'

Vue.use(Resource)
Vue.use(Cookie)
Vue.use(Toast)
Vue.use(Vue2Filters)

Vue.component('icon', Icon)
Vue.component('sidebar', Sidebar)
Vue.component('headful', vueHeadful)

import router from './core/router.js'
import store from './core/store'
import middleware from './core/middleware.js'

middleware(router)

var Root = Vue.extend(App)
new Root({ // eslint-disable-line
  router: router,
  store: store,
}).$mount("#drunkenfall")
