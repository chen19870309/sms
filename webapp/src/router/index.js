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
  console.log('beforeEach')
  // if user not exist, then go to login
  console.log(to)
  if (to.name === 'editer') {
    let user = store.getters.currentUser
    console.log(user)
    if (user === undefined || user.Id === undefined) {
      // router.replace({name: 'login'})
    }
  }
  next()
})
router.afterEach(to => {
  iView.LoadingBar.finish()
})
export default router
