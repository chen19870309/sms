// pages/xuexi/study/study.js
const app = getApp();
Page({

  /**
   * 页面的初始数据
   */
  data: {
    MyWords: [],
    CurWord: 0
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    console.log(options)
    if(options.cur != undefined) {//打开的分享页面带过来的字
      app.globalData.word = options.cur
    }
    if(options.scope != undefined) {
    app.globalData.scope = options.scope
    }
    if(options.group != undefined) {
    app.globalData.group = options.group
    }
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
      path: 'pages/xuexi/study/study?scope='+scope+'&group='+group+'&cur='+world.Word
    }
  }
})