<template>
<div>
<!-- 粒子背景 -->
<vue-particles class="bg"></vue-particles>
<blog-header></blog-header>
<div class="func-group">
<Button type="info" @click="funcjson">格式化json</Button>
<Button type="info" @click="encodeb64">Base64编码</Button>
<Button type="info" @click="decodeb64">Base64解码</Button>
<Button type="info" @click="funcmd5">计算MD5</Button>
<Button type="info" @click="genqrcode">生成二维码</Button>
<Button type="info" @click="cropperpic">图片按比例切割</Button>
</div>
<div class="demo-split">
<Split v-model="split1">
    <div slot="left" class="demo-split-pane">
        <span>{{ inputlabel }}</span><br/>
        <Input v-if="!qrcode && !cropper" v-model="inputdata" type="textarea" :autosize="{minRows: 12,maxRows: 30}" placeholder="Enter something..." />
        <Input v-if="qrcode"  v-model="qr.url" placeholder="二维码内容:" />
        <span v-if="qrcode">添加小图标</span><br/>
        <Input v-if="qrcode"  v-model="qr.icon" placeholder="二维码图标:" />
        <div class="cropper-content" v-if="cropper">
        <label for="h2w">图片长宽比:</label>
        <div id="h2w" ><Input v-model="option.fixedH" type="number" placeholder="长" @on-change="changePicRatio"/>:<Input v-model="option.fixedW" placeholder="宽" type="number"  @on-change="changePicRatio" /></div>
        <Button @click="choiceImg" icon="ios-cloud-upload-outline" type="primary">选择图片</Button>
        <input name="file" type="file" ref="uploader" style="position:absolute; clip:rect(0 0 0 0);" accept="image/png, image/jpeg, image/gif, image/jpg" @change="uploadImg($event)" />
        </div>
    </div>
    <div slot="right" class="demo-split-pane">
        <span>{{ outputlabel }}</span><br/>
        <json-viewer :value="jsonData" v-if="bjson"
        :expand-depth=5
        copyable
        boxed
        sort></json-viewer>
        <Input v-model="outB64" v-if="eb64" type="textarea" :autosize="{minRows: 12,maxRows: 30}"  />
        <Input v-model="outText" v-if="db64" type="textarea" :autosize="{minRows: 12,maxRows: 30}" />
        <Input v-model="outMd5" v-if="emd5" type="textarea" :autosize="{minRows: 12,maxRows: 30}" />
        <vue-qr v-if="qrcode && qrUrl!= '' && qrIcon == ''" :text="qrUrl" :margin="0" colorDark="#003399" colorLight="#ffffcc" :logoScale="0.3" :size="256"></vue-qr>
        <vue-qr v-if="qrcode && qrUrl!= '' && qrIcon != ''" :text="qrUrl" :margin="0" colorDark="#003399" colorLight="#ffffcc" :logoSrc="qrIcon + '?cache'" :logoScale="0.3" :size="256"></vue-qr>
        <div class="cropper-content">
        <div class="cropper"  v-if="cropper && option.img!= ''">
        <vueCropper
            ref="cropper"
            :img="option.img"
            :outputSize="option.size"
            :outputType="option.outputType"
            :info="option.info"
            :full="option.full"
            :canMove="option.canMove"
            :canMoveBox="option.canMoveBox"
            :original="option.original"
            :autoCrop="option.autoCrop"
            :fixed="option.fixed"
            :fixedNumber="option.fixedNumber"
            :centerBox="option.centerBox"
            :infoTrue="option.infoTrue"
            :fixedBox="option.fixedBox"
          ></vueCropper>
          <Progress slot='progressBar' :percent="up_percent"></Progress>
          <Button type="info" @click="qiniuUpdate">上传图片</Button>
        </div>
        </div>
    </div>
</Split>
</div>
<blog-footer></blog-footer>
</div>
</template>

