<template>
    <div class="box">
    <!-- 粒子背景 -->
    <vue-particles></vue-particles>
    <site-header></site-header>
    <Card class="card-box">
    <Form ref="formRegist" :model="formRegist" :rules="formRegistRules" >
    <Form-item class="formLogin-title">
        <h3>注册账号</h3>
    </Form-item>
        <Form-item prop="username">
            <i-input size="large" type="text" v-model="formRegist.username" placeholder="用户名">
                <Icon type="ios-person-outline" slot="prepend"></Icon>
            </i-input>
        </Form-item>
        <Form-item prop="email">
            <i-input size="large" type="text" v-model="formRegist.email" placeholder="用户邮箱">
                <Icon type="ios-mail" slot="prepend"></Icon>
            </i-input>
        </Form-item>
        <Form-item prop="code">
          <Row>
            <i-col span="16">
              <i-input size="large" type="text" v-model="formRegist.code" placeholder="邮箱校验码">
                  <Icon type="ios-finger-print" slot="prepend"></Icon>
              </i-input>
            </i-col>
            <i-col span="6" offset="1">
                <i-button type="primary" @click="checkEmail('formRegist')" :disabled="disemail">校验邮箱</i-button>
            </i-col>
          </Row>
        </Form-item>
        <Form-item prop="password">
            <i-input size="large"  type="password" v-model="formRegist.password" placeholder="密码">
                <Icon type="ios-lock-outline" slot="prepend"></Icon>
            </i-input>
        </Form-item>
        <Form-item class="login-no-bottom">
            <Row >
                <i-col :xs="{ span: 4, offset: 6 }" >
                    <i-button type="primary" @click="handleSubmit('formRegist')" :disabled="disabled">注册</i-button>
                </i-col>
                <i-col :xs="{ span: 4, offset: 6 }">
                    <i-button  type="primary" @click="formRegistReset('formRegist')">重置</i-button>
                </i-col>
            </Row>
        </Form-item>
    </Form>
    </Card>
    </div>
</template>

<script>
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import Cookie from 'js-cookie'
import SiteHeader from '@/components/global/SiteHeader'
export default {
  data () {
    return {
      formRegist: {
        username: '',
        email: '',
        code: '',
        password: '',
      },
      formRegistRules: {
        username: [
          { required: true, message: '请填写用户名', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请填写用户邮箱', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请填写校验码', trigger: 'blur' },
          { type: 'string', min: 6, message: '校验码长度不小于6位', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请填写密码', trigger: 'blur' },
          { type: 'string', min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
        ]
      },
      disabled: false,
      disemail: false,
    }
  },
  methods: {
    checkEmail (name) {
      this.disemail = true
      NetWorking.doPost(API.checkemail, null, this.formRegist).then(response => {
        this.disemail = false
        this.$Notice.info({
          title: '激活码已经发生邮箱请注意查看!',
        })
      }, (message) => {
        this.disemail = false
        this.$Notice.error({
          title: '邮箱校验失败!',
          desc: message
        })
      })
    },
    handleSubmit (name) {
      this.$refs[name].validate((valid) => {
        sessionStorage.setItem('username', JSON.stringify(this.formRegist.username))
        if (valid) {
          // this.$Message.success('提交成功!')
          // this.$router.push({ path: '/menu' })
          console.log(this.formRegist)
          let data = {
            username: this.formRegist.username,
            email: this.formRegist.email,
            code: this.formRegist.code,
            password: this.$md5(this.formRegist.password)
          }
          this.disabled = true
          NetWorking.doPost(API.regist, null, data).then(response => {
            this.disabled = false
            let user = response.data
            Cookie.set('auth_token', user.Secret,{ expires: 1})
            Cookie.set('jwt',response.jwt,{ expires: 1})
            this.$store.dispatch('createUser', user,{ expires: 1})
            this.$router.push({ path: '/login' }).catch(err => {})
          }, (message) => {
            this.disabled = false
            // this.$Message.error('Login Failed!' + message)
            this.$Notice.error({
              title: 'Regist Failed!',
              desc: message
            })
          })
        } else {
          this.$Message.error('表单验证失败!')
        }
      })
    },
    formRegistReset (name) {
      this.disabled = false
      this.$refs[name].resetFields()
    }
  },
  components: {
    SiteHeader
  }
}
</script>
