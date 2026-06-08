<div class="login-wrap">
  <div class="login-card">
    <h2>Welcome back</h2>
    <p>Enter any username to continue — no password required.</p>

    {{if .Error}}
    <div class="alert-error">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/login">
      <input type="hidden" name="redirect" value="{{.RedirectURL}}">

      <div class="form-group">
        <label for="username">Username</label>
        <input
          type="text"
          id="username"
          name="username"
          placeholder="e.g. beta, john, sarah"
          autofocus
          maxlength="30"
        >
      </div>

      <button type="submit" class="btn-primary">Continue</button>
    </form>
  </div>
</div>