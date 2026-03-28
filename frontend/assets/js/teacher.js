(function teacherPage() {
  var teacherList = document.getElementById('teacherList');
  var teacherDetail = document.getElementById('teacherDetail');
  var teacherIdInput = document.getElementById('teacherIdInput');
  var teacherLookupBtn = document.getElementById('teacherLookupBtn');

  if (!teacherList || !teacherDetail || !window.cbxApi) return;

  function loadTeachers() {
    window.cbxApi
      .request({ method: 'GET', url: '/teachers/' })
      .then(function (response) {
        var items = (response.data && response.data.items) || [];
        if (!items.length) {
          teacherList.innerHTML = '<article class="glass-card teacher-card">Chua co teacher public nao.</article>';
          return;
        }

        teacherList.innerHTML = items
          .map(function (item) {
            return (
              '<article class="glass-card teacher-card">' +
              '<h3>' + escapeHtml(item.displayName || '') + '</h3>' +
              '<p>' + escapeHtml(item.schoolName || '') + '</p>' +
              '<p>' + escapeHtml(item.subjectName || '') + '</p>' +
              '<p>ID: ' + String(item.id) + '</p>' +
              '</article>'
            );
          })
          .join('');
      })
      .catch(function (error) {
        teacherList.innerHTML = '<article class="glass-card teacher-card">' + escapeHtml(readError(error)) + '</article>';
      });
  }

  function lookupTeacher() {
    var id = Number(teacherIdInput.value || 0);
    if (!id) {
      teacherDetail.innerHTML = 'Nhap teacher id hop le.';
      return;
    }

    window.cbxApi
      .request({ method: 'GET', url: '/teachers/' + String(id) })
      .then(function (response) {
        var item = response.data || {};
        teacherDetail.innerHTML =
          '<h3>' +
          escapeHtml(item.displayName || '') +
          '</h3>' +
          '<p>School: ' +
          escapeHtml(item.schoolName || '-') +
          '</p>' +
          '<p>Subject: ' +
          escapeHtml(item.subjectName || '-') +
          '</p>' +
          '<p>Claim: ' +
          escapeHtml(item.claimStatus || '-') +
          '</p>';
      })
      .catch(function (error) {
        teacherDetail.innerHTML = escapeHtml(readError(error));
      });
  }

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

  teacherLookupBtn.addEventListener('click', lookupTeacher);
  loadTeachers();
})();
