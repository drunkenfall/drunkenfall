/* eslint-env browser */

import Vue from 'vue'
// import VueNativeSock from 'vue-native-websocket'
import Resource from 'vue-resource'
import Cookie from 'vue-cookie'
import Vue2Filters from 'vue2-filters'
import Toast from 'vue-easy-toast'
import vueHeadful from 'vue-headful'
import * as Icon from 'vue-awesome'
import App from './App'

import store from './core/store'

import Sidebar from './components/Sidebar.vue'

Vue.use(Resource)
Vue.use(Cookie)
Vue.use(Toast)
Vue.use(Vue2Filters)

Vue.component('icon', Icon)
Vue.component('sidebar', Sidebar)
Vue.component('headful', vueHeadful)

import router from './core/router.js'
import middleware from './core/middleware.js'

middleware(router)

var Root = Vue.extend(App)
new Root({ // eslint-disable-line
  router: router,
  store: store,
}).$mount("#app")
