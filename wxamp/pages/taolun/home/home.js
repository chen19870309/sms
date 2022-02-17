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
    count2: 0
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
    //if(wx.getStorage("AUTH_WX")){
    if(app.globalData.AuthWX){
      this.setData({
        NickName: app.globalData.NickName,
        avatarUrl:app.globalData.avatarUrl,
        count1:count1,
        count2:count2
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
    }
  }
})
