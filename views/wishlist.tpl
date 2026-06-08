<div class="container">
  <div class="page-header">
    <h1>Travel Wishlist</h1>
    <p>Edit notes, update trip status, or remove destinations.
      Changes save without reloading the page.</p>
  </div>

  <!-- Wishlist rows — AJAX এই div replace করে -->
  <div id="wishlist-rows">
    {{if .WishlistItems}}
    <div class="wishlist-table-wrap">
      <table class="wishlist-table">
        <thead>
          <tr>
            <th>COUNTRY</th>
            <th>NOTE</th>
            <th>STATUS</th>
            <th>ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          {{range .WishlistItems}}
          <tr id="row-{{.ID}}">
            <td><strong>{{.CountryName}}</strong></td>
            <td>
              <input
                type="text"
                class="note-input"
                value="{{.Note}}"
                placeholder="Add a note..."
                data-id="{{.ID}}"
              >
            </td>
            <td>
              <select class="status-select" data-id="{{.ID}}">
                <option value="Planned" {{if eq .Status "Planned"}}selected{{end}}>Planned</option>
                <option value="Visited" {{if eq .Status "Visited"}}selected{{end}}>Visited</option>
              </select>
            </td>
            <td>
              <button class="btn-save" onclick="saveItem('{{.ID}}')">Save</button>
              <button class="btn-delete" onclick="deleteItem('{{.ID}}')">Delete</button>
            </td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </div>
    {{else}}
    <div class="empty-state">
      <p>Your wishlist is empty. Browse
        <a href="/countries" style="color:#6c63ff;">countries</a>
        and add destinations!</p>
    </div>
    {{end}}
  </div>
</div>

<script src="/static/js/wishlist.js"></script>