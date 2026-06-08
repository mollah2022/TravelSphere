(function () {
  'use strict';

  var statsDiv = document.getElementById('dashboard-stats');
  if (!statsDiv) return;

  // ── Stats refresh function ──
  // Wishlist change হলে এই function call হয়
  function refreshStats() {
    fetch('/api/dashboard/summary')
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (!data.success || !data.data) return;
        renderStats(data.data);
      })
      .catch(function () {
        // Silent fail
      });
  }

  // ── Stats HTML render ──
  function renderStats(summary) {
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

  // ── Page load এ একবার refresh করো ──
  // Dashboard সবসময় latest stats দেখাবে
  refreshStats();

  // ── Global expose করো ──
  // wishlist.js থেকে call করার জন্য
  window.refreshDashboardStats = refreshStats;
})();