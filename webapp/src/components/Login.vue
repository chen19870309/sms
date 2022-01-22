<template>
    <div class="box" >
    <!-- 粒子背景 -->
    <vue-particles></vue-particles>
    <site-header></site-header>
    <Card class="card-box">
    <Form ref="formLogin" :model="formLogin" :rules="formLoginRules" >
    <Form-item class="formLogin-title">
        <h3>系统登录</h3>
    </Form-item>

        <Form-item prop="username">
            <i-input size="large" type="text" v-model="formLogin.username" placeholder="用户名">
                <Icon type="ios-person-outline" slot="prepend"></Icon>
            </i-input>
        </Form-item>
        <Form-item prop="password">
            <i-input size="large"  type="password" v-model="formLogin.password" placeholder="密码">
                <Icon type="ios-lock-outline" slot="prepend"></Icon>
            </i-input>
        </Form-item>
          <Form-item class="login-no-bottom">
            <Checkbox-group v-model="formLogin.remember">
                <Checkbox label="记住密码" name="remember"></Checkbox>
            </Checkbox-group>
        </Form-item>
        <Form-item class="login-no-bottom">
            <Row >
                <i-col :xs="{ span: 4, offset: 6 }" >
                    <i-button type="primary" @click="handleSubmit('formLogin')" :disabled="disabled">登录</i-button>
                </i-col>
                <i-col :xs="{ span: 4, offset: 6 }">
                    <i-button  type="primary" @click="formLoginReset('formLogin')">重置</i-button>
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
      formLogin: {
        username: '',
        password: '',
        remember: []
      },
      formLoginRules: {
        username: [
          { required: true, message: '请填写用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请填写密码', trigger: 'blur' },
          { type: 'string', min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
        ]
      },
      disabled: false
    }
  },
  methods: {
    handleSubmit (name) {
      this.$refs[name].validate((valid) => {
        sessionStorage.setItem('username', JSON.stringify(this.formLogin.username))
        if (valid) {
          // this.$Message.success('提交成功!')
          // this.$router.push({ path: '/menu' })
          console.log(this.formLogin)
          this.disabled = true
          let data = {
            username: this.formLogin.username,
            password: this.$md5(this.formLogin.password)
          }
          NetWorking.doPost(API.login, null, data).then(response => {
            this.disabled = false
            let user = response.data
            Cookie.set('auth_token', user.Secret,{ expires: 1})
            Cookie.set('user',  JSON.stringify(user), { expires: 1})
            this.$store.dispatch('createUser', user)
            this.$router.push({ path: this.$store.getters.nextUrl }).catch(err => {})
          }, (message) => {
            this.disabled = false
            // this.$Message.error('Login Failed!' + message)
            this.$Notice.error({
              title: 'Login Failed!',
              desc: message
            })
          })
        } else {
          this.$Message.error('表单验证失败!')
        }
        if (this.formLogin.remember[0] === '记住密码') {
          sessionStorage.setItem('username', JSON.stringify(this.formLogin.username))
          sessionStorage.setItem('password', JSON.stringify(this.formLogin.password))
        } else {
          sessionStorage.removeItem('username')
          sessionStorage.removeItem('password')
        }
      })
    },
    formLoginReset (name) {
      this.disabled = false
      this.$refs[name].resetFields()
    }
  },
  mounted () {
    if (sessionStorage.getItem('username')) {
      this.formLogin.username = JSON.parse(sessionStorage.getItem('username'))
    }
    if (sessionStorage.getItem('password')) {
      this.formLogin.password = JSON.parse(sessionStorage.getItem('password'))
    }
  },
  components: {
    SiteHeader
  }
}
</script>
