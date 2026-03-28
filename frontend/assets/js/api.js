(function apiClientSetup() {
  var baseURL = localStorage.getItem('cbx_api_base') || 'http://localhost:8080/api/v1';

  window.cbxApi = {
    baseURL: baseURL,
    tokenKey: 'cbx_access_token',
    setBaseURL: function (nextURL) {
      localStorage.setItem('cbx_api_base', nextURL);
      this.baseURL = nextURL;
    },
    getToken: function () {
      return localStorage.getItem(this.tokenKey) || '';
    },
    setToken: function (token) {
      localStorage.setItem(this.tokenKey, token);
    },
    clearToken: function () {
      localStorage.removeItem(this.tokenKey);
    },
    request: function (config) {
      var finalConfig = Object.assign({}, config);
      finalConfig.url = this.baseURL + config.url;
      finalConfig.headers = Object.assign({}, config.headers || {});

      var token = this.getToken();
      if (token) {
        finalConfig.headers.Authorization = 'Bearer ' + token;
      }

      return window.axios(finalConfig);
    },
  };
})();
