<template>
    <header id="blog-header">
        <div class="hdr-wrapper section-inner">
            <div class="hdr-left">
                <div class="site-branding">
                    <a href="/">Hola-CosMos</a>
                </div>
            </div>
            <div class="hdr-right">
                <Dropdown>
                    <a href="javascript:void(0)">
                        写文章
                    </a>
                    <DropdownMenu slot="list">
                        <DropdownItem><router-link to="/menu">文章目录</router-link></DropdownItem>
                        <DropdownItem><a href="#" @click.prevent="newblog" >新建文章</a></DropdownItem>
                        <DropdownItem ><router-link to="/cache">草稿箱</router-link></DropdownItem>
                    </DropdownMenu>
                </Dropdown>
                <Divider type="vertical" />
                <Dropdown>
                    <a href="javascript:void(0)">
                        账号
                        <Icon type="ios-arrow-down"></Icon>
                    </a>
                    <DropdownMenu slot="list">
                        <DropdownItem v-show='user.Id != undefined' >
                          <router-link to="/user">
                            <user-center></user-center>
                          </router-link></DropdownItem>
                        <DropdownItem><router-link to="/login">账号登陆</router-link></DropdownItem>
                        <DropdownItem><router-link to="/regist">注册账号</router-link></DropdownItem>
                        <DropdownItem v-show='user.Id != undefined' divided><a href="#" @click.prevent="logout">退出账号</a></DropdownItem>
                    </DropdownMenu>
                </Dropdown>
            </div>
        </div>
    </header>
</template>

<script>
import { mapGetters } from 'vuex'
import UserCenter from '@/components/layout/UserCenter'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
export default {
  name: 'site-header',
  methods: {
    newblog () {
      NetWorking.doGet(API.newblog).then(response => {
        let data = response.data
        this.$store.dispatch('createBlog', data)
        this.$router.push({ path: '/editer/' + data.Code }).catch(err => {
            this.$Message.error('Auto New MarkDown Failed!' + err)
        })
      }, (message) => {
        this.$Message.error('Auto New MarkDown Failed!' + message)
      })
    },
    logout () {
        this.$store.dispatch('deleteUser')
        this.$store.dispatch('deleteBlog')
        this.$router.push({ path: '/menu' })
    },
    showUser () {
      this.$Modal.info({
        title: 'User Info',
        content: `Name：${this.user.Nickname}<br>Secret: ${this.user.Secret}<br>CreateTime：${this.user.Createtime}`
      })
    }
  },
  computed: {
    ...mapGetters({
      user: 'currentUser'
    })
  },
  components: {
      UserCenter
  }
}
</script>
