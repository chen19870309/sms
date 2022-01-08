<template>
  <li class="post-item">
    <a :href="book.url">
      <span class="post-title">{{ book.title }}</span>
      <span class="post-day">{{ book.day }}</span>
    </a>
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
