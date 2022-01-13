<template>
  <div>
  <BlogHeader></BlogHeader>
  <div class="post-page thin body">
    <header class="post-header">
      <div class="post-meta"><span>{{ blog.CreateTime }}</span></div>
      <span>Go 数组比切片好在哪？[{{ blog.Title }}]--by 煎鱼[{{ blog.Code }}]</span>
    </header>
    <hr>
    <markdown-preview :initialValue="blog.Content" :theme="mdtheme"></markdown-preview>
  </div>
  <blog-footer></blog-footer>
   <div id="mobile-menu" class="animated fast">
      <ul>
        <li><a href="#" @click.prevent="newblog" >新建</a></li>
        <li><a href="#" @click.prevent="editblog" >编辑</a></li>
        <li><a href="#" @click.prevent="gomenu">返回</a></li>
      </ul>
    </div>
</div>
</template>

<script>
import { mapGetters } from 'vuex'
import { MarkdownPreview } from 'vue-meditor'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/global/SiteFooter'
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
        this.$router.push({ path: '/editer/' + data.Code }).catch(err => {
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
        console.log("edit blog:" + this.blog.Code)
      })
    }
  },
  computed: {
    ...mapGetters({
      blog: 'currentBlog'
    })
  },
  components: {
    MarkdownPreview,
    BlogHeader,
    BlogFooter
  }
}
</script>
