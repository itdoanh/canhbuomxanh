(function aiopsPage() {
  if (!window.cbxApi) return;

  var reloadLexicon = document.getElementById('reloadLexicon');
  var summaryLexicon = document.getElementById('summaryLexicon');
  var analyzeForm = document.getElementById('analyzeForm');
  var analyzeFlagForm = document.getElementById('analyzeFlagForm');
  var loadEnrichedQueue = document.getElementById('loadEnrichedQueue');

  var lexiconOutput = document.getElementById('lexiconOutput');
  var analyzeOutput = document.getElementById('analyzeOutput');
  var analyzeFlagOutput = document.getElementById('analyzeFlagOutput');
  var enrichedQueueOutput = document.getElementById('enrichedQueueOutput');
  var analyzeMsg = document.getElementById('analyzeMsg');
  var analyzeFlagMsg = document.getElementById('analyzeFlagMsg');

  if (reloadLexicon) {
    reloadLexicon.addEventListener('click', function () {
      window.cbxApi.request({ method: 'POST', url: '/aiops/lexicon/reload' })
        .then(function (response) { lexiconOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { lexiconOutput.textContent = readError(error); });
    });
  }

  if (summaryLexicon) {
    summaryLexicon.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/aiops/lexicon/summary' })
        .then(function (response) { lexiconOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { lexiconOutput.textContent = readError(error); });
    });
  }

  if (analyzeForm) {
    analyzeForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = { text: analyzeForm.text.value };
      window.cbxApi.request({ method: 'POST', url: '/aiops/analyze', data: payload })
        .then(function (response) {
          analyzeMsg.textContent = 'Analyze done';
          analyzeOutput.textContent = JSON.stringify(response.data, null, 2);
        })
        .catch(function (error) {
          analyzeMsg.textContent = readError(error);
        });
    });
  }

  if (analyzeFlagForm) {
    analyzeFlagForm.addEventListener('submit', function (event) {
      event.preventDefault();
      var payload = {
        sourceType: analyzeFlagForm.sourceType.value,
        sourceId: Number(analyzeFlagForm.sourceId.value),
        text: analyzeFlagForm.text.value,
      };
      window.cbxApi.request({ method: 'POST', url: '/aiops/analyze-and-flag', data: payload })
        .then(function (response) {
          analyzeFlagMsg.textContent = 'Analyze+flag done';
          analyzeFlagOutput.textContent = JSON.stringify(response.data, null, 2);
        })
        .catch(function (error) {
          analyzeFlagMsg.textContent = readError(error);
        });
    });
  }

  if (loadEnrichedQueue) {
    loadEnrichedQueue.addEventListener('click', function () {
      window.cbxApi.request({ method: 'GET', url: '/aiops/queue/enriched' })
        .then(function (response) { enrichedQueueOutput.textContent = JSON.stringify(response.data, null, 2); })
        .catch(function (error) { enrichedQueueOutput.textContent = readError(error); });
    });
  }

  function readError(error) {
    if (error && error.response && error.response.data && error.response.data.error) {
      return error.response.data.error;
    }
    return 'Yeu cau that bai';
  }
})();
