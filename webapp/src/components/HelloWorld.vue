<template>
  <div class="hello">
    <Layout >
      <site-header class = "header-fixed-top"></site-header>
      <div >
        <Layout class="body-fixed">
          <Sider :hidden="false"></Sider>
            <Content>
              <router-view/>
            </Content>
        </Layout>
      </div>
    </Layout>
  </div>
</template>

<script>
import SiteFooter from '@/components/global/SiteFooter'
import SiteHeader from '@/components/global/SiteHeader'
export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  },
  components: {
    SiteHeader,
    SiteFooter
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
