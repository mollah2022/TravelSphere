<header class="site-header">
  <div class="header-inner">
    <a href="/" class="logo">🌍 TravelSphere</a>
    <nav class="main-nav">
      <a href="/" class="{{if .NavHome}}active{{end}}">Home</a>
      <a href="/countries" class="{{if .NavCountries}}active{{end}}">Countries</a>
      {{if .IsLoggedIn}}
        <a href="/wishlist" class="{{if .NavWishlist}}active{{end}}">Wishlist</a>
        <a href="/dashboard" class="{{if .NavDashboard}}active{{end}}">Dashboard</a>
        <a href="/logout" class="btn-logout">Logout ({{.Username}})</a>
      {{else}}
        <a href="/login" class="btn-login">Login</a>
      {{end}}
    </nav>
  </div>
</header>