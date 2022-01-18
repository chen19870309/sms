
<template>
  <footer class="section-inner animated fadeIn">
    <hr style="height:1px;border:none;border-top:1px dashed #5566CC;">
                <p>&copy; 2018 - 2021 <a :href="site.url">{{ site.name }}</a> &#183; <a href="http://www.beian.miit.gov.cn/">{{ site.beian }}</a></p>
                <p>
                        Made with <a href="https://gohugo.io/" target="_blank" rel="noopener">Hugo</a> &#183; Theme <a href="https://github.com/Track3/hermit" target="_blank" rel="noopener">Hermit</a> &#183; <a href="https://eddycjy.com/posts/index.xml" target="_blank" title="rss"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-rss"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg></a>
                </p>
  </footer>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters({
      site: 'currentSite'
    })
  }
}
</script>
