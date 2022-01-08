import Vue from 'vue'
import Vuex from 'vuex'

import currentUser from './modules/user'
import currentSite from './modules/site'
Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  modules: {
    currentUser,
    currentSite
  },
  strict: debug
})
