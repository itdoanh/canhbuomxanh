(function authPage() {
  var registerForm = document.getElementById('registerForm');
  var loginForm = document.getElementById('loginForm');
  var message = document.getElementById('authMessage');

  if (!registerForm || !loginForm || !window.cbxApi) return;

  registerForm.addEventListener('submit', function (event) {
    event.preventDefault();
    var payload = {
      fullName: registerForm.fullName.value,
      email: registerForm.email.value,
      password: registerForm.password.value,
      role: registerForm.role.value,
    };

    window.cbxApi
      .request({ method: 'POST', url: '/auth/register', data: payload })
      .then(function () {
        message.textContent = 'Dang ky thanh cong. Ban co the dang nhap ngay.';
        registerForm.reset();
      })
      .catch(function (error) {
        message.textContent = readError(error);
      });
  });

  loginForm.addEventListener('submit', function (event) {
    event.preventDefault();
    var payload = {
      email: loginForm.email.value,
      password: loginForm.password.value,
    };

    window.cbxApi
      .request({ method: 'POST', url: '/auth/login', data: payload })
      .then(function (response) {
        window.cbxApi.setToken(response.data.accessToken || '');
        message.textContent = 'Dang nhap thanh cong. Token da duoc luu.';
        loginForm.reset();
      })
      .catch(function (error) {
        message.textContent = readError(error);
      });
  });

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
