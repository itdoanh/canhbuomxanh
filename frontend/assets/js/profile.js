(function profilePage() {
  var profileCard = document.getElementById('profileCard');
  var profileForm = document.getElementById('profileForm');
  var profileMessage = document.getElementById('profileMessage');

  if (!profileCard || !profileForm || !window.cbxApi) return;

  function loadMe() {
    window.cbxApi
      .request({ method: 'GET', url: '/profile/me' })
      .then(function (response) {
        var me = response.data || {};
        profileCard.innerHTML =
          '<h2>' +
          escapeHtml(me.fullName || '-') +
          '</h2>' +
          '<p>Email: ' +
          escapeHtml(me.email || '-') +
          '</p>' +
          '<p>Role: ' +
          escapeHtml(me.role || '-') +
          '</p>' +
          '<p>Avatar: ' +
          escapeHtml(me.avatarUrl || '-') +
          '</p>';

        profileForm.fullName.value = me.fullName || '';
        profileForm.avatarUrl.value = me.avatarUrl || '';
      })
      .catch(function (error) {
        profileCard.innerHTML = '<h2>Profile chua san sang</h2><p>' + escapeHtml(readError(error)) + '</p>';
      });
  }

  profileForm.addEventListener('submit', function (event) {
    event.preventDefault();
    var payload = {
      fullName: profileForm.fullName.value,
      avatarUrl: profileForm.avatarUrl.value,
    };

    window.cbxApi
      .request({ method: 'PUT', url: '/profile/me', data: payload })
      .then(function () {
        profileMessage.textContent = 'Cap nhat thanh cong.';
        loadMe();
      })
      .catch(function (error) {
        profileMessage.textContent = readError(error);
      });
  });

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }

  function escapeHtml(text) {
    return String(text)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  loadMe();
})();
