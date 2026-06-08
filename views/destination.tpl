<div class="container">

  <!-- Country Hero Card -->
  <div class="destination-hero">
    <div class="destination-hero-inner">
      <img
        class="destination-flag"
        src="{{.Country.FlagURL}}"
        alt="{{.Country.FlagAlt}}"
      >
      <div class="destination-info">
        <div class="region-badge">{{.Country.Region}}</div>
        <h1>{{.Country.Name}}</h1>
        <p class="official-name">{{.Country.OfficialName}}</p>

        <div class="destination-meta-grid">
          <div class="meta-item">
            <label>CAPITAL</label>
            <span>{{.Country.Capital}}</span>
          </div>
          <div class="meta-item">
            <label>POPULATION</label>
            <span>{{.FormattedPopulation}}</span>
          </div>
          <div class="meta-item">
            <label>REGION</label>
            <span>{{.Country.Region}}<br>{{.Country.Subregion}}</span>
          </div>
          <div class="meta-item">
            <label>CURRENCY</label>
            <span>{{.Country.Currencies}}</span>
          </div>
          <div class="meta-item">
            <label>LANGUAGES</label>
            <span>{{.Country.Languages}}</span>
          </div>
        </div>

        {{if .IsLoggedIn}}
        <button
          class="btn-add-wishlist"
          id="add-wishlist-btn"
          data-country="{{.Country.Name}}"
        >
          Add to Wishlist
        </button>
        {{else}}
        <a href="/login?redirect=/countries/{{.Country.Slug}}"
           class="btn-add-wishlist"
           style="display:inline-block;text-align:center;">
          Login to Add Wishlist
        </a>
        {{end}}

        <!-- AJAX feedback আসবে এখানে -->
        <div id="wishlist-feedback"></div>
      </div>
    </div>
  </div>

  <!-- Weather + Attractions -->
  <div class="two-col">

    <!-- Travel Weather -->
    <div class="weather-card">
      <h3>Travel weather</h3>
      {{if .Weather.Available}}
      <div class="weather-info">
        <div class="weather-temp">{{.Weather.TempC}}°C</div>
        <div class="weather-condition">
          {{if .Weather.Icon}}
          <img src="{{.Weather.Icon}}" alt="weather" style="width:32px;vertical-align:middle;">
          {{end}}
          {{.Weather.Condition}}
        </div>
        <div class="weather-detail">
          Feels like {{.Weather.FeelsLikeC}}°C &nbsp;|&nbsp;
          Humidity {{.Weather.Humidity}}% &nbsp;|&nbsp;
          Wind {{.Weather.WindKph}} kph
        </div>
      </div>
      {{else}}
      <div class="weather-unavailable">
        Weather data is optional. Add <code>WEATHER_API_KEY</code>
        to your <code>.env</code> file to enable live conditions.
      </div>
      {{end}}
    </div>

    <!-- Attractions & Landmarks -->
    <div class="attractions-card">
      <h3>Attractions &amp; landmarks</h3>
      {{if .Attractions}}
        {{range .Attractions}}
        <div class="attraction-item">
          <div>
            <span class="attraction-name">{{.Name}}</span>
            <span class="attraction-kinds">{{.Kinds}}</span>
          </div>
        </div>
        {{end}}
      {{else}}
      <div class="empty-state">
        <p>No attractions found for this destination.</p>
      </div>
      {{end}}
    </div>

  </div>
</div>

<script>
  // Country name pass করো JS এ
  var COUNTRY_NAME = "{{.Country.Name}}";
  var IS_LOGGED_IN = {{if .IsLoggedIn}}true{{else}}false{{end}};
</script>
<script src="/static/js/destination.js"></script>