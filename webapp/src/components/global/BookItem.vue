<template>
  <li class="post-item">
    <router-link :to=book.url>
      <span class="post-title">{{ book.title }}</span>
      <span class="post-day">({{book.author}}){{ book.day }}</span>
    </router-link>
  </li>
</template>

<script>
export default {
  props: {
    book: {
      type: Object,
      default: () => {}
    }
  }
}
</script>
