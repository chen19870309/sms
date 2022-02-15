// pages/taolun/log/log.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    logs: [
      {
        version:'v1.1.0',
        date: '2022/02/14',
        ctxs:[
          '1. 新增日程页面',
          '2. 记录用户的学习和测试行为',
          '3. 优化学习页面上下左右的翻页功能',
        ]
      },
    {
      version:'v1.0.0',
      date: '2022/02/08',
      ctxs:[
        '1. 主界面/学习/测试',
        '2. 生字本管理',
        '3. 用户反馈页面',
        '4. 添加生字本和已学会模块'
      ]
    }
    ]
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

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

  },
  pageBack() {
    wx.navigateBack({
      delta: 1
    });
  }
})