<template>
  <div>
      <BlogHeader></BlogHeader>
      <h1>草稿箱</h1>
      <div class="post-page thin body">
      <Table border :columns="columns" :data="data"></Table>
      </div>
      <blog-footer></blog-footer>
  </div>
</template>

<script>
import NetWorking from '@/utils/networking'
import * as API from '@/utils/api'
import BlogHeader from '@/components/global/SiteHeader'
import BlogFooter from '@/components/global/SiteFooter'
export default {
  data () {
    return {
      columns: [
        {
          title: '文章标题',
          key: 'Title',
          render: (h, params) => {
            return h('div', [
              h('Icon', {
                props: {
                  type: 'person'
                }
              }),
              h('strong', params.row.Title+'|'+params.row.Code )
            ]);
          }
        },
        {
          title: '编辑时间',
          key: 'UpdateTime'
        },
        {
          title: '摘要',
          key: 'Sum'
        },
        {
          title: 'Action',
          key: 'action',
          width: 150,
          align: 'center',
          render: (h, params) => {
            return h('div', [
              h('Button', {
                props: {
                  type: 'primary',
                  size: 'small'
                },
                style: {
                  marginRight: '5px'
                },
                on: {
                  click: () => {
                    this.show(params.index)
                  }
                }
              }, 'View'),
              h('Button', {
                props: {
                  type: 'error',
                  size: 'small'
                },
                on: {
                  click: () => {
                    this.remove(params.index)
                }
              }
            }, 'Delete')
          ])
          }
        }
      ],
      data: []
    }
  },
  created () {
    NetWorking.doGet(API.blogcaches).then(response => {
      this.data = response.data
    }, (message) => {
      this.$Notice.error({
        title: '草稿箱信息获取失败',
        desc: 'Get Cache Box Failed!' + message
      })
    })
  },
  methods: {
    show (index) {
      this.$router.push({ path: '/editer/' + this.data[index].Code })
    },
    remove (index) {
      this.data.splice(index, 1)
    }
  },
  components: {
    BlogHeader,
    BlogFooter
  }
}
</script>
