<template>
  <div class="posts-group">
    <div class="posts-chepter" :id="chepter.id">{{ chepter.name }}</div>
      <ul class="posts-list">
        <book-item v-for="book in chepter.books" :key="book.id" :book="book"></book-item>
      </ul>
  </div>
</template>

<script>
import BookItem from '@/components/global/BookItem'
export default {
  props: {
    chepter: {
      type: Object,
      default: function () {
        return {
          id: 0,
          name: '',
          books: []
        }
      }
    }
  },
  components: {
    BookItem
  }
}
</script>
