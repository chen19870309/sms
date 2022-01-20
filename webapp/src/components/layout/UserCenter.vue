<template>
  <div class="user-center clearfix">
    <span class="headimg center-block"><img :src="''+user.Icon"></span>
    <div class="user-info clearfix">
        <h4>{{ user.Nickname }}</h4>
        <small>{{ user.Remark }}</small>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
    name: 'user-center',
    methods: {
      
    },
    computed: {
    ...mapGetters({
      user: 'currentUser'
    })
  }
}
</script>

<style scoped src="../../assets/css/personal2.css"></style>
<style scoped src="../../assets/css/bootstrap.min.css"></style>
