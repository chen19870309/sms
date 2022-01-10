import HelloWorld from '@/components/HelloWorld'
import Login from '@/components/Login'
import Editer from '@/components/BlogEditer'
import Page from '@/components/BlogPage'
import Menu from '@/components/BookMenu'

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
  }
]
