import HelloWorld from '@/components/HelloWorld'
import Login from '@/components/Login'
import Editer from '@/components/BlogEditer'
import Page from '@/components/BlogPage'
import Menu from '@/components/BookMenu'
import P401 from '@/components/error/401'
import P404 from '@/components/error/404'
import P500 from '@/components/error/500'

export default[
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
    path: '/editer/:code',
    component: Editer
  },
  {
    path: '/page/:code',
    component: Page
  },
  {
    path: '/menu',
    component: Menu
  },
  {
    path: '/401',
    component: P401
  },
  {
    path: '/404',
    component: P404
  },
  {
    path: '/500',
    component: P500
  },
  {
    path: '*',
    component: P404
  }
]
