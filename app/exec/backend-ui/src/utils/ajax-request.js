import axios from 'axios';
import { MessageBox, Message } from 'element-ui'
//requestAPI 请求方法 （地址，参数，请求方法）
function requestAPI(url, params, method = 'GET') {

    console.log()

    //同步方法（成功返回，失败手动返回）
    return new Promise((resolve, reject) => {
        if (method === 'GET') {
            axios({
                url: window.location.origin + url,
                params,
                method,
                // `headers` 是即将被发送的自定义请求头
                headers: { 'content-type': 'application/json' },
            })
                .then(data => {
                    resolve(data.data);
                })
                .catch(err => {
                    console.log(err)
                    reject(err);
                });
        } else if (method === 'POST' || method === 'PUT' || method === 'DELETE') {
            axios({
                url: window.location.origin + url,
                data: params,
                method,
                // `headers` 是即将被发送的自定义请求头
                headers: { 'content-type': 'application/json' },
            })
                .then(data => {
                    resolve(data.data);
                })
                .catch(err => {


                    Message({
                        showClose: true,
                        message: err.response.data.message,
                        type: 'error'
                    });
                    reject(err);
                });
        }
    });
}
export default requestAPI;
