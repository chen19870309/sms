<template>
<div>
<blog-header></blog-header>
<div class="func-group">
<Button type="info" @click="funcjson">格式化json</Button>
<Button type="info" @click="encodeb64">Base64编码</Button>
<Button type="info" @click="decodeb64">Base64解码</Button>
<Button type="info" @click="funcmd5">计算MD5</Button>
</div>
<div class="demo-split">
<Split v-model="split1">
    <div slot="left" class="demo-split-pane">
        <span>{{ inputlabel }}</span><br/>
        <Input v-model="inputdata" type="textarea" :autosize="{minRows: 12,maxRows: 30}" placeholder="Enter something..." />
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
    </div>
</Split>
</div>
<blog-footer></blog-footer>
</div>
</template>

<script>
import JsonViewer from 'vue-json-viewer'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/layout/BlogFooter'
const Base64 = require('js-base64').Base64;
export default {
  data () {
    return {
      split1: 0.5,
      inputdata: '',
      inputlabel: '请输入json',
      outputlabel: '格式化后的json',
      bjson: true,
      eb64: false,
      db64: false,
      emd5: false,
    }
  },
  computed: {
    jsonData: function() {
      if(this.inputdata === '') {
        return {}
      }else{
        return JSON.parse(this.inputdata)
      }
    },
    outB64: function() {
      if(this.inputdata === '') {
        return ""
      }else{
        return Base64.encode(this.inputdata)
      }
    },
    outText: function() {
      if(this.inputdata === '') {
        return ""
      }else{
        return Base64.decode(this.inputdata)
      }
    },
    outMd5: function() {
      if(this.inputdata === '') {
        return ""
      }else{
        return this.$md5(this.inputdata)
      }
    }
  },
  methods: {
    funcjson () {
      this.inputdata= ''
      this.inputlabel= '请输入json'
      this.outputlabel= '格式化后的json'
      this.bjson = true
      this.eb64 = false
      this.db64 = false
    },
    encodeb64 () {
      this.inputdata= ''
      this.inputlabel= '请输入明文'
      this.outputlabel= 'Base64格式'
      this.bjson = false
      this.eb64 = true
      this.db64 = false
      this.emd5 = false
    },
    decodeb64 () {
      this.inputdata= ''
      this.inputlabel= '请输入Base64'
      this.outputlabel= '明文格式'
      this.bjson = false
      this.eb64 = false
      this.db64 = true
      this.emd5 = false
    },
    funcmd5 () {
      this.inputdata= ''
      this.inputlabel= '请输入明文'
      this.outputlabel= 'MD5:'
      this.bjson = false
      this.eb64 = false
      this.db64 = false
      this.emd5 = true
    }
  },
  components: {
    JsonViewer,
    BlogHeader,
    BlogFooter
  }
}
</script>

<style>
.demo-split{
  text-align: left;
  margin: 3%;
  margin-top:0px;
  min-height: 620px;
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
}
</style>


