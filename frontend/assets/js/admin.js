(function adminPage() {
  if (!window.cbxApi) return;

  var loadSystemOverview = document.getElementById('loadSystemOverview');
  var loadSpamSignals = document.getElementById('loadSpamSignals');
  var getAIOpsConfig = document.getElementById('getAIOpsConfig');
  var updateAIOpsForm = document.getElementById('updateAIOpsForm');
  var roleForm = document.getElementById('roleForm');

  var systemOverviewOutput = document.getElementById('systemOverviewOutput');
  var spamSignalsOutput = document.getElementById('spamSignalsOutput');
  var aiopsConfigOutput = document.getElementById('aiopsConfigOutput');
  var aiopsMsg = document.getElementById('aiopsMsg');
  var roleMsg = document.getElementById('roleMsg');

  if (loadSystemOverview) {
    loadSystemOverview.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/admin/system/overview' })
        .then(function (response) { systemOverviewOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { systemOverviewOutput.textContent = readError(error); });
    });
  }

  if (loadSpamSignals) {
    loadSpamSignals.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/admin/spam-vote/teachers' })
        .then(function (response) { spamSignalsOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { spamSignalsOutput.textContent = readError(error); });
    });
  }

  if (getAIOpsConfig) {
    getAIOpsConfig.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/admin/aiops/config' })
        .then(function (response) { aiopsConfigOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { aiopsConfigOutput.textContent = readError(error); });
    });
  }

  if (updateAIOpsForm) {
    updateAIOpsForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {};

      var keywordsCsv = updateAIOpsForm.keywordsCsv.value.trim();
      if (keywordsCsv) {
        payload.bannedKeywords = keywordsCsv.split(',').map(function (item) { return item.trim(); }).filter(Boolean);
      }

      var suggestWeight = Number(updateAIOpsForm.suggestWeight.value || 0);
      if (suggestWeight > 0) {
        payload.suggestWeight = suggestWeight;
      }

      window.cbxApi.request({ method: 'POST', url: '/admin/aiops/config', data: payload })
        .then(function (response) { aiopsMsg.textContent = JSON.stringify(response.data); })
        .catch(function (error) { aiopsMsg.textContent = readError(error); });
    });
  }

  if (roleForm) {
    roleForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {
        userId: Number(roleForm.userId.value),
        role: roleForm.role.value,
      };
      window.cbxApi.request({ method: 'POST', url: '/admin/access/role', data: payload })
        .then(function (response) { roleMsg.textContent = JSON.stringify(response.data); })
        .catch(function (error) { roleMsg.textContent = readError(error); });
    });
  }

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
