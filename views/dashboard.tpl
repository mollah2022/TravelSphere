<div class="container">
  <div class="page-header">
    <h1>Travel Dashboard</h1>
    <p>Your saved trips at a glance. Stats refresh automatically
      when your wishlist changes.</p>
  </div>

  <!-- Stats — AJAX এই div replace করে -->
  <div id="dashboard-stats">
    <div class="stats-grid">
      <div class="stat-card">
        <label>TOTAL SAVED</label>
        <div class="stat-number">{{.Summary.Total}}</div>
      </div>
      <div class="stat-card">
        <label>PLANNED</label>
        <div class="stat-number">{{.Summary.Planned}}</div>
      </div>
      <div class="stat-card">
        <label>VISITED</label>
        <div class="stat-number">{{.Summary.Visited}}</div>
      </div>
    </div>
  </div>

  <!-- Saved Destinations List -->
  <h2 class="section-title">Saved destinations</h2>
  {{if .Destinations}}
    {{range .Destinations}}
    <div class="destination-list-item">
      <strong>{{.CountryName}}</strong>
      —
      <span class="{{if eq .Status "Visited"}}status-visited{{else}}status-planned{{end}}">
        {{.Status}}
      </span>
      {{if .Note}}
        &middot; <span style="color:#6b7280;">{{.Note}}</span>
      {{end}}
    </div>
    {{end}}
  {{else}}
  <div class="empty-state">
    <p>No saved destinations yet.</p>
  </div>
  {{end}}
</div>

<script src="/static/js/dashboard.js"></script>