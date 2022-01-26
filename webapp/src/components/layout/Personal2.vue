<template>
<div>
<nav id="mHeader" class="navbar navbar-inverse navbar-fixed-top">
  <site-header></site-header>
</nav>
<div class="container user">
    <div class="user-cont clearfix">
        <div class="col-md-4 user-left">
            <div class="user-left-n clearfix">
                <h6>{{ user.Username }}</h6>
                <a href="#" class="user-headimg f" @click.prevent="UploadModel = true"><img :src="''+user.Icon" v-if="headIcon"></a>
                <div class="user-name f">
                    <h4>{{ user.Nickname }}</h4>
                    <p>{{ user.Remark }}</p>
                </div>
            </div>
            <div class="user-left-n clearfix">
                <div class="alert alert-warning">
                    <i class="fa fa-lightbulb-o"></i>&nbsp;<strong>明日复明日</strong><br>
                    <p>很多人接受现实的荒野</p>
                    <p>不知不觉已垂暮之年</p>
                    <p>一年十年就在一辈子之间</p>
                    <p>忘记了要去寻找你的世界</p>
                    <p>生活不止眼前的苟且</p>
                    <p>还有诗和远方的田野</p>
                </div>
            </div>
            <div class="user-left-n clearfix">
                <a href="#" class="btn btn-success infos"  @click.prevent="newblog" ><Icon type="md-paper" />&nbsp;开始写作</a>
                <a href="#" class="btn btn-warning" @click.prevent="Model = true"><Icon type="md-alert" />&nbsp;修改密码</a>
            </div>
        </div>
        <div class="col-md-8 user-right">
            <div class="user-right-n clearfix">
                <Tabs size="small">
                    <TabPane id="user-article" label="我的文章">
                        <List header="文章列表" border size="small">
                          <ListItem v-for="mk in markds" :key="mk.id">
                              <Icon type="md-eye" v-show="mk.status == 1" />
                              <Icon type="ios-person" v-show="mk.status == 2" />
                              <a href="#" @click.prevent="editblog(mk.code)">#[{{ mk.id }}]|{{mk.updatetime}}|.[{{ mk.title }}]</a></ListItem>
                        </List>
                        <Page id="user-article-page" :total="page.total" :current="page.index" @on-change="indexpage" />
                    </TabPane>
                    <TabPane label="修改信息">
                        <Form :model="formLeft" label-position="left" :label-width="100">
                            <FormItem label="昵称">
                                <Input v-model="formLeft.nickname" />
                            </FormItem>
                            <FormItem label="备注">
                                <Input v-model="formLeft.remark" type="textarea" :autosize="{minRows: 2,maxRows: 10}" placeholder="Enter something..." />
                            </FormItem>
                            <FormItem>
                                <Button type="primary"  @click="handleSubmit('formLeft')" >更新</Button>
                            </FormItem>
                        </Form>
                    </TabPane>
                </Tabs>
            </div>
        </div>
    </div>
</div>
<blog-footer></blog-footer>
<Modal
    :title=Title
    v-model=Model
    @on-ok="ok"
    @on-cancel="cancel"
    :mask-closable="false">
    <p>
    <i-input type="password" v-model="oldpwd" placeholder="原密码">
      <Icon type="ios-lock-outline" slot="prepend"></Icon>
    </i-input>
    </p>
    <p>
    <i-input type="password" v-model="newpwd" placeholder="新密码">
      <Icon type="ios-lock-outline" slot="prepend"></Icon>
    </i-input>
    </p>
</Modal>
<Modal
    title='上传头像'
    v-model=UploadModel
    @on-cancel="cancel"
    :mask-closable="false">
<upload :uptoken='qiniuuptoken'
        :filename='filename'
        browse_button='pickfile'
        domain='http://r5uiv7l5f.hd-bkt.clouddn.com'
        bucket_name='sp2022'
        @on-percent='filePercent'
        @on-change='uploaded'>
        <Button type="ghost" id='pickfile' slot='button'>选择文件</Button>
        <p></p>
        <Progress slot='progressBar' :percent="up_percent"></Progress>
