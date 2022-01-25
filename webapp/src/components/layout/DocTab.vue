<template>
  <div  class="col-lg-4 col-md-3 col-sm-4">
    <a :href="''+doc.url"><img :src="''+doc.pic" alt="" :title="doc.title" height="180"></a>
    <div class="art-info">
      <h4><a :href="''+doc.url">{{ doc.title }}</a></h4>
      <small>{{ doc.sum }}</small>
    </div>
    <div class="art-fields">
      <Icon type="md-pricetag" />
      &nbsp;<small>{{ doc.day }}</small>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    doc: {
      type: Object,
      default: () => {}
    }
  }
}
</script>

<style scoped src="../../assets/css/personal2.css"></style>
<style scoped src="../../assets/css/bootstrap.min.css"></style>