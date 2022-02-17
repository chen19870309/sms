// pages/xuexi/exam/exam.js
const app = getApp();
Page({
  /**
   * 页面的初始数据
   */
  data: {
    colors:['red','orange','olive','mauve'],
    pic: '',
    word: '',
    pinYin: '',
    sound: '',
    message: '',
    content: '',
    modelName: '',
    ready: false,
    list: [],
    didwords: [],
    first:true,
    loading:true
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    console.log("scope=",options)
    if (options.cnt <=0) {
      this.setData({ready:false})
      if (options.scope == '生字本') {
        this.showModal('🎉🎉🎉学完了🎉🎉🎉!','没有更多的生字了^_^')
      }else{
        this.showModal('你选择的内容为空!')
      }
    }else{
    this.setData({ready:true})
    if (options.scope!=undefined && options.group!=undefined){
      app.globalData.scope =  options.scope
      app.globalData.group =  options.group
      this.setData({loading:true})
      let that = this
      app.getwords(()=>{
        app.cacheWords(()=>{
            that.setData({loading:false})
            that.initExam()
            that.showCheckModal('开始学习还是测试？')
        })
      })
    }
  }
  },
  showModal(message,ctx) {
    this.setData({
      modalName: 'Info',
      message: message,
      content: ctx
    })
  },
  showCheckModal(message) {
    this.setData({
      modalName: 'check',
      message: message
    })
  },
  startExam() {
    this.setData({
      modalName: null
    })
    app.globalData.Round = 3*app.globalData.Total
    this.initExam()
    this.speeker()
  },
  startStudy() {
    wx.navigateTo({
      url: '/pages/xuexi/study/study',
    })
  },
  hideModal(e) {
    this.setData({
      modalName: null
    })
    wx.navigateTo({
      url: '/pages/index/index?page=achieve',
    })
  },
  initExam() {
    var cur = app.globalData.CurWord
    var selections = []
    for (var i = 0;i<4;i++){
      selections.push({
        color: this.data.colors[i],
        word: app.globalData.MyWords[(cur+i)%app.globalData.Total].Word
      })
    }
    selections.sort(function() {
      return .5 - Math.random();
    });
    this.setData({
     list: selections,
     animation: '',
     pic: app.globalData.MyWords[cur].Pic,
     word: app.globalData.MyWords[cur].Word,
     pinYin: app.globalData.MyWords[cur].PinYin
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
    var cur = app.globalData.CurWord
    var world = app.globalData.MyWords[cur]
    var scope = app.globalData.scope
    var group = app.globalData.group
    return {
      title: '小西的测试:'+scope,
      desc: '一起来打卡学习吧!',
      path: 'pages/xuexi/exam/exam?scope='+scope+'&group='+group+'&cur='+world.Word
    }
  },
  speeker: function() {
    var cur = app.globalData.CurWord
    var bgm = app.globalData.bgm
    if (app.globalData.MyWords[cur].Sound != ''){
      bgm.src=app.globalData.MyWords[cur].Sound
      bgm.play()
    }
  },
  toggle: function(e) {
    //console.log(e);
    let that = this
    var word = e.currentTarget.dataset.word;
    var cur = app.globalData.CurWord
    var w = app.globalData.MyWords[cur].Word
    if (word == w){
      app.globalData.bgm.pause()
      if (that.data.first) {//一次成功
        that.data.didwords.push(word)
        console.log("got:",word)
      }else{
        that.setData({first:true})
      }
      app.globalData.Round --
      if (app.globalData.Round <= 0) {
        var len=that.data.didwords.length
        var score = 100*(len/3/app.globalData.Total)
        var maps = that.countWord(that.data.didwords)
        console.log(maps)
        let userid = parseInt(app.globalData.userid)
        for(var i=0,j=app.globalData.MyWords.length;i<j;i++){
          let w = app.globalData.MyWords[i]
          var key = w.Word
          var status = 0
          if (maps[key] == 3) {
            status = 1
          }
          wx.request({
            method: 'post',
            url: app.globalData.Host+'/weixin/words',
            data: {
              Id: w.Id,
              userid: userid,
              status: status
            },
            success: function(res) {
              console.log('post words success!',res)
            },
            fail: function(res) {
              console.log('post words failed!',res)
            }
          })
        }
        var str = '本次得分是：'+Math.ceil(score)+'分\n'
        that.showModal('🎉本次测试完成🎉',str)
        if (score == 100) {
          str += "太棒了🎈"
          app.globalData.bgm.src="/pages/xuexi/exam/success.m4a"
          app.globalData.bgm.play()
        }else{
          str += "字/[答错次数]:"
        }
        //记录diary
        var myDate = new Date();//获取系统当前时间
        var nowDate = myDate.getDate();
        for(var i=0;i<app.globalData.Total;i++){
          var w = app.globalData.MyWords[i].Word
            var c= 3-maps[w]
            if(c > 0) {
              str += " "+w+"["+c+"]"
            }
        }
        app.putDiary(nowDate,'🎉完成测试🎉:\n'+str)
        wx.setStorageSync('scopes', null)
      }else{
      that.nextWord()
      that.initExam()
      }
    }else{
    that.setData({
      animation: word,
      first: false,
    })
    var bgm = app.globalData.bgm
    bgm.src = "/pages/xuexi/exam/fail.m4a"
    bgm.play()
    setTimeout(function() {
      that.setData({
        animation: ''
      })
    }, 1000)
  }
  },
  countWord(arr){
    return arr.reduce(function (prev,next) {
      prev[next] = (prev[next] + 1) || 1
      return prev
    },{})
  },
  nextWord: function() {
    let that = this
    var cur = 0
    var c = 1
     //console.log("movePage:",c,"|",cur,app.globalData.CurWord)
    if (app.globalData.CurWord>=0) {
      cur = (app.globalData.CurWord + c)%app.globalData.Total
    }else{
      cur = (app.globalData.Total + c)%app.globalData.Total
      if (cur < 0){
        cur = app.globalData.Total - 1
      }
      }
      app.globalData.CurWord = cur
      var animation = wx.createAnimation({
        duration:500,
        timingFunction: 'ease',
        delay: 50,
      });
      animation.opacity(0.2).translate(c*100,0).step()
      animation.opacity(1.0).translate(0,0).step()
      that.setData({
        animation: animation.export(),
      })
      var bgm = app.globalData.bgm
      setTimeout( function() {
        that.setData({
          animation: '',
          pic: app.globalData.MyWords[cur].Pic,
          word: app.globalData.MyWords[cur].Word,
          pinYin: app.globalData.MyWords[cur].PinYin
        })
        if (app.globalData.MyWords[cur].Sound != ''){
          bgm.src=app.globalData.MyWords[cur].Sound
          bgm.play()
        }
      },300)
    }
})