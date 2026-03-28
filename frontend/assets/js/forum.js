(function forumPage() {
  var postForm = document.getElementById('postForm');
  var postMessage = document.getElementById('postMessage');
  var postsContainer = document.getElementById('postsContainer');

  if (!postsContainer || !window.cbxApi) return;

  if (postForm) {
    postForm.addEventListener('submit', function (event) {
      event.preventDefault();

      var payload = {
        title: postForm.title.value,
        body: postForm.body.value,
      };

      if (postForm.teacherId.value) {
        payload.teacherId = Number(postForm.teacherId.value);
      }

      window.cbxApi
        .request({ method: 'POST', url: '/forum/posts', data: payload })
        .then(function () {
          postMessage.textContent = 'Da dang bai viet.';
          postForm.reset();
          loadPosts();
        })
        .catch(function (error) {
          postMessage.textContent = readError(error);
        });
    });
  }

  function loadPosts() {
    window.cbxApi
      .request({ method: 'GET', url: '/forum/posts' })
      .then(function (response) {
        var items = (response.data && response.data.items) || [];
        renderPosts(items);
      })
      .catch(function (error) {
        postsContainer.innerHTML = '<article class="glass-card post-card">' + escapeHtml(readError(error)) + '</article>';
      });
  }

  function renderPosts(items) {
    if (!items.length) {
      postsContainer.innerHTML = '<article class="glass-card post-card">Chua co bai viet nao.</article>';
      return;
    }

    postsContainer.innerHTML = items
      .map(function (item) {
        return (
          '<article class="glass-card post-card">' +
          '<h3>' + escapeHtml(item.title || '') + '</h3>' +
          '<p>' + escapeHtml(item.body || '') + '</p>' +
          '<p class="post-meta">Tac gia: ' + escapeHtml(item.authorName || 'unknown') + '</p>' +
          '<button class="btn btn-ghost" data-post-id="' + String(item.id) + '">Tai comments</button>' +
          '<div id="comments-' + String(item.id) + '" class="comment-list"></div>' +
          '</article>'
        );
      })
      .join('');

    postsContainer.querySelectorAll('button[data-post-id]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var postId = btn.getAttribute('data-post-id');
        loadComments(postId);
      });
    });
  }

  function loadComments(postId) {
    var target = document.getElementById('comments-' + String(postId));
    if (!target) return;

    window.cbxApi
      .request({ method: 'GET', url: '/forum/posts/' + String(postId) + '/comments' })
      .then(function (response) {
        var items = (response.data && response.data.items) || [];
        if (!items.length) {
          target.innerHTML = '<div class="comment-item">Chua co binh luan.</div>';
          return;
        }
        target.innerHTML = items
          .map(function (item) {
            return '<div class="comment-item"><strong>' + escapeHtml(item.authorName || 'unknown') + ':</strong> ' + escapeHtml(item.body || '') + '</div>';
          })
          .join('');
      })
      .catch(function (error) {
        target.innerHTML = '<div class="comment-item">' + escapeHtml(readError(error)) + '</div>';
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

  loadPosts();
})();
