<template>
  <div>
  <BlogHeader></BlogHeader>
  <div class="post-page thin body">
    <header class="post-header">
      <div class="post-meta"><span>Sep 17, 2021</span></div>
      <span>Go 数组比切片好在哪？--by 煎鱼</span>
    </header>
    <hr>
    <markdown-preview :initialValue="mddata" :theme="mdtheme"></markdown-preview>
  </div>
  <blog-footer></blog-footer>
</div>
</template>

<script>
import { MarkdownPreview } from 'vue-meditor'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/global/SiteFooter'
export default {
  data () {
    return {
      mdtheme: 'dark',
      mddata: '#  Go 数组比切片好在哪？ \n ##  二级标题\n ```\nfunc(){\n test.abs==nil \n} \n'
    }
  },
  components: {
    MarkdownPreview,
    BlogHeader,
    BlogFooter
  }
}
</script>