</upload>
</Modal>
</div>
</template>

<script>
import SiteHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/layout/BlogFooter'
import Upload from 'qiniu-upload-vue'
import { mapGetters } from 'vuex'
import Cookie from 'js-cookie'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
export default {
  name: 'personal-page',
  data() {
    return {
      formLeft: {
        nickname: '',
        remark: ''
      },
      Title: '修改密码',
      headIcon: true,
      Model: false,
      UploadModel: false,
      oldpwd: '',
      newpwd: '',
      markds: [],
      qiniuuptoken: '',
      up_percent: 0,
      filename: 'test',
      page: {
        index:1,
        size: 10,
        total:0
      }
    }
  },
  created() {
      this.formLeft.nickname = this.user.Nickname
      this.formLeft.remark = this.user.Remark
      NetWorking.doGet(API.userindex).then(response => {
        this.markds = response.data
        this.page.total = response.count
        this.page.index = Cookie.get("user_page_index")
      },(message) => {
          this.$Message.error('Get User MarkDown Files Failed!' + message)
      })
      NetWorking.doGet(API.uptoken).then(response => {
        console.log("UpToken:",response.data)
        this.qiniuuptoken = response.data.token
      },(message) => {
          this.$Message.error('Get User MarkDown Files Failed!' + message)
      })
  },
  methods: {
    handleSubmit (info) {
      NetWorking.doPost(API.edituser,null,this.formLeft).then(response => {
        this.$Notice.success({
          title: 'Update Info Success!!'
        })
      },(message) => {
        this.$Notice.error({
          title: 'Update Info Failed!!',
          desc: message
        })
      })
    },
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
    editblog (code) {
      this.$router.push({ path: '/editer/' + code }).catch(err => {
        console.log('edit blog:' + code)
      })
    },
    indexpage (i) {
      console.log("indexpage:"+i)
      Cookie.set('user_page_index', i)
      NetWorking.doGet(API.userindex).then(response => {
          this.markds = response.data
          this.page.total = response.count
          console.log(this.data,this.page)
      },(message) => {
          this.$Message.error('Get User MarkDown Files Failed!' + message)
      })
    },
    ok () {
      if(this.oldpwd === this.newpwd) {
        this.$Notice.error({title: '新老密码不能相同!'});
      } else {
        if(this.oldpwd === '' || this.newpwd === ''){
          this.$Notice.error({title: '新老密码不能为空!'});
        } else {
            let data = {
              OP: this.$md5(this.oldpwd),
              NP: this.$md5(this.newpwd)
            }
            console.log(data)
            NetWorking.doPost(API.editpwd,null,data).then(response => {
              this.$Message.info('update password ok!')
            }, (message) => {
              this.$Notice.error({
              title: 'Update Password Failed!!',
              desc: message
            })
          })
        }
      }
    },
    cancel () {
      this.$Message.info('Clicked cancel')
    },
    filePercent (val) {
      console.log('filePercent',val)
    },
    uploaded (val) {
      this.up_percent = 100
      this.headIcon = false
      let nUser = this.user
      nUser.Icon = val
      this.$store.dispatch('createUser', nUser)
      setTimeout(()=>{
        this.UploadModel = false
        this.headIcon = true
        this.up_percent = 0
      },1000)
      let data = {
        nickname: this.user.Nickname,
        icon: val
      }
      NetWorking.doPost(API.edituser,null,data).then(response => {
        this.$Notice.success({
          title: '更头像成功!!'
        })
      },(message) => {
        this.$Notice.error({
          title: '更新头像失败!!',
          desc: message
        })
      })
      console.log('Uploaded:',val)
    }
  },
  computed: {
    ...mapGetters({
      user: 'currentUser'
    })
  },
  components: {
    SiteHeader,
    BlogFooter,
    Upload
  }
}
</script>

<style scoped src="../../assets/css/personal2.css"></style>
<style scoped src="../../assets/css/bootstrap.min.css"></style>
