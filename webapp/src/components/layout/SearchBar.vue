<template>
  <div class="container-fluid banner bank">
    <h1 class="text-center h1">Look For Less,Do More</h1>
    <p class="text-center">provides the latest , the greatest and the most comprehensive jQuery plugins </p>
    <div class="search-box center-block">
        <input type="text" v-model="searchVal" class="search center-block" placeholder="搜索文章..." @keyup.enter="doSearch" :disabled="searching"/>
        <i class="fa fa-search"></i>
        <Spin v-if="searching">
          <Icon type="ios-loading" size=18 class="spin-icon-load"></Icon>
        </Spin>
    </div>
  </div>
</template>

<script>
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
export default {
  props: {
      callback: {
        type: Function,
        required: true
      }
  },
  data() {
    return {
      searching: false,
      searchVal: ''
    }
  },
  methods: {
    doSearch () {
      console.log("search :" + this.searchVal)
      this.searching = !this.searching
      let data = {
          text: this.searchVal
      }
      NetWorking.doPost(API.search,null,data).then(response => {
        this.$props.callback(response.data)
        this.done()
      },(message) => {
          this.$Message.error('Search Doc Failed!' + message)
          this.done()
      })
    },
    done () {
      this.searching = false
    }
  }
}
</script>

<style>
.spin-icon-load{
  animation: ani-demo-spin 1s linear infinite;
}
@keyframes ani-demo-spin {
  from { transform: rotate(0deg);}
  50%  { transform: rotate(180deg);}
  to   { transform: rotate(360deg);}
}
</style>


<style scoped src="../../assets/css/personal2.css"></style>
<style scoped src="../../assets/css/bootstrap.min.css"></style>

