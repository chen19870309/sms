<template>
  <div class="article">
    <MarkdownPro :height="800" :autoSave=true @on-save="handleOnSave" :theme="theme"></MarkdownPro>
    <div id="mobile-menu" class="animated fast">
      <ul>
        <li><a href="https://eddycjy.com/posts/">新建</a></li>
        <li><a href="https://eddycjy.com/tags/">发布</a></li>
        <li><a href="https://eddycjy.com/about/">返回</a></li>
      </ul>
    </div>
  </div>
</template>

<script>
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import MarkdownPro from 'vue-meditor'
export default {
  name: 'markdown',
  components: {
    MarkdownPro
  },
  methods: {
    handleOnSave ({value, theme}) {
      console.log(value, theme)
      this.data = value
      this.theme = theme
      let params = {
        markdown: value,
        theme: theme
      }
      NetWorking.doPost(API.save+this.code,null,params).then(response => {
            this.disabled = false
          }, (message) => {
            this.disabled = false
            this.$Message.error('Auto Save MarkDown Failed!' + message)
          })
    }
  },
  data () {
    return {
      data: '',
      theme: 'oneDark',
      code: this.$route.params.code
    }
  }
}
</script>
