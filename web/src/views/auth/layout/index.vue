<template>
  <div class="auth">
    <div class="overview">
      <img :src="backgroundSrc" alt="">
    </div>
    <div class="auth-modal">
      <div class="header__left">
        <img
          v-if="commonInfo.login.logo && commonInfo.login.logo.path"
          style="height: 60px; margin: 0 15px 0 22px"
          :src="basePath + '/user/api' + commonInfo.login.logo.path"
          alt=""/>
        <!--<span style="font-size: 16px;">{{commonInfo.home.title || ''}}</span>-->
        <!--<div style="margin-left: 10px">
          <ChangeLang :isLogin="true" />
        </div>-->
      </div>
<!--      <div class="container__left">-->
<!--        {{ commonInfo.login.welcomeText }}-->
<!--      </div>-->

      <slot :commonInfo="commonInfo"/>

    </div>
  </div>
</template>

<script>
import {mapState, mapActions} from 'vuex'
import ChangeLang from "@/components/changeLang.vue"
import {replaceTitle, replaceIcon, avatarSrc} from "@/utils/util";
import { getCommonInfo } from '@/api/user'

export default {
  components: {ChangeLang},
  data() {
    return {
      backgroundSrc: require('@/assets/imgs/auth_bg.png'),
      basePath: this.$basePath
    }
  },
  computed: {
    ...mapState('login', ['commonInfo']),
    ...mapState('user', ['lang'])
  },
  watch: {
    'lang': {
      handler(val) {
        if (val) {
          /*this.getImgCode()
          this.getLogoInfo()*/
        }
      },
      immediate: true
    }
  },
  created() {
    this.getCommonInfo().then(() => {
      const { tab = {}, login = {} } = this.commonInfo || {}
      const { logo = {}, title = '' } = tab || {}
      const { background = {} } = login || {}

      background.path && this.setAuthBg(background.path)
      title && replaceTitle(title)
      logo.path && replaceIcon(logo.path)
      this.$emit('getCommonInfo', this.commonInfo)
    })
  },
  methods: {
    ...mapActions('login', ['getCommonInfo']),
    setDefaultImage() {
      this.backgroundSrc = require('@/assets/imgs/auth_bg.png')
    },
    setAuthBg(backgroundPath) {
      if (!backgroundPath) {
        this.setDefaultImage()
        return
      }
      this.backgroundSrc = avatarSrc(backgroundPath)
    },
  }
}
</script>

<style lang="scss" scoped>
.overview {
  position: relative;
  height: 100%;
  overflow: hidden;
  //background-color: #000;
  z-index: 10;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    background-size: 100% 100%;
  }

  .overview-desc {
    width: 800px;
    position: absolute;
    bottom: 56px;
    left: 56px;
    color: #fff;
    text-align: center;
    opacity: .8;
    letter-spacing: 1px;

    .desc {
      font-size: 30px;
      text-align: left;

      p:nth-child(1) {
        font-size: 22px;
      }

      p:nth-child(2) {
        font-size: 30px;
        margin: 10px 0;
      }

      p:nth-child(3) {
        font-size: 18px;
      }
    }

    .tabs {
      display: flex;
      font-size: 27px;
      margin-top: 30px;
      color: #fff;

      .tab {
        width: 1.63rem;
        min-width: 163px;
        margin-right: 20px;
        border: 1px solid #fff;
        cursor: pointer;

        p:nth-child(1) {
          font-size: 18px;
          padding: 4px 0 3px 0;
        }

        p:nth-child(2) {
          font-size: 12px;
          font-weight: 400;
          padding: 0 0 6px 0;
        }
      }
    }
  }
}

.auth {
  height: 100%;
}

.auth-modal {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 100%;
  z-index: 1000;

  .header__left {
    position: relative;
    width: 100%;
    min-width: 500px;
    color: #fff;
    font-weight: bold;
    display: flex;
    align-items: center;
    margin-top: 16px;
    margin-left: 10px;
    height: 60px;
  }

  .container__left {
    display: flex;
    align-items: center;
    height: calc(80% - 60px);
    font-size: 35px;
    width: calc(100% - 13% - 400px);
    justify-content: center;
    color: #fff;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.6);
  }

  .auth-box {
    position: absolute;
    width: 400px;
    min-width: 400px;
    height: 460px;
    top: 0;
    bottom: 0;
    right: 13%;
    margin: auto;
    background: rgba(244, 247, 255, 0.5);
    border-radius: 4px;
    //box-shadow: 0px 5px 30px 0px rgba(5,8,27,0.1);
    //border: 1px solid;
    //border-image: linear-gradient(180deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.8)) 1 1;
    backdrop-filter: blur(10px);

    .auth-header {
      color: $color_title;
      text-align: left;
      padding: 30px 30px 0 30px;
      //border-bottom: 1px solid rgba(235,235,235,0.8);
      span {
        display: inline-block;
        vertical-align: bottom;
        font-size: 24px;
      }
    }

    .auth-form {
      padding: 30px;

      /deep/ {
        .el-input {
          height: 30px;
          line-height: 36px;
        }

        .el-input__prefix {
          left: 14px;
        }

        .el-form-item__label {
          font-size: 12px !important;
          padding: 0 !important;
          margin-left: 0 !important;
          line-height: 26px !important;
          color: #425466;
        }

        .el-input__inner {
          font-size: 12px !important;
          background: #fff;
          border: none !important;
          padding-left: 50px;
        }

        i {
          color: #666 !important;
        }

        .el-icon-lock, .el-icon-user, .el-icon-key {
          font-size: 18px;
        }

        input[type=input]::placeholder {
          color: #B3B1BC;
        }

        input::-webkit-input-placeholder {
          color: #B3B1BC;
        }

        input::-moz-placeholder { /* Mozilla Firefox 19+ */
          color: #B3B1BC;
        }

        input:-moz-placeholder { /* Mozilla Firefox 4 to 18 */
          color: #B3B1BC;
        }

        input:-ms-input-placeholder { /* Internet Explorer 10-11 */
          color: #B3B1BC;
        }
      }

      .auth-form-item {
        position: relative;
      }

      /deep/ .auth-icon {
        position: absolute;
        width: 17px;
        z-index: 10;
        top: 13.5px;
        left: 17px;
      }

      /deep/ .pwd-icon {
        position: absolute;
        width: 17px;
        z-index: 10;
        top: 13.5px;
        right: 17px;
        cursor: pointer;
      }

      .nav-bt {
        width: 100%;
        height: 1px;
        line-height: 30px;
        font-size: 14px;
        margin-top: -20px;
      }

      .auth-bt {
        width: 100%;
        height: 40px;
        line-height: 30px;
        font-size: 14px;
        margin-top: 40px;

        .primary-bt {
          padding: 13px 0;
          height: 40px;
          //background: $color;
          border: none;
          font-size: 14px;
          border-radius: 3px !important;
        }

        .disabled.primary-bt {
          opacity: 0.5;
          cursor: not-allowed;
          border-color: #e9e9eb;
        }

        .operation {
          /*color: #fff;*/
          color: #8d8d8d;
          display: flex;
          margin: 30px 70px;

          span {
            display: block;
            flex: 1;
            text-align: center;
            font-size: 12px;
            cursor: pointer;
          }
        }
      }

      .bottom-text {
        text-align: center;
        font-weight: normal;
        color: #888;
        margin-top: 6px;
        font-size: 12px;
      }
    }

    .auth-pwd-input /deep/ {
      .el-input__inner {
        padding-right: 46px;
      }
    }
  }
}
</style>
