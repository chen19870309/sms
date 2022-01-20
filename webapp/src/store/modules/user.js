import * as types from '../mutation-types'
import storage from '@/utils/storage'
const state = {
  userinfo: {},
  hasGetInfo: false,
  token: '',
  nextUrl: '/user'
}

// getters
const getters = {
  currentUser: state => state.userinfo,
  nextUrl: state => state.nextUrl,
  authToken: state => state.Token
}
// actions
const actions = {
  setNext ({commit}, url) {
    commit('setNext', url)
  },
  saveToken ({commit}, token) {
    commit('setToken', token)
  },
  createUser ({commit}, user) {
    console.log('in create user')
    commit('setHasGetInfo', true)
    commit('setToken', user.Secret)
    commit(types.CREATE_USER, user)
  },
  fetchUser ({commit}) {
    commit(types.FETCH_USER)
  },
  deleteUser ({commit}) {
    commit(types.DELETE_USER)
    commit('setHasGetInfo', false)
    commit('setToken', '')
  }
}
// mutations
const mutations = {
  setNext (state, url) {
    state.nextUrl = url
  },
  setHasGetInfo (state, status) {
    state.hasGetInfo = status
  },
  setToken (state, token) {
    state.token = token
    storage.set('auth_token', token)
  },
  [types.CREATE_USER] (state, user) {
    state.userinfo = user
    state.setHasGetInfo = true
    state.token = user.Code
    storage.set('current_user', user)
  },
  [types.FETCH_USER] (state) {
    state.userinfo = storage.get('current_user')
  },
  [types.DELETE_USER] (state) {
    storage.remove('current_user')
    state.userinfo = {}
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
