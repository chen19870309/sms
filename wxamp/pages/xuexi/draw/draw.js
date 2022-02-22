// pages/xuexi/draw/draw.js
const app = getApp();
Component({
  /**
   * 组件的属性列表
   */
  properties: {

  },

  /**
   * 组件的初始数据
   */
  options: {
    addGlobalClass: true,
  },
  data: {
    btcolors:['red','orange','olive','mauve'],
    canvas: '',
    ctx: null,
    word:'?',
    canvasWidth:'600',
    canvasHeight: 0,
    boldVal:18,
    colors: ['black', 'pink', 'red', 'skyblue', 'greenyellow', '#00FF00', '#0000FF', '#FF00FF','#00FFFF',
      '#FFFF00', '#70DB93', '#5C3317 ', '#9F5F9F', '#B5A642', '#FF7F00','#42426F'],
    curColor:'#000000',
    isEraser:false,
    isSuc:true,
    isTap:true,
    words:['紫','气','东','来','时','来','远','转','','','','','','','',''],
    list:[]
  },

  ready() {
    let that = this
    var cur = app.globalData.CurWord
    var selections = []
    var colors = this.data.btcolors
    let len = this.data.canvasWidth/2
      that.setData({
        canvasWidth: len,
        canvasHeight: len
      })
      that.isMouseDown=false
      that.lastLoc={ x: 0, y: 0 }
      that.lastTimestamp = 0;
      that.lastLineWidth = -1;
    colors.sort(function() {
      return .5 - Math.random();
    });
    for (var i = 0;i<4;i++){
      selections.push({
        color: colors[i],
        word: app.globalData.MyWords[(cur+i)%app.globalData.Total].Word
      })
    }
    selections.sort(function() {
      return .5 - Math.random();
    });
      this.setData({
        list: selections
      })
      that.drawBgColor();   
  },
  /**
   * 组件的方法列表
   */
  methods: {
    getCanvas(id) {
      let selector = '#xuexi_'+id
      let that = this
      return new Promise((resolve, reject) => {
          if(that.data.ctx != null) {
            resolve({
              word: that.data.words[id],
              ctx: that.data.ctx,
              _size: that.data._size
            })
          }else{
          that.createSelectorQuery()
          .select(selector)
          .fields({
              size: true,
              node: true
          })      
          .exec((res) => {  
              if (res && res[0]) {
                const canavs = res[0].node
                const ctx = canavs.getContext('2d')
                let dpr = wx.getSystemInfoSync().pixelRatio
                canavs.width = res[0].width *dpr
                canavs.height = res[0].height *dpr
                ctx.scale(dpr,dpr)
                that.setData({
                  ctx: ctx,
                      _size: {
                        w: res[0].width,
                        h: res[0].height,
                        d: dpr
                      }
                })
                  resolve({
                      word: that.data.words[id],
                      ctx: ctx,
                      _size: {
                        w: res[0].width,
                        h: res[0].height,
                        d: dpr
                      }
                  })
              } else {
                  resolve(null)
              }           
          });
      }});
    },
    clearCanavs(res) {
      res.ctx.clearRect(0, 0, res._size.w, res._size.h);
    },
    drawWord(res) {
      let ctx = res.ctx
      ctx.fillStyle = 'rgba(250, 250, 250, 0.6)'
      ctx.strokeStyle = 'red'
      ctx.lineWidth = 1
      let w = res._size.w;
      let h = w;
      ctx.fillRect(0,0,w,h)
      res.ctx.beginPath()
      ctx.strokeRect(0,0,w,h)

      ctx.setLineDash([3, 3]);
      ctx.moveTo(0 , 0);
      ctx.lineTo(w , h);

      ctx.moveTo(w , 0);
      ctx.lineTo(0 , h);

      ctx.moveTo(w/2,0)
      ctx.lineTo(w/2 , h);

      ctx.moveTo(0,h/2)
      ctx.lineTo(w , h/2);
      res.ctx.closePath()
      ctx.stroke();

      ctx.font= (w-15)*this.data.boldVal/20 + "px " + "KaiTi,KaiTi_GB2312,AR PL UKai CN,simkai,楷体,STKaiti";
      ctx.fillStyle = this.data.curColor;
      ctx.strokeStyle = this.data.curColor;
      ctx.textBaseline="middle";
      ctx.textAlign="center";
      //ctx.strokeText('好', w/2, h/2);
      if(res.word != undefined && res.word != null && res.word != ''){
        ctx.fillText(res.word, w/2, h/2);
      }
    },
    drawBgColor(){
      let that = this
      for(let i=0;i<4;i++) {
        this.getCanvas(i).then((res)=>{
          console.log(res)
          that.clearCanavs(res)
          that.drawWord(res)
        })
      }
    },
    changeBold:function(e){
      console.log(e.detail.value)
      this.setData({boldVal: e.detail.value})
      this.drawBgColor()
    },
    selectColor:function(e){
      console.log("selectColor",e)
      this.setData({ curColor: e.currentTarget.dataset.value })
      this.setData({ isEraser: false })
      this.drawBgColor()
    },
    beginStroke(event) {
        this.isMouseDown = true
        this.lastLoc = { x: event.touches[0].x, y: event.touches[0].y }
        this.lastTimestamp = event.timeStamp;
        this.setData({ isTap: true })
        // //draw
        // let that = this
        // this.getCanvas('#xuexi').then((res)=>{
        //   res.ctx.arc(that.lastLoc.x, that.lastLoc.y, that.data.boldVal / 2, 0, 2 * Math.PI)
        //   res.ctx.fillStyle=that.data.curColor;
        //   res.ctx.fill();
        // })
   
        // if (event.touches.length>1){
        //   var xMove = event.touches[1].x - event.touches[0].x;
        //   var yMove = event.touches[1].y - event.touches[0].y;
        //   this.lastDistance = Math.sqrt(xMove * xMove + yMove * yMove);
   
        // }
       
    },
   
    endStroke(event) {
      // console.log(this.data.isTap)
      // if (this.data.isTap){
      //   this.lastLoc = { x: event.changedTouches[0].x, y: event.changedTouches[0].y }
      //   this.lastTimestamp = event.timeStamp;
      //   //draw
      //   this.context.arc(this.lastLoc.x, this.lastLoc.y, this.data.boldVal / 2, 0, 2 * Math.PI)
      //   this.context.setFillStyle(this.data.curColor);
      //   this.context.fill();
      //   wx.drawCanvas({
      //     canvasId: 'canvas',
      //     reserve: true,
      //     actions: this.context.getActions() // 获取绘图动作数组
      //   })
   
      // }
      this.isMouseDown= false
    },
   
    moveStroke(event) {
      let that = this
      if (this.isMouseDown && event.touches.length == 1) {
        var touch = event.touches[0];
        var curLoc = { x: touch.x, y: touch.y };
        var curTimestamp = event.timeStamp;
        var s = this.calcDistance(curLoc, this.lastLoc)
        var t = curTimestamp - this.lastTimestamp;
        var lineWidth = this.calcLineWidth(t, s)
        console.log("dith=",s)
        if(s<1) return
        //draw
        this.getCanvas('#xuexi').then((res)=>{
          console.log(res.ctx)
        res.ctx.strokeStyle=that.data.curColor;
        res.ctx.lineWidth=lineWidth;
        res.ctx.beginPath()
        res.ctx.moveTo(this.lastLoc.x, this.lastLoc.y)
        res.ctx.lineTo(curLoc.x, curLoc.y)   
        // locHistory.push({ x: curLoc.x, y: curLoc.y, with: lineWidth, t: t })
        res.ctx.lineCap= "round"
        res.ctx.lineJoin="round"
        res.ctx.closePath()
        res.ctx.stroke();
        //res.ctx.save()
        //res.ctx.clip()
        //that.drawWord(res)
        //res.ctx.restore()
        //res.ctx.save()
        })
    
        this.lastLoc=curLoc;
        // this.setData({ lastTimestamp: curTimestamp })
        // this.setData({ lastLineWidth: lineWidth })
      
      } else if (event.touches.length > 1){
        this.setData({isTap:false})
   
        var xMove = event.touches[1].x - event.touches[0].x;
        var yMove = event.touches[1].y - event.touches[0].y;
        var newdistance = Math.sqrt(xMove*xMove + yMove*yMove);
        // if (newdistance - this.lastDistance>0){
        //   this.setData({ canvasWidth: this.data.canvasWidth * 1.2 })
        //   this.setData({ canvasHeight: this.data.canvasHeight * 1.2 })
        // }else{
        //   this.setData({ canvasWidth: this.data.canvasWidth * 0.8 })
        //   this.setData({ canvasHeight: this.data.canvasHeight * 0.8})
        // }
   
      }
    },
    calcLineWidth(t, s){
      var v = s / t;
      var resultLineWidth = this.data.boldVal;
      if(v <= 0.1) {
        resultLineWidth = resultLineWidth * 1.2;
      }else if(v >= 10) {
        resultLineWidth = resultLineWidth/1.2
      }else{
        resultLineWidth = resultLineWidth - (v - 0.1) / (10 - 0.1) * (resultLineWidth * 1.2 - resultLineWidth / 1.2)
      }
      return resultLineWidth
    },
   calcDistance(loc1, loc2) {
      return Math.sqrt((loc1.x - loc2.x) * (loc1.x - loc2.x) + (loc1.y - loc2.y) * (loc1.y - loc2.y))
    },
   clearCanvas:function(){
     this.drawBgColor()
     this.setData({ isEraser:false})
    },
    toggle: function(e) {
      let that = this
      var word = e.currentTarget.dataset.word;
      var cur = app.globalData.CurWord
      var w = app.globalData.MyWords[cur].Word
      if (word == w){
        that.setData({
          word:w
        })
        that.drawBgColor()
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
   eraser:function () {
     this.setData({ isEraser: !this.data.isEraser})
     this.setData({ curColor:'#ffffff'})
    // this.context.clearActions();
    // this.context.draw()
   },
   saveImg:function(){
     var that = this;
     var keys = wx.getStorageSync('myImg') || [];
     console.log(that.data.keyLength)
     wx.canvasToTempFilePath({
       canvasId: 'canvas',
       success: function (response) {
         keys.unshift(response.tempFilePath)
         wx.setStorage({
           key: 'myImg',
           data: keys,
         })
         that.setData({ isSuc: false })
         
       }
     })
   }
  },
  canvasToTempFilePath: function(
    {
        canvasId,
        quality = 1,
        fileType = 'jpg',
    },
    context
  ) {
    new Promise((resolve, reject) => {
      wx.canvasToTempFilePath({
          canvasId,
          quality,
          fileType,
          success: resolve,
          fail: reject,
      }, context);
  });
}
})
