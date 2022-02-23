//app.js
App({
  onLaunch: function () {
    const fs = wx.getFileSystemManager()
    wx.getSystemInfo({
      success: e => {
        console.log(e)
        this.globalData.StatusBar = e.statusBarHeight;
        let capsule = wx.getMenuButtonBoundingClientRect();
        if (capsule) {
          console.log(capsule)
          this.globalData.Custom = capsule;
          this.globalData.CustomBar = capsule.bottom + capsule.top - e.statusBarHeight;
        } else {
          this.globalData.CustomBar = e.statusBarHeight + 50;
        }
      }
    })
    this.globalData.userid = wx.getStorageSync('userid')
    this.globalData.AuthWX = wx.getStorageSync("AUTH_WX") ? true:false
    this.globalData.NickName = wx.getStorageSync('NICKNAME') 
    this.globalData.AvatarUrl = wx.setStorageSync("AvatarUrl")
    try{
      fs.accessSync(this.globalData.AvatarUrl)
    }catch(e){

    }
    console.log("this.globalData.userid[",this.globalData.userid ,"]")
    // if(this.globalData.userid == undefined || this.globalData.userid <= 0) {
    //   this.authJWT()
    // }
    // 查看是否授权
    let that = this
    // wx.getSetting({
    //   success (res){
    //     if (res.authSetting['scope.userInfo']) {
    //       // 已经授权，可以直接调用 getUserInfo 获取头像昵称
    //       wx.getUserInfo({
    //         success: that.login
    //       })
    //     }else{
    //       wx.navigateTo({
    //         url: 'pages/auth/auth',
    //       })
    //     }
    //     // let appInfo = wx.getAppBaseInfo();
    //     // console.log("appInfo",appInfo)
    //     // let deviceInfo = wx.getDeviceInfo();
    //     // console.log("deviceInfo",deviceInfo)
    //     that.cacheWords()
    //   }
    // })
  },
  globalData: {
    Host: 'https://www.xiaoxibaby.xyz',
    group: 'PB',   //普通公开群组
    scope: 'common',  //基础字范围
    appInfo: {},
    deviceInfo: {},
    AuthWX:false,
    NickName: '',
    AvatarUrl: '',
    usercode: '',
    userid: '1',
    jwt: '',
    userInfo: {},
    Scopes: [],
    MyWords: [],
    word:'',
    CurWord: 0,
    Total: 0,
    LoadCount: 0,
    Round: 0, //测试或学习的次数
    Loading:false,
    FromShare: false
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
      if(that.getLocalData('jwt_token') != null){
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
            url: that.globalData.Host+'/weixin/wx77fbd12265db4add/login',
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
               that.setLocalData('jwt_token',res.data.jwt_token,3600)
               that.globalData.userid = res.data.data.Id
               wx.setStorageSync('userid', res.data.data.Id)
               wx.setStorageSync('NICKNAME', res.data.data.Nickname)
               wx.setStorageSync('AvatarUrl', res.data.data.Icon)
               that.globalData.userInfo = res.data.data
               console.log( res.data.data )
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
    let arrs = that.getLocalData('scopes')
    if(arrs != undefined && arrs != null && arrs.length > 0) {
      that.globalData.Scopes = arrs
      if (callback != undefined) {
        callback()
      }
      return
    }
    that.globalData.Loading = true
    this.authJWT().then(()=>{
      wx.request({
        url: that.globalData.Host+'/weixin/scopes',
        data: {
          userid: that.globalData.userid
        },
        header: {
          'content-type': 'application/json',
          'Authorization': 'token '+that.getLocalData('jwt_token')
        },
        success: function(res){
          that.globalData.Loading = false
          //that.globalData.Scopes = res.data.data
          let sps = res.data.data
          var arr = []
          console.log("sps",sps)
          if (sps != undefined) {
            for(var key in sps){
              let item = sps[key]
              console.log("item",item)
              arr.push({
                id: item.Id,
                title: item.Scope,
                name: item.Gp=='private'?'共'+item.Cnt+'字':'共'+item.Cnt+'字/'+item.Ucnt+'字',
                cnt: item.Cnt,
                color: item.Color,
                icon: item.Icon,
                scope: item.Scope,
                group: item.Gp
              })
            }
            arr.sort(function(a,b){
              return a.id - b.id
            })
            that.globalData.Scopes = arr
            that.setLocalData('scopes',arr,60)//缓存1分钟
          }
          if (callback != undefined) {
            callback()
          }
        },
        fail: function(res) {
          console.log("get words fail:",res)
          that.globalData.Loading = false
          if (callback != undefined) {
            callback()
          }
        }
      })
    })
  },
  setLocalData(key,data,expire) {
    var t1 = Date.parse(new Date())
    var t2 = t1 + expire*1000
    console.log(t2,key)
    wx.setStorageSync(key, data)
    wx.setStorageSync(key+'_TTL', t2)
  },
  getLocalData(key) {
    var t1 = Date.parse(new Date())
    var t2 = wx.getStorageSync(key+'_TTL')
    console.log(t2,key)
    var data = wx.getStorageSync(key)
    if(data && t2 > t1) {
      return data
    }
    return null
  },
  getwords: function(callback) {
    let that = this
    that.globalData.Loading = true
    let mode = wx.getStorageSync('STDMODE')
    if (mode != true) {
      mode = false
    }
    wx.request({
      url: that.globalData.Host+'/weixin/words',
      data: {
        scope: that.globalData.scope,
        group: that.globalData.group,
        userid: that.globalData.userid,
        word: that.globalData.word,
        mode: mode
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'token '+that.getLocalData('jwt_token')
      },
      success: function(res){
        console.log("get words success:",res.data.data)
        res.data.data.sort(function() {
          return .5 - Math.random();
        });
        if(that.globalData.word != ''){
          for(var i=0,len=res.data.data.length; i < len; i++){
            let item = res.data.data[i]
            if (item.Word == that.globalData.word) {
              that.globalData.CurWord = i
              break
            }
          }
        }else{
          that.globalData.CurWord = 0
        }
        that.globalData.MyWords = res.data.data
        that.globalData.Total = res.data.data.length
        that.cacheWords(callback)
        that.globalData.Loading = false
      },
      fail: function(res) {
        that.globalData.Loading = false
        console.log("get words fail:",res)
      }
    })
  },
  getDiary: function(callback) {
    var myDate = new Date();//获取系统当前时间
    let that = this
    let arrs = that.getLocalData('diarys')
    if(arrs != undefined && arrs != null && arrs.length>0) {
      callback(arrs)
      return 
    }
    that.globalData.Loading = true
    wx.request({
      url: that.globalData.Host+'/weixin/diary',
      data: {
        userid: that.globalData.userid,
        year: myDate.getFullYear(),
        month: myDate.getMonth()+1 
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'token '+that.getLocalData('jwt_token')
      },
      success: function(res){
        console.log("get diary success",res)
        let arr = res.data.data
        arr.sort(function(a,b){
          return b.day - a.day
        })
        that.globalData.Loading = false
        that.setLocalData('diarys',arr,60)
        callback(arr)
      },
      fail: function(res){
        console.log("get diary failed!",res)
        that.globalData.Loading = false
        callback()
      }
    })
  },
  putDiary: function(nowDate,val,callback) {
    var myDate = new Date();//获取系统当前时间
    let that = this
    wx.request({
      method:'POST',
      url: that.globalData.Host+'/weixin/diary',
      data: {
        userid: that.globalData.userid,
        year: myDate.getFullYear(),
        month: myDate.getMonth()+1 ,
        day: nowDate,
        remark: val
      },
      header: {
        'content-type': 'application/json',
        'Authorization': 'token '+that.globalData.jwt
      },
      success: function(res){
        console.log("put diary success",res)
        if (callback != undefined){
          callback(res)
        }
      },
      fail: function(res){
        console.log("put diary failed!",res)
        if (callback != undefined){
          callback(res)
        }
      }
    })
  },
  cacheWords: function(callback){
    var that = this
    const fs = wx.getFileSystemManager()
    this.globalData.LoadCount = 0
    for(var i=0,len=this.globalData.MyWords.length; i < len; i++){
      let item = this.globalData.MyWords[i]
      if (!item.Pic.startsWith("http://tmp/")){
        let imgkey = 'image_'+item.Id+'_'+item.Word
        var path = wx.getStorageSync(imgkey)
        if (path != undefined && path != '' && !path.endsWith("json")){
          try {
            fs.accessSync(path)        
            console.log("get cache path:",path)
            that.globalData.MyWords[i].Pic = path
            that.globalData.LoadCount++
            that.setLocalData('LoadCount',that.globalData.LoadCount,60)
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
            item.Pic = res.tempFilePath
            that.globalData.LoadCount++
            that.setLocalData('LoadCount',that.globalData.LoadCount,60)
          }
        })
      }
    }
    for(var i=0,len=this.globalData.MyWords.length; i < len; i++){
      let item = this.globalData.MyWords[i]
      if (!item.Sound.startsWith("http://tmp/")){
        let soundkey = 'sound_'+item.Id+'_'+item.Word
        var path = wx.getStorageSync(soundkey)
        if (path != undefined && path != '' && !path.endsWith("json")  && (path.endsWith("mp3") || path.endsWith("m4a"))){
          try {
            fs.accessSync(path)        
            console.log("get ",soundkey,"cache path:",path)
            that.globalData.MyWords[i].Sound = path
            that.globalData.LoadCount++
            that.setLocalData('LoadCount',that.globalData.LoadCount,60)
            continue
          } catch(e) {
            console.error(e)
          }
        }
        wx.downloadFile({ 
          url: decodeURI(item.Sound),
          success(res){
            let path = res.tempFilePath  
            // fs.renameSync(res.tempFilePath,path)
            console.log('声音缓存成功1', path)
            wx.setStorageSync(soundkey, path)
            item.Sound = path
            that.globalData.LoadCount++
            that.setLocalData('LoadCount',that.globalData.LoadCount,60)
          }
        })
      }
    }
    
    if (callback != undefined){
      callback()
    }
  },
  sleep: function(t) {
    return new Promise((resolve)=> setTimeout(resolve,t))
  }
})