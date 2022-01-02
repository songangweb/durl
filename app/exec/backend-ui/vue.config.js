module.exports = {
    devServer: {
        proxy: {
            '/request': {
                target: 'http://backend.durl.fun',
                changeOrigin: true,
                pathRewrite: {
                    '^/request': '',
                },
            },
        },
    },
};
