export default {
    data() {
      return {
        playername: "",
        password: "",
        loading: false,
        error: false,
      }
    },
    methods: {
        async tryLogin() {
          if (this.playername.length == 0 || this.password.length == 0) {
              this.error = true;
              return
          }
          this.error = false;
          this.loading = true;
          let url = `/login`;
          const requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ playername: this.playername, password: this.password })
          };
          try {
            const login_response = await fetch(url, requestOptions);
            if (!login_response.ok) {
              this.error = true;
              this.loading = false;
            } else {
                this.loading = false;
                window.location.href = "/dashboard.html"
            }
          } catch (err) {
            this.loading = false;
            this.error = true;
          }
          
        },
    },
    template: `

    <div class="form-floating m-1">
        <input type="text" class="form-control" v-model="playername" name="playername" placeholder="Playername" required>
        <label for="floatingInput">User</label>
    </div>
    <div class="form-floating m-1">
        <input type="password" class="form-control" v-model="password" name="password" placeholder="Password" required>
        <label for="floatingPassword">
            Password
        </label>
    </div>

    <div v-if="error">
      <div class="alert alert-danger p-2 m-1" role="alert">
          Login failed.
      </div>
    </div>

    <div class="checkbox mb-3">
        <label>
            <input type="checkbox" value="remember-me"> Remember me
        </label>
    </div>
    <div class="d-flex justify-content-between">
        <button class="btn btn-lg btn-primary w-50 m-1" type="button" v-on:click="tryLogin">
        <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
            Login
        </button>
        <button class="btn btn-lg btn-primary w-50 m-1">Sign up</button>
    </div>
    `
}