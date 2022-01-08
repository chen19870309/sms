<template>
  <div>
    <site-header></site-header>
    <div class="site-main section-inner thin animated fadeIn menus">
      <h1>Menus</h1>
      <book-group v-for="chepter in menu.chepters" :key="chepter.id" :chepter="chepter"></book-group>
    </div>
    <site-footer></site-footer>
  </div>
</template>

<script>
import SiteHeader from '@/components/global/SiteHeader'
import SiteFooter from '@/components/global/SiteFooter'
import BookGroup from '@/components/global/BookGroup'
export default {
  name: 'book-menu',
  components: {
    BookGroup,
    SiteHeader,
    SiteFooter
  },
  data () {
    return {
      menu: {
        chepters: [
          {
            id: 1,
            name: 'chapter1',
            books: [
              {
                id: 1,
                url: '#/page/1',
                title: '我爱踢足球',
                day: 'Sep,7,2021'
              },
              {
                id: 2,
                url: '#2',
                title: '我爱踢足球',
                day: 'Sep,8,2021'
              },
              {
                id: 3,
                url: '#3',
                title: '我爱踢足球',
                day: 'Sep,9,2021'
              },
              {
                id: 4,
                url: '#4',
                title: '我爱踢足球',
                day: 'Sep,10,2021'
              },
              {
                id: 5,
                url: '#5',
                title: '我爱踢足球',
                day: 'Sep,11,2021'
              }
            ]
          },
          {
            id: 2,
            name: 'chapter2',
            books: [
              {
                id: 6,
                url: '#1',
                title: '我爱踢足球',
                day: 'Sep,7,2021'
              },
              {
                id: 7,
                url: '#2',
                title: '我爱踢足球',
                day: 'Sep,8,2021'
              },
              {
                id: 8,
                url: '#3',
                title: '我爱踢足球',
                day: 'Sep,9,2021'
              },
              {
                id: 9,
                url: '#4',
                title: '我爱踢足球',
                day: 'Sep,10,2021'
              },
              {
                id: 10,
                url: '#5',
                title: '我爱踢足球',
                day: 'Sep,11,2021'
              }
            ]
          }
        ]
      }
    }
  }
}
</script>
