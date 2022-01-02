import axios from 'axios';
//requestAPI 请求方法 （地址，参数，请求方法）
function requestAPI(url, params, method = 'GET') {
    //同步方法（成功返回，失败手动返回）
    return new Promise((resolve, reject) => {
        if (method == 'GET') {
            axios({
                url,
                params,
                method,
                // `headers` 是即将被发送的自定义请求头
                headers: { 'content-type': 'application/json' },
            })
                .then(data => {
                    resolve(data.data);
                })
                .catch(err => {
                    reject(err);
                });
        } else if (method == 'POST' || method == 'PUT' || method == 'DELETE') {
            axios({
                url,
                data: params,
                method,
                // `headers` 是即将被发送的自定义请求头
                headers: { 'content-type': 'application/json' },
            })
                .then(data => {
                    resolve(data.data);
                })
                .catch(err => {
                    reject(err);
                });
        }
    });
}
export default requestAPI;
