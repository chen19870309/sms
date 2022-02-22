// pages/xuexi/book/book.js
const app = getApp();
Page({

  /**
   * 页面的初始数据
   */
  data: {
    title:'三只小猪',
    loading:true,
    nowPgae:1,
    startX:0,
    slider:false,
    pic: '',
    animation: '',
    animationData:{},
    cardInfoList: [{ name: 1}, { name: 2}, { name: 3}, { name: 4}]
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    let that = this
    app.getScopes(()=>{
      //console.log("get scopes success:",app.globalData.Scopes)
      that.setData({
        elements : app.globalData.Scopes
      })
    })
    app.getwords(()=>{
      app.cacheWords(()=>{
      var cur = app.globalData.CurWord
        that.setData({
          animation: '',
          loading: false,
          pic: app.globalData.MyWords[cur].Pic
        })
    })
  })
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

  touchStart: function(e){
    this.setData({
        startX: e.changedTouches[0].clientX,
    })
  },
  touchEnd: function(e) {
    let that=this;
    let startX = this.data.startX;
    let endX = e.changedTouches[0].clientX;
    if (this.data.slider)return;
    if (app.globalData.Total == 1) return
    let c = 1
    if (startX - endX > 30){
      c = -1
    }
    var cur = app.globalData.CurWord
    if (app.globalData.CurWord>0) {
      cur = (app.globalData.CurWord + c)%app.globalData.Total
    }else{
      cur = (app.globalData.Total + c)%app.globalData.Total
      if (cur < 0){
        cur = app.globalData.Total - 1
      }
    }
    app.globalData.CurWord = cur


        //创建动画   5s将位置移动到-150%,-150%
        let animation = wx.createAnimation({
            duration: 500,
        });
        animation.translate(c*500,-500).rotate(c*100).step();
        this.setData({
            animation: animation.export()
        });

        // 移动完成后
        setTimeout(function(){
            //创建动画   将位置归位
            let animation = wx.createAnimation({
                duration: 500,
            });
            animation.translate(0,0).rotate(0).step();

            that.setData({
                cardInfoList: that.data.cardInfoList,
                animation: animation.export(),
                slider:false,
                nowPgae:that.data.nowPgae+1,
                pic: app.globalData.MyWords[cur].Pic
            });
        },500)
    }

})