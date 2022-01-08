import * as types from '../mutation-types'
import storage from '@/utils/storage'
const state = {
  siteinfo: {
    name: '煎鱼',
    url: 'https://eddycjy.com',
    beian: '粤ICP备19076352号'
  }
}

// getters
const getters = {
  currentSite: state => state.siteinfo
}
// actions
const actions = {
  saveSite ({commit}, site) {
    commit(types.SAVE_SITE, site)
  },
  fetchSite ({commit}) {
    commit(types.FETCH_SITE)
  },
  deleteSite ({commit}) {
    commit(types.DELETE_SITE)
  }
}
// mutations
const mutations = {
  [types.CREATE_USER] (state, user) {
    state.user = user
    storage.set('current_site', user)
  },
  [types.FETCH_USER] (state) {
    state.user = storage.get('current_site')
  },
  [types.DELETE_USER] (state) {
    storage.remove('current_site')
    state.user = {}
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
