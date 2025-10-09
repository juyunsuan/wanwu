/**
 * 通用 mixins 方法
 * 提供项目中常用的工具方法和生命周期钩子
 */
export default {
  data() {
    return {
    }
  },

  methods: {
    /**
     * 格式化日期
     * @param {Date|string|number} date - 日期
     * @param {string} format - 格式字符串
     * @returns {string} - 格式化后的日期字符串
     */
    $formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
      if (!date) return ''
      const d = new Date(date)
      if (isNaN(d.getTime())) return ''
      
      const year = d.getFullYear()
      const month = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      const hours = String(d.getHours()).padStart(2, '0')
      const minutes = String(d.getMinutes()).padStart(2, '0')
      const seconds = String(d.getSeconds()).padStart(2, '0')
      
      return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes)
        .replace('ss', seconds)
    },

    /**
     * 深拷贝对象
     * @param {any} obj - 要拷贝的对象
     * @returns {any} - 拷贝后的对象
     */
    $deepClone(obj) {
      if (obj === null || typeof obj !== 'object') return obj
      if (obj instanceof Date) return new Date(obj.getTime())
      if (obj instanceof Array) return obj.map(item => this.$deepClone(item))
      if (typeof obj === 'object') {
        const clonedObj = {}
        for (const key in obj) {
          if (obj.hasOwnProperty(key)) {
            clonedObj[key] = this.$deepClone(obj[key])
          }
        }
        return clonedObj
      }
    },

    /**
     * 防抖函数
     * @param {Function} func - 要防抖的函数
     * @param {number} delay - 延迟时间（毫秒）
     * @returns {Function} - 防抖后的函数
     */
    $debounce(func, delay = 300) {
      let timeoutId
      return function (...args) {
        clearTimeout(timeoutId)
        timeoutId = setTimeout(() => func.apply(this, args), delay)
      }
    },

    /**
     * 节流函数
     * @param {Function} func - 要节流的函数
     * @param {number} delay - 延迟时间（毫秒）
     * @returns {Function} - 节流后的函数
     */
    $throttle(func, delay = 300) {
      let lastCall = 0
      return function (...args) {
        const now = Date.now()
        if (now - lastCall >= delay) {
          lastCall = now
          func.apply(this, args)
        }
      }
    },
    /**
     * 获取文件大小格式化字符串
     * @param {number} bytes - 字节数
     * @returns {string} - 格式化后的文件大小
     */
    $formatFileSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    },

    /**
     * 验证邮箱格式
     * @param {string} email - 邮箱地址
     * @returns {boolean} - 是否有效
     */
    $isValidEmail(email) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return emailRegex.test(email)
    },

    /**
     * 验证手机号格式
     * @param {string} phone - 手机号
     * @returns {boolean} - 是否有效
     */
    $isValidPhone(phone) {
      const phoneRegex = /^1[3-9]\d{9}$/
      return phoneRegex.test(phone)
    },

    /**
     * 滚动到页面顶部
     */
    $scrollToTop() {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    },

    /**
     * 复制文本到剪贴板
     * @param {string} text - 要复制的文本
     * @returns {Promise} - 复制结果
     */
    async $copyToClipboard(text) {
      try {
        if (navigator.clipboard) {
          await navigator.clipboard.writeText(text)
          this.$success('复制成功')
        } else {
          // 兼容旧浏览器
          const textArea = document.createElement('textarea')
          textArea.value = text
          document.body.appendChild(textArea)
          textArea.select()
          document.execCommand('copy')
          document.body.removeChild(textArea)
          this.$success('复制成功')
        }
      } catch (error) {
        this.$error('复制失败')
        console.error('复制失败:', error)
      }
    },

    /**
     * 处理引用点击事件
     * @param {Event} e - 点击事件
     * @param {Object} options - 配置选项
     * @param {number} options.sessionStatus - 会话状态
     * @param {Object} options.sessionData - 会话数据
     * @param {string} options.citationSelector - 引用元素选择器，默认为 '.citation'
     * @param {string} options.subTagSelector - 子标签选择器，默认为 '.subTag'
     * @param {string} options.scrollElementId - 滚动容器ID，默认为 'timeScroll'
     * @param {Function} options.onToggleCollapse - 切换折叠状态的回调函数
     */
    $handleCitationClick(e, options = {}) {
      const {
        sessionStatus = 0,
        sessionData = null,
        citationSelector = '.citation',
        subTagSelector = '.subTag',
        scrollElementId = 'timeScroll',
        onToggleCollapse = null
      } = options

      // 检查会话状态
      if (sessionStatus === 0) return

      // 查找最近的引用元素
      const citationElement = e.target.closest(citationSelector)
      if (!citationElement) return

      // 获取标签索引
      const tagIndex = parseInt(citationElement.textContent, 10)
      if (isNaN(tagIndex) || tagIndex <= 0) return

      // 查找所有子标签
      const allSubTag = document.querySelectorAll(subTagSelector)
      if (allSubTag.length === 0) return

      // 检查索引是否有效
      if (tagIndex > allSubTag.length) return

      // 获取目标元素
      const targetElement = allSubTag[tagIndex - 1]
      if (!targetElement) return

      // 获取父级索引和折叠状态
      const parentsIndex = targetElement.dataset.parentsIndex
      const collapse = targetElement.dataset.collapse

      // 检查会话数据结构
      if (!sessionData || 
          !sessionData.history || 
          !sessionData.history[parentsIndex] || 
          !sessionData.history[parentsIndex].searchList || 
          !sessionData.history[parentsIndex].searchList[tagIndex - 1]) {
        return
      }

      // 切换折叠状态
      if (collapse === 'false') {
        if (onToggleCollapse && typeof onToggleCollapse === 'function') {
          // 使用自定义回调函数
          onToggleCollapse(sessionData.history[parentsIndex].searchList[tagIndex - 1], true)
        } else {
          // 默认使用 Vue.set 或直接赋值
          if (this.$set) {
            this.$set(
              sessionData.history[parentsIndex].searchList[tagIndex - 1],
              'collapse',
              true
            )
          } else {
            sessionData.history[parentsIndex].searchList[tagIndex - 1].collapse = true
          }
        }
      }

      // 滚动到底部
      const timeScrollElement = document.getElementById(scrollElementId)
      if (timeScrollElement) {
        timeScrollElement.scrollTop = timeScrollElement.scrollHeight
      }

      // 阻止事件冒泡
      e.stopPropagation()
    }
  },

  computed: {
    /**
     * 是否为空对象
     * @returns {Function} - 判断函数
     */
    $isEmpty() {
      return (obj) => {
        if (obj === null || obj === undefined) return true
        if (typeof obj === 'string') return obj.trim() === ''
        if (Array.isArray(obj)) return obj.length === 0
        if (typeof obj === 'object') return Object.keys(obj).length === 0
        return false
      }
    }
  },

  mounted() {

  },

  beforeDestroy() {

  }
}
