// pages/achieve/home/home.js
const app = getApp();
Component({
  /**
   * 组件的属性列表
   */
  properties: {

  },
  options: {
    addGlobalClass: true,
  },

  /**
   * 组件的初始数据
   */
  data: {
    elements: []
  },

  ready() {
    let that = this
    app.getScopes(()=>{
      //console.log("get scopes success:",app.globalData.Scopes)
      that.setData({
        elements : app.globalData.Scopes
      })
    })
  },

  /**
   * 组件的方法列表
   */
  methods: {

  }
})
