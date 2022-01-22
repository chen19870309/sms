<template>
  <div class="article">
    <Markdown :height="800" :autoSave=true @on-save="handleOnSave"  :theme="theme" :value="blog.Content" :interval="30000" @on-ready="prepareMd" ></Markdown>
    <div id="mobile-menu" class="animated fast">
      <ul>
        <li><a href="#" @click.prevent="UploadModel = true">图片</a></li>
        <li><a href="#" @click.prevent="newblog" >新建</a></li>
        <li><a href="#" @click.prevent="newpush" >发布</a></li>
        <li><a href="#" @click.prevent="goback">返回</a></li>
      </ul>
    </div>
  <Modal
    title='上传图片'
    v-model=UploadModel
    @on-cancel="cancel"
    :mask-closable="false">
<upload :uptoken='qiniu.token'
        :filename='qiniu.prefix'
        browse_button='pickfile'
        :domain='qiniu.domain'
        :bucket_name='qiniu.backet'
        @on-percent='filePercent'
        @on-change='uploaded'>
        <Button type="primary" id='pickfile' slot='button'>选择文件</Button>
        <p></p>
        <Progress slot='progressBar' :percent="up_percent"></Progress>
</upload>
</Modal>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import Markdown from 'vue-meditor'
import Upload from 'qiniu-upload-vue'
export default {
  name: 'markdown',
  components: {
    Markdown,
    Upload
  },
  methods: {
    goback () {
      this.$router.go(-1)
    },
    filePercent (val) {
      console.log(val)
    },
    uploaded (val) {
      this.up_percent = 100
      this.markdown.insertContent('![image]('+val+')\n')
      setTimeout(() => {
        this.UploadModel = false
        this.up_percent = 0
      }, 1000)
    },
    prepareMd (md) {
      this.markdown = md
    },
    newpush () {
      let params = {
        data: this.blog.Content,
        theme: this.theme,
        author_id: this.user.Id
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
    },
  },
  data () {
    return {
      theme: 'oneDark',
      up_percent: 0,
      UploadModel: false,
      qiniu: {},
      markdown: {},
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
    NetWorking.doGet(API.uptoken).then(response => {
        console.log('UpToken:', response.data)
        this.qiniu = response.data
      },(message) => {
          this.$Message.error('Get User MarkDown Files Failed!' + message)
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
