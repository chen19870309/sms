import Vue from 'vue'
import iView from 'iview'
import Router from 'vue-router'
import routes from './routers.js'
import store from '@/store'

Vue.use(Router)

const router = new Router({
  routes,
  mode: 'history'
})
router.beforeEach((to, from, next) => {
  iView.LoadingBar.start()
  console.log('beforeEach:' + store.state.currentUser.hasGetInfo)
  // if user not exist, then go to login
  console.log(to.path)
  if (to.name === 'editer' || to.name === 'cache') {
    if (!store.state.currentUser.hasGetInfo) {
      store.commit('setNext', to.path)
      router.replace({name: 'login'})
    }
  }
  next()
})
router.afterEach(to => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
  iView.LoadingBar.finish()
})
export default router
