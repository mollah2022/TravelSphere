<nav class="navbar">
  <a href="/" class="navbar-brand">TravelSphere</a>

  <div class="navbar-menu">
    <a href="/" {{if .NavHome}}class="active"{{end}}>Home</a>
    <a href="/countries" {{if .NavCountries}}class="active"{{end}}>Countries</a>
    {{if .IsLoggedIn}}
      <a href="/wishlist" {{if .NavWishlist}}class="active"{{end}}>Wishlist</a>
      <a href="/dashboard" {{if .NavDashboard}}class="active"{{end}}>Dashboard</a>
    {{end}}
  </div>

  <div class="navbar-right">
    {{if .IsLoggedIn}}
      <span>Hi, <strong>{{.Username}}</strong></span>
      <a href="/logout" class="btn-logout">Logout</a>
    {{else}}
      <a href="/login" class="btn-login">Login</a>
    {{end}}
  </div>
</nav>