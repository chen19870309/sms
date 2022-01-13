<template>
  <div class="article">
    <Markdown :height="800" :autoSave=true @on-save="handleOnSave" :theme="theme" :value="blog.Content" :interval="60000"></Markdown>
    <div id="mobile-menu" class="animated fast">
      <ul>
        <li><a href="#" @click.prevent="newblog" >新建</a></li>
        <li><a href="#" @click.prevent="newpush" >发布</a></li>
        <li><a href="#" @click.prevent="goback">返回</a></li>
      </ul>
    </div>
  </div>
</template>

<script>
import Cookie from 'js-cookie'
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
    goback () {
      this.$router.go(-1)
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
        this.$store.dispatch('deleteUser')
      })
    },
    handleOnSave ({value, theme}) {
      this.$store.dispatch('updateBlog', value)
      this.theme = theme
      let params = {
        data: value,
        theme: theme,
        author_id: this.user.Id
      }
      NetWorking.doPost(API.save + this.blog.Code, null, params).then(response => {
      }, (message) => {
        this.$Notice.error({
          title: '自动保存失败',
          desc: 'Auto Save MarkDown Failed!' + message
        })
        this.$store.dispatch('deleteUser')
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
      let data = response.data
      this.$store.dispatch('createBlog', data)
      this.$router.push({ path: '/editer/' + data.Code }).catch(err => {})
    }, (message) => {
      this.$Message.error('Load MarkDown Failed!' + message)
    })
  },
  computed: {
    ...mapGetters({
      blog: 'currentBlog',
      user: 'currentUser'
    })
  }
}
</script>
