// pages/index/index.js
const app = getApp();
var bgm = wx.createInnerAudioContext();
Page({

  /**
   * 页面的初始数据
   */
  data: {
    PageCur: '',
    MyWords: [],
    CurWord: 0
  },
  NavChange(e) {
    this.setData({
      PageCur: e.currentTarget.dataset.cur
    })
    if (e.currentTarget.dataset.cur=='kechengbiao') {
      // wx.navigateTo({
      //   url: '/pages/auth/auth',
      // })
      //this.showModal()
    }
  },
  showModal(e) {
    this.setData({
      modalName: 'Modal'
    })
  },
  hideModal(e) {
    this.setData({
      modalName: null
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    console.log(options)
    if(options.cur != undefined) {
    app.globalData.CurWord = options.cur
    }
    if(options.scope != undefined) {
    app.globalData.scope = options.scope
    }
    if(options.group != undefined) {
    app.globalData.group = options.group
    }
    this.setData({
      PageCur: 'achieve'
    })
    app.globalData.bgm = bgm
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {
  
  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
    wx.getSetting({
      success (res){
        if (res.authSetting['scope.userInfo']) {
          // 已经授权，可以直接调用 getUserInfo 获取头像昵称
          wx.getUserInfo({
            success: function(res){
              //console.log("onShow!",res)
            }
          })
        }else{
          wx.navigateTo({
            url: 'pages/auth/auth',
          })
        }
      }
    })
  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {
  
  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {
  
  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {
    app.getScopes()//刷新scope
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {
  
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {
    var cur = app.globalData.CurWord
    var world = app.globalData.MyWords[cur]
    var scope = app.globalData.scope
    var group = app.globalData.group
    return {
      title: '小西的学习卡片:'+world.Word,
      desc: '一起来打卡学习吧!',
      path: 'pages/index/index?scope='+scope+'&group='+group+'&cur='+cur
    }
  }
})