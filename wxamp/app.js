//app.js
App({
  onLaunch: function () {
    wx.getSystemInfo({
      success: e => {
        this.globalData.StatusBar = e.statusBarHeight;
        let capsule = wx.getMenuButtonBoundingClientRect()["top"];
        if (capsule) {
          this.globalData.Custom = capsule;
          this.globalData.CustomBar = capsule.bottom + capsule.top - e.statusBarHeight;
        } else {
          this.globalData.CustomBar = e.statusBarHeight + 50;
        }
      }
    })
    // 查看是否授权
    let that = this
    wx.getSetting({
      success (res){
        if (res.authSetting['scope.userInfo']) {
          // 已经授权，可以直接调用 getUserInfo 获取头像昵称
          wx.getUserInfo({
            success: that.login
          })
        }else{
          wx.navigateTo({
            url: 'pages/auth/auth',
          })
        }
        let appInfo = wx.getAppBaseInfo();
        console.log("appInfo",appInfo)
        let deviceInfo = wx.getDeviceInfo();
        console.log("deviceInfo",deviceInfo)
        that.cacheWords()
      }
    })
  },
  globalData: {
    group: 'PB',   //普通公开群组
    scope: 'common',  //基础字范围
    appInfo: {},
    deviceInfo: {},
    NickName: '',
    AvatarUrl: '',
    usercode: '',
    userid: '1',
    jwt: '',
    userInfo: {},
    Scopes: [],
    MyWords: [],
    CurWord: 0,
    Total: 0,
    Round: 0, //测试或学习的次数
  },
  login: function(res) {
    if(res != undefined) {
     console.log(res.userInfo)
     this.globalData.NickName = res.userInfo.nickName
     this.globalData.AvatarUrl = res.userInfo.avatarUrl
     this.globalData.userInfo = res.userInfo
    }
  },
  authJWT: function() {
    let that = this
    return new Promise((resolve,reject) => {
      if(that.globalData.jwt != ''){
        resolve()
      }else{
      wx.login({
        timeout: 2000,
        fail: function(res) {
          console.log("wx.login failed:",res)
        },
        success: function(res){
          console.log("wx.login success:",res)
          var code = res.code
          wx.request({
            url: 'https://www.xiaoxibaby.xyz/weixin/wx77fbd12265db4add/login',
            method: 'POST',
            data: {
              code: code
            },
            header: {
              'content-type': 'application/json'
            },
            success: function(res){
              console.log("success:",res)
              if(res.data.data != undefined && res.data.data != null){
               that.globalData.jwt = res.data.jwt_token
               that.globalData.userid = res.data.data.Id
               that.globalData.userInfo = res.data.data
               console.log( that.globalData.jwt )
               console.log( that.globalData.userid )
              }
              resolve(res)
            },
            fail: function(res) {
              console.log("fail:",res)
              that.globalData.userid = '1'
              reject(res)
            }
          })
        }
      })
    }
    })
  },
  getScopes: function(callback) {
    let that = this
    this.authJWT().then(()=>{
      wx.request({
        url: 'https://www.xiaoxibaby.xyz/weixin/scopes',
        data: {
          userid: that.globalData.userid
        },
        success: function(res){
          //that.globalData.Scopes = res.data.data
          let sps = res.data.data
          var arr = []
          console.log("sps",sps)
          if (sps != undefined) {
            for(var key in sps){
              let item = sps[key]
              console.log("item",item)
              arr.push({
                title: item.Scope,
                name: item.Gp=='private'?'共'+item.Cnt+'字':'共'+item.Cnt+'字/'+item.Ucnt+'字',
                cnt: item.Cnt,
                color: item.Color,
                icon: item.Icon,
                scope: item.Scope,
                group: item.Gp
              })
            }
            that.globalData.Scopes = arr
          }
          if (callback != undefined) {
            callback()
          }
        },
        fail: function(res) {
          console.log("get words fail:",res)
        }
      })
    })
  },
  getwords: function(callback) {
    let that = this
    wx.request({
      url: 'https://www.xiaoxibaby.xyz/weixin/words',
      data: {
        scope: that.globalData.scope,
        group: that.globalData.group,
        userid: that.globalData.userid
      },
      success: function(res){
        console.log("get words success:",res.data.data)
        res.data.data.sort(function() {
          return .5 - Math.random();
        });
        that.globalData.CurWord = 0
        that.globalData.MyWords = res.data.data
        that.globalData.Total = res.data.data.length
        that.cacheWords()
        if (callback != undefined){
          callback()
        }
      },
      fail: function(res) {
        console.log("get words fail:",res)
      }
    })
  },
  cacheWords: function(){
    let that = this
    const fs = wx.getFileSystemManager()
    for(var i=0,len=this.globalData.MyWords.length; i < len; i++){
      let item = this.globalData.MyWords[i]
      if (!item.Pic.startsWith("http://tmp/")){
        let imgkey = 'image_cache_'+item.Word
        var path = wx.getStorageSync(imgkey)
        if (path != undefined && path != '' && !path.endsWith("json")){
          try {
            fs.accessSync(path)        
            console.log("get cache path:",path)
            this.globalData.MyWords[i].Pic = path
            continue
          } catch(e) {
            console.error(e)
          }
        }
        wx.downloadFile({ 
          url: item.Pic,
          success(res){
            console.log('图片缓存成功1', res.tempFilePath)
            wx.setStorageSync(imgkey, res.tempFilePath)
          }
        })
      }
    }
    for(var i=0,len=this.globalData.MyWords.length; i < len; i++){
      let item = this.globalData.MyWords[i]
      if (!item.Sound.startsWith("http://tmp/")){
        let soundkey = 'sound_cache_'+item.Word
        var path = wx.getStorageSync(soundkey)
        if (path != undefined && path != '' && !path.endsWith("json")){
          try {
            fs.accessSync(path)        
            console.log("get ",soundkey,"cache path:",path)
            this.globalData.MyWords[i].Sound = path
            continue
          } catch(e) {
            console.error(e)
          }
        }
        wx.downloadFile({ 
          url: item.Sound,
          success(res){
            console.log('声音缓存成功1', res.tempFilePath)
            wx.setStorageSync(soundkey, res.tempFilePath)
          }
        })
      }
    }
  }
})