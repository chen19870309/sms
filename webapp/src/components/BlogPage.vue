<template>
  <div>
  <BlogHeader class="header-fixed-top"></BlogHeader>
  <div class="post-page thin body">
    <markdown-preview :initialValue="blog.Content" :theme="mdtheme"></markdown-preview>
  </div>
   <div id="mobile-menu" class="animated fast" v-show='user.Id != undefined'>
      <ul>
        <li><a href="#" @click.prevent="newblog" >新建</a></li>
        <li><a href="#" @click.prevent="editblog" >编辑</a></li>
        <li><a href="#" @click.prevent="gomenu">返回</a></li>
      </ul>
    </div>
  <blog-footer></blog-footer>
</div>
</template>

<script>
import { mapGetters } from 'vuex'
import { MarkdownPreview } from 'vue-meditor'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/layout/BlogFooter'
export default {
  data () {
    return {
      mdtheme: 'oneDark'
    }
  },
  created () {
    NetWorking.doGet(API.posts + this.$route.params.code).then(response => {
      console.log(response.data)
      let data = response.data
      this.$store.dispatch('createBlog', data)
    }, (message) => {
      this.$Message.error('Load  MarkDown Failed!' + message)
      this.$router.push({ path: '/404' })
    })
  },
  methods: {
    gomenu () {
      this.$router.push({path: '/menu'})
    },
    newblog () {
      NetWorking.doGet(API.newblog).then(response => {
        let data = response.data
        this.$store.dispatch('createBlog', data)
        this.$router.push({path: '/editer/' + data.Code}).catch( err => {
          console.log('pass router')
        })
      }, (message) => {
        this.$Notice.error({
          title: '新建文章失败',
          desc: 'Auto New MarkDown Failed!' + message
        })
        this.$store.dispatch('deleteUser')
      })
    },
    editblog () {
      this.$router.push({ path: '/editer/' + this.blog.Code }).catch(err => {
        console.log('edit blog:' + this.blog.Code)
      })
    }
  },
  computed: {
    ...mapGetters({
      blog: 'currentBlog',
      user: 'currentUser'
    })
  },
  components: {
    MarkdownPreview,
    BlogHeader,
    BlogFooter
  }
}
</script>
