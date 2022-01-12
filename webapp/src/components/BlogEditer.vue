<template>
  <div class="article">
    <Markdown :height="800" :autoSave=true @on-save="handleOnSave" :theme="theme" :value="blog.Content" :interval="60000"></Markdown>
    <div id="mobile-menu" class="animated fast">
      <ul>
        <li><a href="#" @click.prevent="newblog" >新建</a></li>
        <li><a href="#" @click.prevent="newpush" >发布</a></li>
        <li><a href="#" @click.prevent="gohome">返回</a></li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import Markdown from 'vue-meditor'
export default {
  name: 'markdown',
  components: {
    Markdown
  },
  methods: {
    gohome () {
      this.$router.push({path: '/'})
    },
    newpush () {
      let params = {
        data: this.blog.Content,
        theme: this.theme,
        author_id: 1
      }
      NetWorking.doPut(API.posts + this.blog.Code,null,params).then(response => {
        this.$router.push({ path: '/page/' + this.blog.Code })
      }, (message) => {
        this.$Message.error('Put MarkDown Failed!' + message)
      })
    },
    newblog () {
      NetWorking.doGet(API.newblog).then(response => {
        let data = response.data
        this.$store.dispatch('createBlog', data)
        this.$router.push({ path: '/editer/' + data.Code })
      }, (message) => {
        this.$Message.error('Auto New MarkDown Failed!' + message)
      })
    },
    handleOnSave ({value, theme}) {
      console.log(value, theme)
      this.blog.Content = value
      this.$store.dispatch('createBlog', this.blog)
      this.theme = theme
      let params = {
        data: value,
        theme: theme,
        author_id: 1
      }
      NetWorking.doPost(API.save + this.blog.Code, null, params).then(response => {
        this.disabled = false
      }, (message) => {
        this.disabled = false
        this.$Notice.error({
          title: '自动保存失败',
          desc: 'Auto Save MarkDown Failed!' + message
        })
      })
    }
  },
  data () {
    return {
      theme: 'oneDark'
    }
  },
  created () {
    NetWorking.doGet(API.posts + this.$route.params.code).then(response => {
      console.log(response.data)
      let data = response.data
      this.$store.dispatch('createBlog', data)
      this.$router.push({ path: '/editer/' + data.Code })
    }, (message) => {
      this.$Message.error('Load MarkDown Failed!' + message)
    })
  },
  computed: {
    ...mapGetters({
      blog: 'currentBlog'
    })
  }
}
</script>
