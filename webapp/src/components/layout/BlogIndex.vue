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
        <router-link to="/tools"> <Icon type="logo-nodejs" /> &nbsp;JSON在线工具 </router-link>
    </div>
    <div class="container-fluid"></div>
</div>

<div class="container-fluid content-box" id="article">
    <div class="container content">
        <div class="col-lg-12 col-md-12 col-sm-12" id="art">
            <doc-tab v-for="doc in docs" :key="doc.Id" :doc=doc></doc-tab>
        </div>
        <div class="view-more">
            <Button type="success" v-show="total > more"  :loading="loading"  @click="moreindex" >
              <span v-if="!loading">加载更多...</span>
              <span v-else>Loading...</span>
            </Button>
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
      docs: [],
      more: 8,
      total: 0,
      loading: false
    }
  },
  methods: {
    query (data) {
      this.docs = data
    },
    moreindex () {
      this.more = this.more + 4
      this.loading = true
      NetWorking.doGet(API.mainindex + '?more=' + this.more).then(response => {
        console.log(response.data)
        this.docs = response.data
        this.total = response.count
        this.loading = false
      },(message) => {
        this.loading = false
        this.$Message.error('Load More Failed!' + message)
        this.$router.push({ path: '/404' })
      })
    }
  },
  created() {
    NetWorking.doGet(API.mainindex).then(response => {
      console.log(response.data)
      this.docs = response.data
      this.total = response.count
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
