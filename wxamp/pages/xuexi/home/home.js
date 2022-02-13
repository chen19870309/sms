// pages/xuexi/home/home.js
const app = getApp();
var startX, endX,startY, endY;
var moveFlag = true;
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    group: String,
    scope: String
  },

  /**
   * 组件的初始数据
   */
  options: {
    addGlobalClass: true,
  },
  data: {
    pic: '',
    word: '',
    pinYin: '',
    message: ''
  },
  ready() {
    let that = this
    // that.setData({
    //   animation: '',
    //   pic: app.globalData.MyWords[0].Pic,
    //   word: app.globalData.MyWords[0].Word,
    //   pinYin: app.globalData.MyWords[0].PinYin
    // })
    app.getwords(()=>{
      var cur = app.globalData.CurWord
      app.globalData.Round = 5*app.globalData.Total
      that.setData({
      animation: '',
      pic: app.globalData.MyWords[cur].Pic,
      word: app.globalData.MyWords[cur].Word,
      pinYin: app.globalData.MyWords[cur].PinYin
    })
    that.speeker()
    })
  },
  /**
   * 组件的方法列表
   */
  methods: {
    showModal(msg) {
      this.setData({
        modalName: 'Info',
        message: msg
      })
    },
    hideModal(e) {
      this.setData({
        modalName: null
      })
      let scope = app.globalData.scope
      let group = app.globalData.group
      wx.navigateTo({
        url: '/pages/index/index?scope='+scope+'&group='+group
      })
    },
    goexam() {
      let scope = app.globalData.scope
      let group = app.globalData.group
      wx.navigateTo({
        url: '/pages/xuexi/exam/exam?scope='+scope+'&group='+group
      })
    },
    touchStart(e) {
      startX = e.touches[0].pageX;
      startY = e.touches[0].pageY;
      moveFlag = true;
    },
    speeker(){
      var cur = app.globalData.CurWord
      var bgm = app.globalData.bgm
      if (app.globalData.MyWords[cur].Sound != ''){
        bgm.src=app.globalData.MyWords[cur].Sound
        bgm.play()
      }
    },
    touchMove(e) {
      endX = e.touches[0].pageX;
      endY = e.touches[0].pageY;
      if (moveFlag) {
        moveFlag = false;
        //console.log("startX - endX:",startX - endX);
        if (startX - endX > 10) {
          this.movePage(-1)
          return
        }
        if (endX - startX > 10) {
          this.movePage(1)
          return 
        }
        if (startY - endY > 10) {
          this.movePage(-2)
          return
        }
        if (endY - startY > 10) {
          this.movePage(2)
          return 
        }
      }
    },
    touchEnd(e){
      moveFlag = true;
    },
    movePage(c) {
      let that = this
      var cur = 0
      let updown = false
      if(c == -2 || c == 2){
        updown = true
        c = c/2
      }
      app.globalData.Round --
      if (app.globalData.Round <= 0) {
        that.showModal('好棒😄,本次已经学完'+app.globalData.Total+'个字了!休息一会吧🎉🎉🎉')
        //记录diary
        var myDate = new Date();//获取系统当前时间
        var nowDate = myDate.getDate();
        var str = ""
        for(var i=0;i<app.globalData.Total;i++){
          str += " "+app.globalData.MyWords[i].Word
          if(i>0&&i%5==0) {
            str += "\n"
          }
        }
        app.putDiary(nowDate,'🎉本次学习的内容是:\n'+str)
        return 
      }
      if (app.globalData.Total == 1) return
     console.log("movePage:",c,"|",cur,updown)
      if (app.globalData.CurWord>0) {
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
      if(updown){
        animation.opacity(0.2).translate(0,c*100).step()
        animation.opacity(1.0).translate(0,0).step()
      }else{
        animation.opacity(0.2).translate(c*100,0).step()
        animation.opacity(1.0).translate(0,0).step()
      }
      that.setData({
        animation: animation.export(),
      })
      var bgm = app.globalData.bgm;
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
  }
})
