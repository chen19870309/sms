// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from '@/store'

import less from 'less'

import iView from 'iview'
import 'iview/dist/styles/iview.css'

import './assets/css/webapp.css'

Vue.use(iView)
Vue.use(less)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
