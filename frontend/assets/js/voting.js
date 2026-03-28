(function votingPage() {
  if (!window.cbxApi) return;

  var voteForm = document.getElementById('voteForm');
  var alertForm = document.getElementById('alertForm');
  var releaseForm = document.getElementById('releaseForm');
  var recomputeForm = document.getElementById('recomputeForm');
  var summaryForm = document.getElementById('summaryForm');

  var voteMessage = document.getElementById('voteMessage');
  var alertMessage = document.getElementById('alertMessage');
  var releaseMessage = document.getElementById('releaseMessage');
  var recomputeMessage = document.getElementById('recomputeMessage');
  var summaryOutput = document.getElementById('summaryOutput');

  if (voteForm) {
    voteForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {
        teacherId: Number(voteForm.teacherId.value),
        rawScore: Number(voteForm.rawScore.value),
        semesterKey: voteForm.semesterKey.value,
        voteMode: voteForm.voteMode.value,
      };
      window.cbxApi
        .request({ method: 'POST', url: '/voting/votes', data: payload })
        .then(function (response) {
          voteMessage.textContent = JSON.stringify(response.data);
          voteForm.reset();
        })
        .catch(function (error) {
          voteMessage.textContent = readError(error);
        });
    });
  }

  if (alertForm) {
    alertForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {
        teacherId: Number(alertForm.teacherId.value),
        detail: alertForm.detail.value,
      };
      window.cbxApi
        .request({ method: 'POST', url: '/voting/alerts/forced-vote', data: payload })
        .then(function () {
          alertMessage.textContent = 'Alarm da duoc gui.';
          alertForm.reset();
        })
        .catch(function (error) {
          alertMessage.textContent = readError(error);
        });
    });
  }

  if (releaseForm) {
    releaseForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = { semesterKey: releaseForm.semesterKey.value };
      window.cbxApi
        .request({ method: 'POST', url: '/voting/votes/release', data: payload })
        .then(function (response) {
          releaseMessage.textContent = JSON.stringify(response.data);
        })
        .catch(function (error) {
          releaseMessage.textContent = readError(error);
        });
    });
  }

  if (recomputeForm) {
    recomputeForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {};
      if (recomputeForm.teacherId.value) {
        payload.teacherId = Number(recomputeForm.teacherId.value);
      }
      window.cbxApi
        .request({ method: 'POST', url: '/voting/badges/recompute', data: payload })
        .then(function (response) {
          recomputeMessage.textContent = JSON.stringify(response.data);
        })
        .catch(function (error) {
          recomputeMessage.textContent = readError(error);
        });
    });
  }

  if (summaryForm) {
    summaryForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var teacherId = Number(summaryForm.teacherId.value);
      window.cbxApi
        .request({ method: 'GET', url: '/voting/badges/teachers/' + String(teacherId) })
        .then(function (response) {
          summaryOutput.textContent = JSON.stringify(response.data, null, 2);
        })
        .catch(function (error) {
          summaryOutput.textContent = readError(error);
        });
    });
  }

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
