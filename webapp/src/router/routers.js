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
    component: BlogIndex,
    meta: {
      title: ' ( ゜- ゜)つロ ～～一起学习吧~'
    }
  },
  {
    path: '/user',
    name: 'user',
    component: Personal,
    meta: {
      title: '个人中心'
    }
  },
  {
    path: '/hello',
    component: HelloWorld
  },
  {
    path: '/tools',
    component: JsonTool,
    meta: {
      title: 'web工具'
    }
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: {
      title: '用户登陆'
    }
  },
  {
    path: '/regist',
    name: 'regist',
    component: Regist,
    meta: {
      title: '账号注册'
    }
  },
  {
    path: '/editer/:code',
    name: 'editer',
    component: Editer,
    meta: {
      title: '写文章'
    }
  },
  {
    path: '/page/:code',
    name: 'page',
    component: Page,
    meta: {
      title: '(๑•̀ㅂ•́)و✧'
    }
  },
  {
    path: '/menu',
    name: 'menu',
    component: Menu,
    meta: {
      title: '导航菜单'
    }
  },
  {
    path: '/cache',
    name: 'cache',
    component: BlogCache,
    meta: {
      title: '草稿箱'
    }
  },
  {
    path: '/401',
    component: P401,
    meta: {
      title: 'ヽ(*。>Д<)o゜'
    }
  },
  {
    path: '/404',
    component: P404,
    meta: {
      title: 'ヽ(*。>Д<)o゜'
    }
  },
  {
    path: '/500',
    component: P500,
    meta: {
      title: 'ヽ(*。>Д<)o゜'
    }
  },
  {
    path: '*',
    component: P404,
    meta: {
      title: 'ヽ( ￣д￣;)ノ∑(っ'
    }
  }
]
