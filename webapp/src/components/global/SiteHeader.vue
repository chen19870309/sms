<template>
    <header id="blog-header">
        <div class="hdr-wrapper section-inner">
            <div class="hdr-left">
                <div class="site-branding">
                    <a href="#">Hola-CosMos</a>
                </div>
            </div>
            <div class="hdr-right">
                <Dropdown>
                    <a href="javascript:void(0)">
                        下拉菜单
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
                        theme
                        <Icon type="ios-arrow-down"></Icon>
                    </a>
                    <DropdownMenu slot="list">
                        <DropdownItem><router-link to="/login">账号登陆</router-link></DropdownItem>
                        <DropdownItem>炸酱面</DropdownItem>
                        <DropdownItem disabled>豆汁儿</DropdownItem>
                        <DropdownItem>冰糖葫芦</DropdownItem>
                        <DropdownItem divided>退出账号</DropdownItem>
                    </DropdownMenu>
                </Dropdown>
            </div>
        </div>
    </header>
</template>

<script>
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
    }
  }
}
</script>
