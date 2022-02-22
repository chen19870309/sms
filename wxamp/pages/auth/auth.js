// pages/auth/auth.js
const app = getApp();
Page({

  /**
   * 页面的初始数据
   */
  data: {

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

  onGetUserInfo: function (e) {
    //console.log(e.detail)
    if (e.detail.userInfo) {
      app.globalData.userInfo = e.detail.userInfo;
      wx.request({
        url: app.globalData.Host+'/weixin/wx77fbd12265db4add/userinfo',
        method: 'POST',
        data: {
          userid: app.globalData.userid,
          data: e.detail.userInfo
        },
        success: function(res){
          console.log(res.data)
          //wx.setStorage("NICKNAME",res.data.data.Nickname)
          //wx.setStorage("AVATAR",e.detail.userInfo.avatarUrl)
          //wx.setStorage('AUTH_WX',true)
          app.globalData.NickName = res.data.data.Nickname
          app.globalData.avatarUrl = res.data.data.Icon
          app.globalData.AuthWX = true
        }
      })
      wx.navigateTo({
        url: '/pages/index/index?page=taolun',
      })
    }
  }
})