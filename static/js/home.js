(function () {
  'use strict';

  var searchInput = document.getElementById('home-search');
  var suggestionsBox = document.getElementById('search-suggestions');

  if (!searchInput || !suggestionsBox) return;

  var debounceTimer = null;

  // ── Debounce helper ──
  // User টাইপ করা থামলে 300ms পরে API call করে
  function debounce(fn, delay) {
    return function () {
      var args = arguments;
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(function () {
        fn.apply(null, args);
      }, delay);
    };
  }

  // ── Suggestions fetch ──
  function fetchSuggestions(query) {
    query = query.trim();

    // Empty হলে dropdown বন্ধ করো
    if (query.length === 0) {
      hideSuggestions();
      return;
    }

    fetch('/api/suggestions?q=' + encodeURIComponent(query))
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (!data.success || !data.data || data.data.length === 0) {
          hideSuggestions();
          return;
        }
        renderSuggestions(data.data);
      })
      .catch(function () {
        hideSuggestions();
      });
  }

  // ── Render suggestion items ──
  function renderSuggestions(countries) {
    var html = '';
    countries.forEach(function (c) {
      html += '<div class="suggestion-item" data-slug="' +
        escapeHtml(c.slug) + '">' +
        escapeHtml(c.name) +
        ' <span style="color:#9ca3af;font-size:0.8rem;">— ' +
        escapeHtml(c.capital) +
        '</span></div>';
    });
    suggestionsBox.innerHTML = html;
    suggestionsBox.classList.add('show');

    // প্রতিটা suggestion item এ click listener
    var items = suggestionsBox.querySelectorAll('.suggestion-item');
    items.forEach(function (item) {
      item.addEventListener('click', function () {
        var slug = item.getAttribute('data-slug');
        // Country detail page এ navigate করো (SSR)
        window.location.href = '/countries/' + slug;
      });
    });
  }

  // ── Hide suggestions ──
  function hideSuggestions() {
    suggestionsBox.classList.remove('show');
    suggestionsBox.innerHTML = '';
  }

  // ── Input event ──
  searchInput.addEventListener('input', debounce(function () {
    fetchSuggestions(searchInput.value);
  }, 300));

  // ── Outside click এ close ──
  document.addEventListener('click', function (e) {
    if (!searchInput.contains(e.target) &&
        !suggestionsBox.contains(e.target)) {
      hideSuggestions();
    }
  });

  // ── Escape key এ close ──
  searchInput.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') hideSuggestions();
  });

  // ── XSS prevent helper ──
  function escapeHtml(str) {
    if (!str) return '';
    return String(str)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }
})();