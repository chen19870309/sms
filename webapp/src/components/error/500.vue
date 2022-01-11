<template>
  <error-content code="500" desc="Oh~~鬼知道服务器经历了什么~" :src="src"/>
</template>

<script>
import error500 from '@/assets/img/error-500.svg'
import errorContent from './ErrContent.vue'
export default {
  name: 'error_500',
  components: {
    errorContent
  },
  data () {
    return {
      src: error500
    }
  }
}
</script>