<script>
import vueQr from 'vue-qr'
import JsonViewer from 'vue-json-viewer'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/layout/BlogFooter'
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
const Base64 = require('js-base64').Base64;
export default {
  data () {
    return {
      split1: 0.5,
      inputdata: '',
      inputlabel: '请输入json',
      outputlabel: '格式化后的json',
      up_percent: 0,
      exam: null,
      qr: {
        url:'',
        icon: ''
      },
      qiniu: {

      },
      bjson: true,
      eb64: false,
      db64: false,
      emd5: false,
      qrcode: false,
      cropper:false,
      preview: '',
      option: {
        img: '', // 裁剪图片的地址
        size: 1,
        info: true, // 裁剪框的大小信息
        outputSize: 0.8, // 裁剪生成图片的质量
        outputType: 'jpeg', // 裁剪生成图片的格式
        canScale: true, // 图片是否允许滚轮缩放
        autoCrop: true, // 是否默认生成截图框
        // autoCropWidth: 140, // 默认生成截图框宽度
        // autoCropHeight: 200, // 默认生成截图框高度
        fixedBox: false, // 固定截图框大小 不允许改变
        fixed: true, // 是否开启截图框宽高固定比例
        fixedNumber: [10, 7], // 截图框的宽高比例
        fixedH:10,
        fixedW:7,
        fileaa:null,
        full: true, // 是否输出原图比例的截图
        canMoveBox: false, // 截图框能否拖动
        original: false, // 上传图片按照原始比例渲染
        centerBox: true, // 截图框是否被限制在图片里面
        infoTrue: false // true 为展示真实输出图片宽高 false 展示看到的截图框宽高
      },
    }
  },
  computed: {
    jsonData: function() {
      if (this.inputdata === '') {
        return {}
      } else {
        return JSON.parse(this.inputdata)
      }
    },
    outB64: function() {
      if (this.inputdata === '') {
        return ""
      } else {
        return Base64.encode(this.inputdata)
      }
    },
    outText: function() {
      if (this.inputdata === '') {
        return ""
      } else {
        return Base64.decode(this.inputdata)
      }
    },
    outMd5: function() {
      if (this.inputdata === '') {
        return ""
      } else {
        return this.$md5(this.inputdata)
      }
    },
    qrUrl: function() {
      if (this.qr.url === '') {
        return ""
      } else {
        return this.qr.url
      }
    },
    qrIcon: function() {
      if (this.qr.icon === '') {
        return ""
      } else {
        return this.qr.icon
      }
    }
  },
  methods: {
    funcjson () {
      this.inputdata = ''
      this.inputlabel = '请输入json'
      this.outputlabel = '格式化后的json'
      this.bjson = true
      this.eb64 = false
      this.db64 = false
      this.emd5 = false
      this.qrcode = false
      this.cropper = false
    },
    encodeb64 () {
      this.inputdata = ''
      this.inputlabel = '请输入明文'
      this.outputlabel = 'Base64格式'
      this.bjson = false
      this.eb64 = true
      this.db64 = false
      this.emd5 = false
      this.qrcode = false
      this.cropper = false
    },
    decodeb64 () {
      this.inputdata = ''
      this.inputlabel = '请输入Base64'
      this.outputlabel = '明文格式'
      this.bjson = false
      this.eb64 = false
      this.db64 = true
      this.emd5 = false
      this.qrcode = false
      this.cropper = false
    },
    funcmd5 () {
      this.inputdata = ''
      this.inputlabel = '请输入明文'
      this.outputlabel = 'MD5:'
      this.bjson = false
      this.eb64 = false
      this.db64 = false
      this.emd5 = true
      this.qrcode = false
      this.cropper = false
    },
    genqrcode () {
      this.inputlabel = '生成二维码参数:'
      this.outputlabel = 'QRCode:'
      this.bjson = false
      this.eb64 = false
      this.db64 = false
      this.emd5 = false
      this.qrcode = true
      this.cropper = false
    },
    cropperpic () {
      this.inputlabel = '按比例切割图片:'
      this.outputlabel = 'PicUrl:'
      this.bjson = false
      this.eb64 = false
      this.db64 = false
      this.emd5 = false
      this.qrcode = false
      this.cropper = true
      this.up_percent = 0
    },
    choiceImg(){
      this.$refs.uploader.dispatchEvent(new MouseEvent('click')) 
    },
    uploadImg(e) {
      if (this.qiniu.token == undefined || this.qiniu.token == ''){
        this.chechQiniuToken()
      }
      this.exam=e
      var file = e.target.files[0]
      if (!/\.(gif|jpg|jpeg|png|bmp|GIF|JPG|PNG)$/.test(e.target.value)) {
        alert('图片类型必须是.gif,jpeg,jpg,png,bmp中的一种')
        return false
      }
      var reader = new FileReader()
      reader.onload = (e) => {
        let data
        if (typeof e.target.result === 'object') {
          // 把Array Buffer转化为blob 如果是base64不需要
          data = window.URL.createObjectURL(new Blob([e.target.result]))
        } else {
          data = e.target.result
        }
        this.option.fixedNumber = [this.option.fixedW, this.option.fixedH]
        this.autoCrop = true
        this.fixed = true
        this.option.img = data
      }
      console.log(file)
      reader.readAsDataURL(file)
    },
    fileSize(base64){
      let fileSize;
      //找到等号，把等号去掉
      if (base64.indexOf('=') > 0) {
        var indexOf = base64.indexOf('=');
        base64 = base64.substring(0, indexOf);//把末尾的’=‘号去掉
      }
      fileSize = parseInt(base64.length - (base64.length / 8) * 2);
      return fileSize;
    },
    qiniuUpdate(){
      this.$refs.cropper.getCropData((data)=>{
        let that = this 
        this.option.img = data;
        let pic = data.replace(/^.*?,/, '');
        let size = this.fileSize(pic);
        var file = this.exam.target.files[0]
        let key = this.qiniu.prefix + file.name;
        let url = "https://upload.qiniup.com/putb64/"+size+"/key/"+Base64.encode(key);   
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange=function(){
          if (xhr.readyState==4){//上传成功
            console.log(xhr.responseText)
            var res = JSON.parse(xhr.responseText)
            if (res.key != undefined && res.key != ''){
            that.up_percent = 100
            that.$Notice.success({
              title: '图片上传成功'+that.qiniu.domain + key,
            });
            that.outputlabel = that.qiniu.domain + key
            }else{
              that.$Notice.error({
              title: '图片上传失败',
              desc: '错误信息:' + xhr.responseText
            });
            }
          }
        }
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/octet-stream");
        xhr.setRequestHeader("Authorization", "UpToken "+this.qiniu.token);
        xhr.send(pic);
      })
    },
    chechQiniuToken(){
      NetWorking.doGet(API.uptoken).then(response => {
        console.log('UpToken:', response.data)
        this.qiniu = response.data
      }, (message) => {
        this.$Message.error('Get Qiniu Token Failed!' + message)
      })
    },
    changePicRatio() {
      if (this.option.fixedH >0 && this.option.fixedW > 0){
      if(this.exam != null){
        this.$refs.cropper.clearCrop();
        this.$refs.cropper.startCrop();
        //this.uploadImg (this.exam);
      }else if(this.option.img != ''){
        this.$refs.cropper.clearCrop();
        this.$refs.cropper.startCrop();
      }
      }
    }
  },
  components: {
    JsonViewer,
    BlogHeader,
    BlogFooter,
    vueQr
  }
}
</script>

<style>
.demo-split{
  text-align: left;
  margin: 3%;
  margin-top:0px;
  min-height: 600px;
  border: 1px solid #dcdee2;
  display: grid;
}
.demo-split-pane{
  margin: 5px;
  padding: 10px;
  min-height: 600px;
}
.func-group{
  margin: 3%;
  padding: 10px;
  border: 1px solid #ccdede;
  background-color: beige;
}
.cropper-content .cropper {
  width: auto;
  height: 300px;
}
</style>
