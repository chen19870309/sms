<template>
  <div>
    <site-header></site-header>
    <div class="site-main section-inner thin animated fadeIn menus">
      <h1>Menus</h1>
      <book-group v-for="chepter in menu.chepters" :key="chepter.id" :chepter="chepter"></book-group>
    </div>
    <blog-footer></blog-footer>
  </div>
</template>

<script>
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import SiteHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/layout/BlogFooter'
import BookGroup from '@/components/global/BookGroup'
export default {
  name: 'book-menu',
  components: {
    BookGroup,
    SiteHeader,
    BlogFooter
  },
  data () {
    return {
      menu: {}
    }
  },
  created () {
    NetWorking.doGet(API.menu).then(response => {
      console.log(response.data)
      this.menu = response.data
    }, (message) => {
      this.$Message.error('Load Menu Failed!' + message)
    })
  }
}
</script>
