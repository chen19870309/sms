<template>
  <div class="error-page">
    <div class="content-con">
      <img :src="src" :alt="code">
      <div class="text-con">
        <h4>{{ code }}</h4>
        <h5>{{ desc }}</h5>
      </div>
      <err-btn-group class="back-btn-group"></err-btn-group>
    </div>
  </div>
</template>

<script>
import './error.less'
import errBtnGroup from './ErrBtnGroup.vue'
export default {
  name: 'error_content',
  components: {
    errBtnGroup
  },
  props: {
    code: String,
    desc: String,
    src: String
  }
}
</script>