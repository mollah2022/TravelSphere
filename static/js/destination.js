(function () {
  'use strict';

  var btn          = document.getElementById('add-wishlist-btn');
  var feedbackDiv  = document.getElementById('wishlist-feedback');

  // Button বা feedback div না থাকলে (not logged in) চলবে না
  if (!btn || !feedbackDiv) return;

  // Country name JS variable থেকে নাও (template এ set করা)
  var countryName = (typeof COUNTRY_NAME !== 'undefined') ? COUNTRY_NAME : '';

  // ── Feedback দেখাও ──
  function showFeedback(msg, isSuccess) {
    feedbackDiv.innerHTML =
      '<span class="' +
      (isSuccess ? 'feedback-success' : 'feedback-error') +
      '">' + escapeHtml(msg) + '</span>';
  }

  // ── Add to Wishlist button click ──
  btn.addEventListener('click', function () {
    if (!countryName) {
      showFeedback('Country name missing.', false);
      return;
    }

    // Button disable করো — double submit prevent
    btn.disabled = true;
    btn.textContent = 'Adding...';
    feedbackDiv.innerHTML = '';

    var payload = JSON.stringify({
      country_name: countryName,
      status: 'Planned',
      note: ''
    });

    fetch('/api/wishlist', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: payload
    })
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (data.success) {
          showFeedback('✓ Added to wishlist!', true);
          btn.textContent = '✓ Added';
          // Button সবুজ করো
          btn.style.background = '#10b981';
        } else {
          showFeedback(data.error || 'Failed to add.', false);
          btn.disabled = false;
          btn.textContent = 'Add to Wishlist';
        }
      })
      .catch(function () {
        showFeedback('Network error. Please try again.', false);
        btn.disabled = false;
        btn.textContent = 'Add to Wishlist';
      });
  });

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