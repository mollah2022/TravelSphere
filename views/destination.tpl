<div class="container destination-page">

  <div class="destination-hero">
    <img
      src="{{.Country.FlagURL}}"
      alt="{{.Country.FlagAlt}}"
      class="destination-flag"
    >
    <div class="destination-hero-info">
      <h1>{{.Country.Name}}</h1>
      <p class="destination-official">{{.Country.OfficialName}}</p>
      <div class="destination-badges">
        <span class="badge">{{.Country.Region}}</span>
        <span class="badge">{{.Country.Subregion}}</span>
        <span class="badge">{{.Country.CCA2}}</span>
      </div>
    </div>
  </div>

  <div class="destination-grid">

    <div class="info-card">
      <h3>Basic Info</h3>
      <ul class="info-list">
        <li><span>Capital</span><strong>{{.Country.Capital}}</strong></li>
        <li><span>Population</span><strong>{{.FormattedPopulation}}</strong></li>
        <li><span>Languages</span><strong>{{.Country.Languages}}</strong></li>
        <li><span>Currencies</span><strong>{{.Country.Currencies}}</strong></li>
      </ul>
    </div>

    {{if .Weather.Available}}
    <div class="info-card weather-card">
      <h3>Current Weather in {{.Country.Capital}}</h3>
      <div class="weather-main">
        <img src="{{.Weather.Icon}}" alt="{{.Weather.Condition}}">
        <span class="weather-temp">{{.Weather.TempC}}°C</span>
      </div>
      <ul class="info-list">
        <li><span>Condition</span><strong>{{.Weather.Condition}}</strong></li>
        <li><span>Feels Like</span><strong>{{.Weather.FeelsLikeC}}°C</strong></li>
        <li><span>Humidity</span><strong>{{.Weather.Humidity}}%</strong></li>
        <li><span>Wind</span><strong>{{.Weather.WindKph}} km/h</strong></li>
      </ul>
    </div>
    {{end}}

  </div>

  {{if .IsLoggedIn}}
  <div class="wishlist-action">
    <h3>Add to Wishlist</h3>
    <div class="wishlist-form" id="wishlist-form">
      <input
        type="text"
        id="wishlist-note"
        placeholder="Add a note (optional)"
        maxlength="200"
      >
      <select id="wishlist-status">
        <option value="Planned">Planned</option>
        <option value="Visited">Visited</option>
      </select>
      <button
        class="btn-primary"
        onclick="addToWishlist('{{.Country.Name}}')"
      >
        + Add to Wishlist
      </button>
    </div>
    <div id="wishlist-msg"></div>
  </div>
  {{else}}
  <div class="wishlist-action">
    <p><a href="/login?redirect=/countries/{{.Country.Slug}}">Login</a> to save this destination.</p>
  </div>
  {{end}}

  {{if .Attractions}}
  <div class="attractions-section">
    <h3>Nearby Attractions</h3>
    <div class="attractions-grid">
      {{range .Attractions}}
      <div class="attraction-card">
        <div class="attraction-card-name">{{.Name}}</div>
        <div class="attraction-card-kinds">{{.Kinds}}</div>
        <div class="attraction-card-dist">{{.Distance}}m away</div>
      </div>
      {{end}}
    </div>
  </div>
  {{end}}

</div>

<script src="/static/js/destination.js"></script>