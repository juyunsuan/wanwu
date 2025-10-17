import service from "@/utils/request";
const hasLang = true
const BASE_URL = '/user/api/v1'

// 登录
export const login = (data) => {
    return service({
        url: `${BASE_URL}/base/login`,
        method: "post",
        data,
        hasLang
    });
};

// 2FA登录
// 第一级验证：密码
export const login2FA1 = (data) => {
    return service({
        url: `${BASE_URL}/base/login/email`,
        method: "post",
        data
    });
};
// 第二级验证：验证码
// 邮箱验证码
export const login2FA2Code = (data) => {
    return service({
        url: `${BASE_URL}/user/login/email/code`,
        method: "post",
        data,
    });
};
// 首次登录
export const login2FA2new = (data) => {
    return service({
        url: `${BASE_URL}/user/login`,
        method: "put",
        data
    });
}
// 非首次登录
export const login2FA2exist = (data) => {
    return service({
        url: `${BASE_URL}/user/login`,
        method: "post",
        data
    });
}

// 获取图形验证码
export const getImgVerCode = () => {
    return service({
        url: `${BASE_URL}/base/captcha`,
        method: "get",
        hasLang
    });
};

// 邮箱注册验证码发送
export const registerCode = (data) => {
    return service({
        url: `${BASE_URL}/base/register/email/code`,
        method: "post",
        data,
    });
};

// 用户邮箱注册
export const register = (data) => {
    return service({
        url: `${BASE_URL}/base/register/email`,
        method: "post",
        data,
    });
};

// 重置密码邮箱验证码发送
export const resetCode = (data) => {
    return service({
        url: `${BASE_URL}/base/password/email/code`,
        method: "post",
        data,
    });
};

// 重置密码
export const reset = (data) => {
    return service({
        url: `${BASE_URL}/base/password/email`,
        method: "post",
        data,
    });
};

export const getLangList = () => {
    return service({
        url: `${BASE_URL}/base/language/select`,
        method: "get",
    });
};


export const changeLang = (data) => {
    return service({
        url: `${BASE_URL}/user/language`,
        method: "put",
        data
    });
};

export const getUserDetail = (data) => {
    return service({
        url: `${BASE_URL}/user/info`,
        method: "get",
        params:data,
    });
};

export const getPermission = (data) => {
    return service({
        url: `${BASE_URL}/user/permission`,
        method: "get",
        params:data
    });
};

export const restUserPassword= (data) => {
    return service({
        url: `${BASE_URL}/user/admin/password`,
        method: "put",
        data,
    });
};
export const restOwnPassword= (data) => {
    return service({
        url: `${BASE_URL}/user/password`,
        method: "put",
        data,
    });
};

export const restAvatar= (data, config) => {
    return service({
        url: `${BASE_URL}/user/avatar`,
        method: "put",
        data,
        config
    });
};

export const docDownload = () => {
    return service({
        url: `${BASE_URL}/doc_center`,
        method: "get",
    });
};

// 公用上传 avatar
export const uploadAvatar = (data, config) => {
    return service({
        url: `${BASE_URL}/avatar`,
        method: "post",
        data,
        config
    })
}

// 平台信息
export const getCommonInfo= () => {
    return service({
        url: `${BASE_URL}/base/custom`,
        method: "get",
        hasLang
    });
}