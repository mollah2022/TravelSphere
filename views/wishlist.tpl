<div class="container">
  <div class="page-header">
    <h1>Your Travel Wishlist</h1>
    <p>Manage your planned and visited destinations.</p>
  </div>

  <div id="wishlist-container">
    {{if .WishlistItems}}
    <div class="wishlist-grid" id="wishlist-grid">
      {{range .WishlistItems}}
      <div class="wishlist-card" id="item-{{.ID}}">
        <div class="wishlist-card-header">
          <span class="wishlist-country">{{.CountryName}}</span>
          <span class="badge badge-{{.Status}}">{{.Status}}</span>
        </div>
        {{if .Note}}
        <p class="wishlist-note">{{.Note}}</p>
        {{end}}
        <div class="wishlist-card-footer">
          <span class="wishlist-date">Added: {{.CreatedAt.Format "Jan 2, 2006"}}</span>
          <div class="wishlist-actions">
            <select
              class="status-select"
              onchange="updateStatus('{{.ID}}', this.value)"
            >
              <option value="Planned" {{if eq .Status "Planned"}}selected{{end}}>Planned</option>
              <option value="Visited" {{if eq .Status "Visited"}}selected{{end}}>Visited</option>
            </select>
            <button
              class="btn-delete"
              onclick="deleteItem('{{.ID}}')"
            >
              Remove
            </button>
          </div>
        </div>
      </div>
      {{end}}
    </div>
    {{else}}
    <div class="empty-state" id="empty-state">
      <p>Your wishlist is empty.</p>
      <a href="/countries" class="btn-primary">Explore Countries</a>
    </div>
    {{end}}
  </div>
</div>

<script src="/static/js/wishlist.js"></script>