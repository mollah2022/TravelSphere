(function () {
  'use strict';

  var searchInput  = document.getElementById('country-search');
  var regionSelect = document.getElementById('region-filter');
  var resultsDiv   = document.getElementById('country-results');

  if (!searchInput || !regionSelect || !resultsDiv) return;

  var debounceTimer = null;

  // ── Debounce ──
  function debounce(fn, delay) {
    return function () {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(fn, delay);
    };
  }

  // ── Loading state ──
  function showLoading() {
    resultsDiv.innerHTML =
      '<div class="loading-text">' +
      '<span class="spinner"></span>' +
      '<p style="margin-top:0.75rem;">Loading countries...</p>' +
      '</div>';
  }

  // ── Country card HTML তৈরি ──
  // Server এর মতো একই card structure
  function buildCountryCard(c) {
    return '<a href="/countries/' + escapeHtml(c.slug) +
      '" class="country-card">' +
      '<img class="country-card-flag" src="' + escapeHtml(c.flag_url) +
      '" alt="Flag of ' + escapeHtml(c.name) + '" loading="lazy">' +
      '<div class="country-card-body">' +
      '<div class="country-card-name">' + escapeHtml(c.name) + '</div>' +
      '<div class="country-card-meta">' +
      'Capital: <span>' + escapeHtml(c.capital) + '</span><br>' +
      'Population: <span>' + escapeHtml(String(c.population)) + '</span><br>' +
      'Currency: <span>' + escapeHtml(c.currencies) + '</span><br>' +
      'Languages: <span>' + escapeHtml(c.languages) + '</span>' +
      '</div></div></a>';
  }

  // ── Fetch + render ──
  function fetchCountries() {
    var search = searchInput.value.trim();
    var region = regionSelect.value;

    // Query string তৈরি করো
    var params = new URLSearchParams();
    if (search) params.set('search', search);
    if (region && region !== 'all') params.set('region', region);

    showLoading();

    fetch('/api/countries?' + params.toString())
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (!data.success) {
          showError(data.error || 'Failed to load countries.');
          return;
        }

        var countries = data.data;
        if (!countries || countries.length === 0) {
          resultsDiv.innerHTML =
            '<div class="empty-state"><p>No countries found.</p></div>';
          return;
        }

        // Grid rebuild করো
        var html = '<div class="countries-grid">';
        countries.forEach(function (c) {
          html += buildCountryCard(c);
        });
        html += '</div>';
        resultsDiv.innerHTML = html;
      })
      .catch(function () {
        showError('Network error. Please try again.');
      });
  }

  // ── Error state ──
  function showError(msg) {
    resultsDiv.innerHTML =
      '<div class="empty-state">' +
      '<p style="color:#ef4444;">' + escapeHtml(msg) + '</p>' +
      '</div>';
  }

  // ── Event listeners ──
  searchInput.addEventListener('input', debounce(fetchCountries, 400));
  regionSelect.addEventListener('change', fetchCountries);

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