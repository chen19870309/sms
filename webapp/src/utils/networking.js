import axios from 'axios'
import storage from './storage'
import router from '../router'
/**
 * 使用cookie,带上cookie之后跨域就不行了
 */
// axios.defaults.withCredentials = true

/**
 * 拦截器，为请求头添加jwt信息
 */
let http = axios.create({
  baseURL: 'https://www.xiaoxibaby.xyz/blog',
  timeout: 10000
})

http.interceptors.request.use(
  config => {
    if (storage.get('JWT_TOKEN')) { // 判断是否存在token，如果存在的话，则每个http header都加上token
      config.headers.Authorization = `token ${storage.get('JWT_TOKEN')}`
    }
    return config
  },
  err => {
    return Promise.reject(err)
  }
)
http.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      console.log('axios:' + error.response.status)
      switch (error.response.status) {
        case 401:
          // 返回 401 清除token信息并跳转到登录页面
          storage.remove('JWT_TOKEN')
          router.replace({
            path: 'login',
            query: {redirect: router.currentRoute.fullPath}
          })
      }
    }
    return Promise.reject(error.response.data) // 返回接口返回的错误信息
  }
)

export default {
  httpMethod: {
    GET: 'get',
    POST: 'post',
    PUT: 'put',
    DELETE: 'delete'
  },
  instance: http,
  /**
   *
   * @param url
   * @param pathParams
   * @param options
   * @returns {*|Promise}
   */
  doGet(url, pathParams = null, options = null) {
    return this.doRequest(url, this.httpMethod.GET, pathParams, null, options)
  },

  /**
   * Http Post
   * @param url
   * @param pathParams
   * @param body
   * @param options
   * @returns {*|Promise}
   */
  doPost(url, pathParams = null, body = null, options = null) {
    return this.doRequest(url, this.httpMethod.POST, pathParams, body, options)
  },

  /**
   * Http Put
   * @param url
   * @param pathParams
   * @param body
   * @param options
   * @returns {*|Promise}
   */
  doPut(url, pathParams = null, body = null, options = null) {
    return this.doRequest(url, this.httpMethod.PUT, pathParams, body, options)
  },

  /**
   * Http Delete
   * @param url
   * @param pathParams
   * @param options
   * @returns {*|Promise}
   */
  doDelete(url, pathParams = null, options = null) {
    return this.doRequest(url, this.httpMethod.DELETE, pathParams, null, options)
  },

  /**
   * Http Request
   * @param url
   * @param method
   * @param pathParams
   * @param body
   * @param options
   * @returns {Promise}
   */
  doRequest(url, method, pathParams, body, options) {
    let wrapURL = this._wrapUrl(url, pathParams)
    let request = null
    switch (method) {
      case this.httpMethod.GET: {
        if (options !== null) {
          request = this.instance.get(wrapURL, {params: options})
        } else {
          request = this.instance.get(wrapURL)
        }
        break
      }
      case this.httpMethod.POST: {
        request = this.instance.post(wrapURL, body, options)
        break
      }
      case this.httpMethod.PUT: {
        request = this.instance.put(wrapURL, body, options)
        break
      }
      case this.httpMethod.DELETE: {
        request = this.instance.delete(wrapURL, options)
        break
      }
    }
    return this._requestPromise(request)
  },

  /**
   * Request Promise
   * @param request {Promise} Axios Http Promise
   * @returns {Promise}
   * @private
   */
  _requestPromise(request) {
    return new Promise(function (resolve, reject) {
      request.then(response => {
        if (response.data.success === true) {
          resolve(response.data)
        } else {
          reject(response.data.message)
        }
      }, error => {
        reject(error)
      })
    })
  },

  /**
   * 从URL匹配出参数pathParams
   * @param {string} url
   * @param {object} params
   * @returns {string}
   * @private
   */
  _wrapUrl(url, params) {
    if (params !== null) {
      let matches = this._getMatches(url)
      for (let match of matches) {
        let value = params[match.replace(':', '')]
        if (value !== null) {
          url = url.replace(match, value)
        }
      }
    }
    return url
  },

  /**
   * 正则Group匹配
   * @param {string} string
   * @returns {Array}
   * @private
   */
  _getMatches(string) {
    let matches = []
    let regex = /(:[a-z_A-Z]+)/
    let match = regex.exec(string)
    while (match !== null) {
      matches.push(match[1])
      string = string.replace(match[1], '')
      match = regex.exec(string)
    }
    return matches
  }
}
