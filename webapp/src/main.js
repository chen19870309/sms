// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from '@/store'
import md5 from 'js-md5'
import less from 'less'
import iView from 'iview'
import VueCropper from 'vue-cropper'
// import 'vue-cropper/dist/index.css'
import 'iview/dist/styles/iview.css'

import './assets/css/webapp.css'

import VueParticles from 'vue-particles'
Vue.use(VueParticles)

Vue.use(iView)
Vue.use(less)
Vue.use(VueCropper)

Vue.config.productionTip = false
Vue.prototype.$md5 = md5

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
