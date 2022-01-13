import * as types from '../mutation-types'
import storage from '@/utils/storage'
const state = {
  bloginfo: {}
}

// getters
const getters = {
  currentBlog: state => state.bloginfo
}
// actions
const actions = {
  updateBlog ({commit}, data) {
    commit('updateBlogCtx', data)
  },
  createBlog ({commit}, blog) {
    commit(types.CREATE_BLOG, blog)
  },
  fetchBlog ({commit}) {
    commit(types.FETCH_BLOG)
  },
  deleteBlog ({commit}) {
    commit(types.DELETE_BLOG)
  }
}
// mutations
const mutations = {
  [types.CREATE_BLOG] (state, blog) {
    state.bloginfo = blog
    storage.set('current_blog', blog)
  },
  [types.FETCH_BLOG] (state) {
    state.bloginfo = storage.get('current_blog')
  },
  [types.DELETE_BLOG] (state) {
    storage.remove('current_blog')
    state.bloginfo = {}
  },
  updateBlogCtx (state, value) {
    state.bloginfo.Content = value
    storage.set('current_blog', state.bloginfo)
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
