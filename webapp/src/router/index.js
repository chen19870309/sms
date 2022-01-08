import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Login from '@/components/Login'
import Editer from '@/components/BlogEditer'
import Page from '@/components/BlogPage'
import Menu from '@/components/BookMenu'
Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/editer',
      component: Editer
    },
    {
      path: '/page/:id',
      component: Page
    },
    {
      path: '/menu',
      component: Menu
    }
  ]
})
