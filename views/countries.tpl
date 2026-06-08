<div class="container">
  <div class="explorer-header">
    <h1>Country Explorer</h1>
    <p>Browse every destination on first load. Search and filter update only the results
      below — no full page reload.</p>
  </div>

  <!-- Search + Filter Bar -->
  <div class="search-filter-bar">
    <div class="filter-group">
      <label>SEARCH</label>
      <input
        type="text"
        id="country-search"
        placeholder="Country or capital..."
        autocomplete="off"
      >
    </div>
    <div class="filter-group">
      <label>REGION</label>
      <select id="region-filter">
        {{range .Regions}}
        <option value="{{if eq . "All Regions"}}all{{else}}{{.}}{{end}}">{{.}}</option>
        {{end}}
      </select>
    </div>
  </div>

  <!-- Country Results — AJAX এই div replace করে -->
  <div id="country-results">
    {{if .Countries}}
    <div class="countries-grid">
      {{range .Countries}}
      <a href="/countries/{{.Slug}}" class="country-card">
        <img
          class="country-card-flag"
          src="{{.FlagURL}}"
          alt="Flag of {{.Name}}"
          loading="lazy"
        >
        <div class="country-card-body">
          <div class="country-card-name">{{.Name}}</div>
          <div class="country-card-meta">
            Capital: <span>{{.Capital}}</span><br>
            Population: <span>{{.Population}}</span><br>
            Currency: <span>{{.Currencies}}</span><br>
            Languages: <span>{{.Languages}}</span>
          </div>
        </div>
      </a>
      {{end}}
    </div>
    {{else}}
    <div class="empty-state"><p>No countries found.</p></div>
    {{end}}
  </div>
</div>

<script src="/static/js/countries.js"></script>