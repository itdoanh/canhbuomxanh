(function teacherPortalPage() {
  if (!window.cbxApi) return;

  var claimForm = document.getElementById('claimForm');
  var optOutForm = document.getElementById('optOutForm');
  var appealForm = document.getElementById('appealForm');
  var loadDashBtn = document.getElementById('loadTeacherDash');
  var loadAppealsBtn = document.getElementById('loadMyAppeals');

  var claimMsg = document.getElementById('claimMsg');
  var optOutMsg = document.getElementById('optOutMsg');
  var appealMsg = document.getElementById('appealMsg');
  var dashOutput = document.getElementById('teacherDashOutput');
  var appealsOutput = document.getElementById('teacherAppealsOutput');

  if (claimForm) {
    claimForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = { teacherId: Number(claimForm.teacherId.value) };
      window.cbxApi.request({ method: 'POST', url: '/teacher-portal/claim', data: payload })
        .then(function (response) { claimMsg.textContent = stringify(response.data); })
        .catch(function (error) { claimMsg.textContent = readError(error); });
    });
  }

  if (optOutForm) {
    optOutForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = { teacherId: Number(optOutForm.teacherId.value) };
      window.cbxApi.request({ method: 'POST', url: '/teacher-portal/opt-out', data: payload })
        .then(function (response) { optOutMsg.textContent = stringify(response.data); })
        .catch(function (error) { optOutMsg.textContent = readError(error); });
    });
  }

  if (appealForm) {
    appealForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {
        teacherId: Number(appealForm.teacherId.value),
        targetType: appealForm.targetType.value,
        targetId: Number(appealForm.targetId.value),
        reason: appealForm.reason.value,
      };
      window.cbxApi.request({ method: 'POST', url: '/teacher-portal/appeals', data: payload })
        .then(function (response) { appealMsg.textContent = stringify(response.data); })
        .catch(function (error) { appealMsg.textContent = readError(error); });
    });
  }

  if (loadDashBtn) {
    loadDashBtn.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/teacher-portal/dashboard' })
        .then(function (response) { dashOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { dashOutput.textContent = readError(error); });
    });
  }

  if (loadAppealsBtn) {
    loadAppealsBtn.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/teacher-portal/appeals' })
        .then(function (response) { appealsOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { appealsOutput.textContent = readError(error); });
    });
  }

  function stringify(data) {
    return JSON.stringify(data);
  }

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
