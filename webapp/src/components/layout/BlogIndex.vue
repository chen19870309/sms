<template>
<div>
<nav id="mHeader" class="navbar navbar-inverse navbar-fixed-top">
  <site-header></site-header>
</nav>

<!--banner部分-->
<search-bar :callback="query"></search-bar>
<!--banner部分-->

<div class="wrap container-fluid">
    <div class="container wrap-cont">
        <router-link to="/menu"><Icon type="ios-menu" /> &nbsp;导航 </router-link> <Divider type="vertical" />
        <router-link to="#"> <Icon type="ios-apps" /> &nbsp;其他  </router-link> <Divider type="vertical" />
        <a href="https://www.bejson.com/"><Icon type="logo-nodejs" /> &nbsp;JSON在线工具</a>
    </div>
    <div class="container-fluid"></div>
</div>

<div class="container-fluid content-box" id="article">
    <div class="container content">
        <div class="col-lg-12 col-md-12 col-sm-12" id="art">
            <doc-tab v-for="doc in docs" :key="doc.Id" :doc=doc></doc-tab>
        </div>
        <div class="view-more">
            <button type="button" id="view-more" class="btn btn-primary center-block">加载更多</button>
        </div>
    </div>
</div>

<blog-footer></blog-footer>
</div>
</template>

<script>
import SiteHeader from '@/components/global/SiteHeader'
import SearchBar from '@/components/layout/SearchBar'
import BlogFooter from '@/components/layout/BlogFooter'
import DocTab from '@/components/layout/DocTab'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
export default {
  data() {
    return {
      login: false,
      docs: []
    }
  },
  methods: {
    query (data) {
      this.docs = data
    }
  },
  created() {
    NetWorking.doGet(API.mainindex).then(response => {
      console.log(response.data)
      this.docs = response.data
    }, (message) => {
      this.$Message.error('Load Main Failed!' + message)
      this.$router.push({ path: '/404' })
    })
  },
  components: {
    SiteHeader,
    SearchBar,
    BlogFooter,
    DocTab
  }
}
</script>

<style scoped src="../../assets/css/personal2.css"></style>
<style scoped src="../../assets/css/bootstrap.min.css"></style>
