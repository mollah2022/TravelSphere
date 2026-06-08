<div class="container">
  <div class="page-header">
    <h1>Travel Dashboard</h1>
    <p>Welcome back, <strong>{{.Username}}</strong>!</p>
  </div>

  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-number">{{.Summary.Total}}</div>
      <div class="stat-label">Total Saved</div>
    </div>
    <div class="stat-card planned">
      <div class="stat-number">{{.Summary.Planned}}</div>
      <div class="stat-label">Planned</div>
    </div>
    <div class="stat-card visited">
      <div class="stat-number">{{.Summary.Visited}}</div>
      <div class="stat-label">Visited</div>
    </div>
  </div>

  <h2 class="section-title">Saved Destinations</h2>

  {{if .Destinations}}
  <div class="dashboard-list">
    {{range .Destinations}}
    <div class="dashboard-item">
      <div class="dashboard-item-info">
        <span class="dashboard-country">{{.CountryName}}</span>
        {{if .Note}}
        <span class="dashboard-note">{{.Note}}</span>
        {{end}}
      </div>
      <div class="dashboard-item-meta">
        <span class="badge badge-{{.Status}}">{{.Status}}</span>
        <span class="dashboard-date">{{.CreatedAt.Format "Jan 2, 2006"}}</span>
      </div>
    </div>
    {{end}}
  </div>
  {{else}}
  <div class="empty-state">
    <p>No destinations saved yet.</p>
    <a href="/countries" class="btn-primary">Start Exploring</a>
  </div>
  {{end}}

</div>

<script src="/static/js/dashboard.js"></script>