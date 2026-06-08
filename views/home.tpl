<section class="hero">
  <div class="hero-inner">
    <h1>Discover your next destination</h1>
    <p>Search countries, explore attractions, and curate your personal travel wishlist.</p>

    <div class="hero-label">WHERE TO NEXT?</div>
    <div class="search-box">
      <input
        type="text"
        id="home-search"
        placeholder="Search country or capital..."
        autocomplete="off"
      >
      <div class="search-suggestions" id="search-suggestions"></div>
    </div>
  </div>
</section>

<div class="container">
  <h2 class="section-title">Featured destinations</h2>

  {{if .FeaturedCountries}}
  <div class="featured-grid">
    {{range .FeaturedCountries}}
    <a href="/countries/{{.Slug}}" class="featured-card">
      <img
        src="{{.FlagURL}}"
        alt="{{.FlagAlt}}"
        loading="lazy"
      >
      <div class="featured-card-body">
        <div class="featured-card-name">{{.Name}}</div>
        <div class="featured-card-sub">{{.Capital}} &middot; {{.Region}}</div>
      </div>
    </a>
    {{end}}
  </div>
  {{else}}
  <div class="empty-state"><p>Featured destinations unavailable.</p></div>
  {{end}}

  <h2 class="section-title" style="margin-top:2.5rem;">Popular attractions</h2>
  {{if .PopularAttractions}}
    {{range .PopularAttractions}}
    <div class="attraction-item">
      <div>
        <span class="attraction-name">{{.Name}}</span>
        <span class="attraction-kinds">{{.Kinds}}</span>
      </div>
      <span class="attraction-country">{{.Country}}</span>
    </div>
    {{end}}
  {{else}}
  <div class="empty-state"><p>Attractions unavailable.</p></div>
  {{end}}
</div>

<script src="/static/js/home.js"></script>