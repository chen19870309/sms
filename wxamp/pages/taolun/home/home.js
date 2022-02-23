// pages/taolun/home/home.js
const app = getApp();
Component({
  options: {
    addGlobalClass: true,
  },
  /**
   * 组件的属性列表
   */
  properties: {
  },

  /**
   * 组件的初始数据
   */
  data: {
    canIUse: wx.canIUse('button.open-type.getUserInfo'),
    count1: 0,
    count2: 0,
    isForm: false,
    touchTime:0
  },

  ready() {
    var count1 = 0,count2 = 0
    let scopes = app.globalData.Scopes
    for(var i in scopes){
      //console.log("call ready!",scopes[i])
      if (scopes[i].title=="生字本") {
        count1 = scopes[i].cnt
      }
      if (scopes[i].title=="已学会") {
        count2 = scopes[i].cnt
      }
    }
    let mode = wx.getStorageSync('STDMODE')
    let sex = wx.getStorageSync('SEX')
    console.log(mode,sex)
    if(wx.getStorageSync("AUTH_WX")){
    //if(app.globalData.AuthWX){
      this.setData({
        NickName: app.globalData.NickName,
        avatarUrl:app.globalData.avatarUrl,
        count1:count1,
        count2:count2,
        mode: mode,
        sex:sex
      })
    }else{
      wx.navigateTo({
        url: '/pages/auth/auth',
      })
    }
  },

  /**
   * 组件的方法列表
   */
  methods: {
    showQrcode() {
      wx.previewImage({
        urls: ['https://xiaoxibaby.xyz/weixinshoukuan.jpg'],
        current: 'https://xiaoxibaby.xyz/weixinshoukuan.jpg' // 当前显示图片的http链接      
      })
    },
    onChooseAvatar(e) {
      console.log(e.detail)
      this.setData({
        avatarUrl:e.detail.avatarUrl ,
      })
    },
    dealSetting() {
      var t1 = Date.parse(new Date())
      if(t1 - this.data.touchTime > 200) {
        if(this.data.isForm) {
          wx.request({
            url: app.globalData.Host+'/weixin/wx77fbd12265db4add/userinfo',
            method: 'POST',
            data: {
              userid: app.globalData.userid,
              data: {
                sex: wx.getStorageSync('SEX'),
                avatarUrl: "",
                nickName: wx.getStorageSync('NICKNAME'),
                mode: wx.setStorageSync('STDMODE')
              }
            }, 
            header: {
              'content-type': 'application/json',
              'Authorization': 'token '+app.getLocalData('jwt_token')
            },
            success: function(res) {
              console.log(res)
            }
          })
        }
        this.setData({
          isForm: !this.data.isForm,
          touchTime: t1
        })
      }
    },
    changeMode(e) {//true为学习模式
      wx.setStorageSync('STDMODE',e.detail.value)
    },
    changeSex(e) {//true位男性
      wx.setStorageSync('SEX',e.detail.value)
    },
    changeNickName(e) {
      wx.setStorageSync('NICKNAME',e.detail.value)
      this.setData({
        NickName: e.detail.value
      })
    }
  }
})
