<template>
  <div class="article">
    <MarkdownPro :height="800"></MarkdownPro>
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
import MarkdownPro from 'vue-meditor'
export default {
  name: 'markdown',
  components: {
    MarkdownPro
  }
}
</script>
