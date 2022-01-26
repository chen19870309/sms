import HelloWorld from '@/components/HelloWorld'
import Login from '@/components/Login'
import Regist from '@/components/Regist'
import Editer from '@/components/BlogEditer'
import BlogCache from '@/components/BlogCache'
import BlogIndex from '@/components/layout/BlogIndex'
import Personal from '@/components/layout/Personal2'
import Page from '@/components/BlogPage'
import Menu from '@/components/BookMenu'
import P401 from '@/components/error/401'
import P404 from '@/components/error/404'
import P500 from '@/components/error/500'

import JsonTool from '../views/WebTool'

export default[
  {
    path: '/',
    name: 'index',
    component: BlogIndex
  },
  {
    path: '/user',
    component: Personal
  },
  {
    path: '/hello',
    component: HelloWorld
  },
  {
    path: '/tools',
    component: JsonTool
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/regist',
    name: 'regist',
    component: Regist
  },
  {
    path: '/editer/:code',
    name: 'editer',
    component: Editer
  },
  {
    path: '/page/:code',
    name: 'page',
    component: Page
  },
  {
    path: '/menu',
    name: 'menu',
    component: Menu
  },
  {
    path: '/cache',
    name: 'cache',
    component: BlogCache
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
