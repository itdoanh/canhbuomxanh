(function moderatorPage() {
  if (!window.cbxApi) return;

  var loadAIQueue = document.getElementById('loadAIQueue');
  var loadAppeals = document.getElementById('loadAppeals');
  var reviewAppealForm = document.getElementById('reviewAppealForm');
  var violationForm = document.getElementById('violationForm');

  var aiQueueOutput = document.getElementById('aiQueueOutput');
  var appealQueueOutput = document.getElementById('appealQueueOutput');
  var reviewAppealMsg = document.getElementById('reviewAppealMsg');
  var violationMsg = document.getElementById('violationMsg');

  if (loadAIQueue) {
    loadAIQueue.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/aiops/queue/enriched' })
        .then(function (response) { aiQueueOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { aiQueueOutput.textContent = readError(error); });
    });
  }

  if (loadAppeals) {
    loadAppeals.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/moderator/appeals' })
        .then(function (response) { appealQueueOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { appealQueueOutput.textContent = readError(error); });
    });
  }

  if (reviewAppealForm) {
    reviewAppealForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var appealId = Number(reviewAppealForm.appealId.value);
      var payload = { action: reviewAppealForm.action.value };
      window.cbxApi.request({ method: 'POST', url: '/moderator/appeals/' + String(appealId) + '/review', data: payload })
        .then(function (response) { reviewAppealMsg.textContent = JSON.stringify(response.data); })
        .catch(function (error) { reviewAppealMsg.textContent = readError(error); });
    });
  }

  if (violationForm) {
    violationForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var userId = Number(violationForm.userId.value);
      var payload = { action: violationForm.action.value };
      window.cbxApi.request({ method: 'POST', url: '/moderator/violations/user/' + String(userId), data: payload })
        .then(function (response) { violationMsg.textContent = JSON.stringify(response.data); })
        .catch(function (error) { violationMsg.textContent = readError(error); });
    });
  }

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
