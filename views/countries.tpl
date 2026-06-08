<div class="container">
  <div class="page-header">
    <h1>Country Explorer</h1>
    <p>{{.TotalCount}} countries loaded</p>
  </div>

  {{if .Error}}
  <div class="alert-error">{{.Error}}</div>
  {{end}}

  <div class="filter-bar">
    <input
      type="text"
      id="country-search"
      placeholder="Search country or capital..."
      autocomplete="off"
    >
    <select id="region-filter">
      {{range .Regions}}
      <option value="{{.}}">{{.}}</option>
      {{end}}
    </select>
  </div>

  <div class="countries-grid" id="countries-grid">
    {{range .Countries}}
    <a href="/countries/{{.Slug}}" class="country-card">
      <img
        src="{{.FlagURL}}"
        alt="{{.FlagAlt}}"
        loading="lazy"
      >
      <div class="country-card-body">
        <div class="country-card-name">{{.Name}}</div>
        <div class="country-card-meta">
          <span>{{.Capital}}</span>
          <span class="badge">{{.Region}}</span>
        </div>
      </div>
    </a>
    {{end}}
  </div>

  <div class="empty-state" id="no-results" style="display:none;">
    <p>No countries match your search.</p>
  </div>
</div>

<script src="/static/js/countries.js"></script>