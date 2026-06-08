(function () {
  'use strict';

  var wishlistRows = document.getElementById('wishlist-rows');
  if (!wishlistRows) return;

  // ── Row HTML তৈরি ──
  function buildRow(item) {
    var isVisited = item.status === 'Visited';
    return '<tr id="row-' + item.id + '">' +
      '<td><strong>' + escapeHtml(item.country_name) + '</strong></td>' +
      '<td>' +
        '<input type="text" class="note-input" ' +
        'value="' + escapeHtml(item.note || '') + '" ' +
        'placeholder="Add a note..." ' +
        'data-id="' + item.id + '">' +
      '</td>' +
      '<td>' +
        '<select class="status-select" data-id="' + item.id + '">' +
          '<option value="Planned"' + (!isVisited ? ' selected' : '') + '>Planned</option>' +
          '<option value="Visited"' + (isVisited ? ' selected' : '') + '>Visited</option>' +
        '</select>' +
      '</td>' +
      '<td>' +
        '<button class="btn-save" onclick="saveItem(\'' + item.id + '\')">Save</button>' +
        '<button class="btn-delete" onclick="deleteItem(\'' + item.id + '\')">Delete</button>' +
      '</td>' +
      '</tr>';
  }

  // ── Table rebuild করো ──
  function rebuildTable(items) {
    if (!items || items.length === 0) {
      wishlistRows.innerHTML =
        '<div class="empty-state">' +
        '<p>Your wishlist is empty. Browse ' +
        '<a href="/countries" style="color:#6c63ff;">countries</a>' +
        ' and add destinations!</p></div>';
      return;
    }

    var html = '<div class="wishlist-table-wrap">' +
      '<table class="wishlist-table"><thead><tr>' +
      '<th>COUNTRY</th><th>NOTE</th><th>STATUS</th><th>ACTIONS</th>' +
      '</tr></thead><tbody>';

    items.forEach(function (item) {
      html += buildRow(item);
    });

    html += '</tbody></table></div>';
    wishlistRows.innerHTML = html;
  }

  // ── Save item (PUT) ──
  window.saveItem = function (id) {
    var row = document.getElementById('row-' + id);
    if (!row) return;

    var noteInput   = row.querySelector('.note-input');
    var statusSelect = row.querySelector('.status-select');

    var note   = noteInput ? noteInput.value.trim() : '';
    var status = statusSelect ? statusSelect.value : 'Planned';

    var saveBtn = row.querySelector('.btn-save');
    if (saveBtn) {
      saveBtn.textContent = 'Saving...';
      saveBtn.disabled = true;
    }

    fetch('/api/wishlist/' + id, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ note: note, status: status })
    })
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (data.success) {
          // Save success — visual feedback
          if (saveBtn) {
            saveBtn.textContent = '✓ Saved';
            saveBtn.style.background = '#10b981';
            setTimeout(function () {
              saveBtn.textContent = 'Save';
              saveBtn.style.background = '';
              saveBtn.disabled = false;
            }, 1500);
          }
          // Dashboard stats refresh করো
          refreshDashboardStats();
        } else {
          alert(data.error || 'Failed to save.');
          if (saveBtn) {
            saveBtn.textContent = 'Save';
            saveBtn.disabled = false;
          }
        }
      })
      .catch(function () {
        alert('Network error. Please try again.');
        if (saveBtn) {
          saveBtn.textContent = 'Save';
          saveBtn.disabled = false;
        }
      });
  };

  // ── Delete item ──
  window.deleteItem = function (id) {
    if (!confirm('Remove this destination from your wishlist?')) return;

    fetch('/api/wishlist/' + id, { method: 'DELETE' })
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (data.success) {
          // Row remove করো
          var row = document.getElementById('row-' + id);
          if (row) row.remove();

          // Table empty হলে empty state দেখাও
          var tbody = wishlistRows.querySelector('tbody');
          if (tbody && tbody.children.length === 0) {
            wishlistRows.innerHTML =
              '<div class="empty-state">' +
              '<p>Your wishlist is empty. Browse ' +
              '<a href="/countries" style="color:#6c63ff;">countries</a>' +
              ' and add destinations!</p></div>';
          }

          // Dashboard stats refresh করো
          refreshDashboardStats();
        } else {
          alert(data.error || 'Failed to delete.');
        }
      })
      .catch(function () {
        alert('Network error. Please try again.');
      });
  };

  // ── Dashboard stats refresh helper ──
  // Wishlist page থেকে dashboard stats update
  function refreshDashboardStats() {
    fetch('/api/dashboard/summary')
      .then(function (res) { return res.json(); })
      .then(function (data) {
        // Dashboard page এ থাকলে stats update করো
        if (data.success && data.data) {
          updateStatCards(data.data);
        }
      })
      .catch(function () {
        // Silent fail — stats refresh optional
      });
  }

  // ── Stat cards update ──
  function updateStatCards(summary) {
    var statsDiv = document.getElementById('dashboard-stats');
    if (!statsDiv) return;

    statsDiv.innerHTML =
      '<div class="stats-grid">' +
        '<div class="stat-card">' +
          '<label>TOTAL SAVED</label>' +
          '<div class="stat-number">' + (summary.total || 0) + '</div>' +
        '</div>' +
        '<div class="stat-card">' +
          '<label>PLANNED</label>' +
          '<div class="stat-number">' + (summary.planned || 0) + '</div>' +
        '</div>' +
        '<div class="stat-card">' +
          '<label>VISITED</label>' +
          '<div class="stat-number">' + (summary.visited || 0) + '</div>' +
        '</div>' +
      '</div>';
  }

  // ── XSS prevent ──
  function escapeHtml(str) {
    if (!str) return '';
    return String(str)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }
})();